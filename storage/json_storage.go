package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"recipe-api/models"
	"sync"
)

// JSONStorage handles JSON file operations for recipes
type JSONStorage struct {
	filePath string
	mutex    sync.RWMutex
}

// NewJSONStorage creates a new JSON storage instance
func NewJSONStorage(filePath string) *JSONStorage {
	storage := &JSONStorage{
		filePath: filePath,
	}
	
	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("Warning: Could not create directory %s: %v\n", dir, err)
	}
	
	// Initialize file if it doesn't exist
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		storage.saveRecipes([]models.Recipe{})
	}
	
	return storage
}

// GetAllRecipes retrieves all recipes from the JSON file
func (js *JSONStorage) GetAllRecipes() ([]models.Recipe, error) {
	js.mutex.RLock()
	defer js.mutex.RUnlock()
	
	data, err := os.ReadFile(js.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipes file: %v", err)
	}
	
	var recipes []models.Recipe
	if err := json.Unmarshal(data, &recipes); err != nil {
		return nil, fmt.Errorf("failed to parse recipes: %v", err)
	}
	
	return recipes, nil
}

// GetRecipeByID retrieves a specific recipe by ID
func (js *JSONStorage) GetRecipeByID(id string) (*models.Recipe, error) {
	recipes, err := js.GetAllRecipes()
	if err != nil {
		return nil, err
	}
	
	for _, recipe := range recipes {
		if recipe.ID == id {
			return &recipe, nil
		}
	}
	
	return nil, fmt.Errorf("recipe with ID %s not found", id)
}

// SaveRecipe adds a new recipe or updates an existing one
func (js *JSONStorage) SaveRecipe(recipe models.Recipe) error {
	js.mutex.Lock()
	defer js.mutex.Unlock()
	
	recipes, err := js.getAllRecipesUnsafe()
	if err != nil {
		return err
	}
	
	// Check if recipe exists (for update)
	found := false
	for i, existingRecipe := range recipes {
		if existingRecipe.ID == recipe.ID {
			recipes[i] = recipe
			found = true
			break
		}
	}
	
	// If not found, add as new recipe
	if !found {
		recipes = append(recipes, recipe)
	}
	
	return js.saveRecipes(recipes)
}

// DeleteRecipe removes a recipe by ID
func (js *JSONStorage) DeleteRecipe(id string) error {
	js.mutex.Lock()
	defer js.mutex.Unlock()
	
	recipes, err := js.getAllRecipesUnsafe()
	if err != nil {
		return err
	}
	
	// Find and remove the recipe
	found := false
	for i, recipe := range recipes {
		if recipe.ID == id {
			recipes = append(recipes[:i], recipes[i+1:]...)
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("recipe with ID %s not found", id)
	}
	
	return js.saveRecipes(recipes)
}

// getAllRecipesUnsafe is an internal method that doesn't use mutex (for internal use only)
func (js *JSONStorage) getAllRecipesUnsafe() ([]models.Recipe, error) {
	data, err := os.ReadFile(js.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read recipes file: %v", err)
	}
	
	var recipes []models.Recipe
	if err := json.Unmarshal(data, &recipes); err != nil {
		return nil, fmt.Errorf("failed to parse recipes: %v", err)
	}
	
	return recipes, nil
}

// saveRecipes saves recipes to the JSON file
func (js *JSONStorage) saveRecipes(recipes []models.Recipe) error {
	data, err := json.MarshalIndent(recipes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal recipes: %v", err)
	}
	
	if err := os.WriteFile(js.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write recipes file: %v", err)
	}
	
	return nil
}
