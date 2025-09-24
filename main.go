package main

import (
	"log"
	"net/http"
	"recipe-api/handlers"
	"recipe-api/storage"
)

func main() {
	// Initialize storage
	store := storage.NewJSONStorage("data/recipes.json")

	// Initialize handlers
	recipeHandler := handlers.NewRecipeHandler(store)

	// Setup routes
	http.HandleFunc("/api/recipes", recipeHandler.HandleRecipes)
	http.HandleFunc("/api/recipes/", recipeHandler.HandleRecipeByID)

	// Serve static files (HTML, CSS, JS)
	http.Handle("/", http.FileServer(http.Dir("static/")))

	// Start server
	log.Println("Server starting on :8080")
	log.Println("API available at: http://localhost:8080/api/recipes")
	log.Println("Web interface at: http://localhost:8080")
	
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
