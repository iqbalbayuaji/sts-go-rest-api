package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recipe-api/models"
	"recipe-api/storage"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RecipeHandler handles HTTP requests for recipes
type RecipeHandler struct {
	storage *storage.JSONStorage
}

// NewRecipeHandler creates a new recipe handler
func NewRecipeHandler(storage *storage.JSONStorage) *RecipeHandler {
	return &RecipeHandler{
		storage: storage,
	}
}

// HandleRecipes handles requests to /api/recipes (GET and POST)
func (rh *RecipeHandler) HandleRecipes(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	switch r.Method {
	case "GET":
		rh.getAllRecipes(w, r)
	case "POST":
		rh.createRecipe(w, r)
	case "PUT":
		rh.updateRecipe(w, r)
	default:
		rh.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandleRecipeByID handles requests to /api/recipes/{id} (DELETE)
func (rh *RecipeHandler) HandleRecipeByID(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extract ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/recipes/")
	if path == "" {
		rh.sendError(w, "Recipe ID is required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		rh.getRecipeByID(w, r, path)
	case "DELETE":
		rh.deleteRecipe(w, r, path)
	default:
		rh.sendError(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// getAllRecipes handles GET /api/recipes
func (rh *RecipeHandler) getAllRecipes(w http.ResponseWriter, r *http.Request) {
	recipes, err := rh.storage.GetAllRecipes()
	if err != nil {
		rh.sendError(w, fmt.Sprintf("Failed to get recipes: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.APIResponse{
		Success: true,
		Message: "Recipes retrieved successfully",
		Data:    recipes,
	}

	rh.sendJSON(w, response, http.StatusOK)
}

// getRecipeByID handles GET /api/recipes/{id}
func (rh *RecipeHandler) getRecipeByID(w http.ResponseWriter, r *http.Request, id string) {
	recipe, err := rh.storage.GetRecipeByID(id)
	if err != nil {
		rh.sendError(w, fmt.Sprintf("Recipe not found: %v", err), http.StatusNotFound)
		return
	}

	response := models.APIResponse{
		Success: true,
		Message: "Recipe retrieved successfully",
		Data:    recipe,
	}

	rh.sendJSON(w, response, http.StatusOK)
}

// createRecipe handles POST /api/recipes
func (rh *RecipeHandler) createRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		rh.sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate recipe
	if err := recipe.Validate(); err != nil {
		rh.sendError(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	// Generate ID and timestamps
	recipe.ID = uuid.New().String()
	recipe.CreatedAt = time.Now()
	recipe.UpdatedAt = time.Now()

	// Save recipe
	if err := rh.storage.SaveRecipe(recipe); err != nil {
		rh.sendError(w, fmt.Sprintf("Failed to save recipe: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.APIResponse{
		Success: true,
		Message: "Recipe created successfully",
		Data:    recipe,
	}

	rh.sendJSON(w, response, http.StatusCreated)
}

// updateRecipe handles PUT /api/recipes
func (rh *RecipeHandler) updateRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		rh.sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Check if recipe exists
	existingRecipe, err := rh.storage.GetRecipeByID(recipe.ID)
	if err != nil {
		rh.sendError(w, "Recipe not found", http.StatusNotFound)
		return
	}

	// Validate recipe
	if err := recipe.Validate(); err != nil {
		rh.sendError(w, fmt.Sprintf("Validation error: %v", err), http.StatusBadRequest)
		return
	}

	// Keep original creation time, update modification time
	recipe.CreatedAt = existingRecipe.CreatedAt
	recipe.UpdatedAt = time.Now()

	// Save updated recipe
	if err := rh.storage.SaveRecipe(recipe); err != nil {
		rh.sendError(w, fmt.Sprintf("Failed to update recipe: %v", err), http.StatusInternalServerError)
		return
	}

	response := models.APIResponse{
		Success: true,
		Message: "Recipe updated successfully",
		Data:    recipe,
	}

	rh.sendJSON(w, response, http.StatusOK)
}

// deleteRecipe handles DELETE /api/recipes/{id}
func (rh *RecipeHandler) deleteRecipe(w http.ResponseWriter, r *http.Request, id string) {
	if err := rh.storage.DeleteRecipe(id); err != nil {
		rh.sendError(w, fmt.Sprintf("Failed to delete recipe: %v", err), http.StatusNotFound)
		return
	}

	response := models.APIResponse{
		Success: true,
		Message: "Recipe deleted successfully",
	}

	rh.sendJSON(w, response, http.StatusOK)
}

// sendJSON sends a JSON response
func (rh *RecipeHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// sendError sends an error response
func (rh *RecipeHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	response := models.APIError{
		Success: false,
		Error:   message,
	}
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
