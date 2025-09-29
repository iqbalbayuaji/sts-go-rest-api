package main

import (
	"log"
	"net/http"
	"recipe-api/auth"
	"recipe-api/handlers"
	"recipe-api/storage"
	"time"
)

func main() {
	// Initialize authentication service
	authService, err := auth.NewAuthService("config.yaml")
	if err != nil {
		log.Fatal("Failed to initialize auth service:", err)
	}

	// Initialize storage
	store := storage.NewJSONStorage("data/recipes.json")

	// Initialize handlers
	recipeHandler := handlers.NewRecipeHandler(store)
	authHandler := handlers.NewAuthHandler(authService)

	// Setup authentication routes (public)
	http.HandleFunc("/api/login", authHandler.HandleLogin)
	http.HandleFunc("/api/logout", authHandler.HandleLogout)

	// Setup protected routes (require authentication)
	http.HandleFunc("/api/recipes", authHandler.AuthMiddleware(recipeHandler.HandleRecipes))
	http.HandleFunc("/api/recipes/", authHandler.AuthMiddleware(recipeHandler.HandleRecipeByID))

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
	log.Println("Web interface at: http://localhost:8080")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
