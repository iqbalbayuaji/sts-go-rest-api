package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"recipe-api/auth"
	"recipe-api/database"
	_ "recipe-api/docs"
	"recipe-api/handlers"
	"recipe-api/models"
	"recipe-api/storage"
	"time"

	httpSwagger "github.com/swaggo/http-swagger"
	"gopkg.in/yaml.v2"
)

// @title Recipe API
// @version 1.0
// @description A REST API for managing recipes with authentication
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database connection
	if err := database.InitDB(config.Database); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Run database migrations
	if err := database.RunMigrations(""); err != nil {
		log.Printf("Warning: Failed to run migrations: %v", err)
	}

	// Initialize storage
	recipeStorage := storage.NewPostgresStorage()
	userStorage := storage.NewPostgresUserStorage()

	// Initialize authentication service
	authService, err := auth.NewAuthService("config.yaml", userStorage)
	if err != nil {
		log.Fatal("Failed to initialize auth service:", err)
	}

	// Initialize handlers
	recipeHandler := handlers.NewRecipeHandler(recipeStorage)
	authHandler := handlers.NewAuthHandler(authService)

	// Setup authentication routes (public)
	http.HandleFunc("/api/login", authHandler.HandleLogin)
	http.HandleFunc("/api/logout", authHandler.HandleLogout)

	// Setup protected routes (require authentication)
	http.HandleFunc("/api/recipes", authHandler.AuthMiddleware(recipeHandler.HandleRecipes))
	http.HandleFunc("/api/recipes/", authHandler.AuthMiddleware(recipeHandler.HandleRecipeByID))

	// Setup Swagger documentation
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Serve static files (HTML, CSS, JS)
	http.Handle("/", http.FileServer(http.Dir("static/")))

	// Start token cleanup routine
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				authService.CleanupExpiredTokens()
				log.Printf("Cleaned up expired tokens. Active tokens: %d", authService.GetActiveTokensCount())
			}
		}
	}()

	// Start server
	log.Println("Server starting on :8080")
	log.Println("Authentication endpoints:")
	log.Println("  POST /api/login - Login with username/password")
	log.Println("  POST /api/logout - Logout (invalidate token)")
	log.Println("Protected API endpoints:")
	log.Println("  GET/POST/PUT /api/recipes - Recipe operations (requires Bearer token)")
	log.Println("  DELETE /api/recipes/{id} - Delete recipe (requires Bearer token)")
	log.Println("API Documentation:")
	log.Println("  Swagger UI: http://localhost:8080/swagger/")
	log.Println("Web interface at: http://localhost:8080")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// loadConfig loads configuration from YAML file
func loadConfig(configPath string) (*models.Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config models.Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Set default values if not specified
	if config.TokenExpiryHours == 0 {
		config.TokenExpiryHours = 24
	}
	if config.JWTSecret == "" {
		config.JWTSecret = "default-secret-change-this"
	}

	return &config, nil
}
