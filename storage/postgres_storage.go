package storage

import (
	"database/sql"
	"fmt"
	"recipe-api/database"
	"recipe-api/models"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// PostgresStorage handles PostgreSQL operations for recipes
type PostgresStorage struct {
	db *sql.DB
}

// NewPostgresStorage creates a new PostgreSQL storage instance
func NewPostgresStorage() *PostgresStorage {
	return &PostgresStorage{
		db: database.GetDB(),
	}
}

// GetAllRecipes retrieves all recipes from the database
func (ps *PostgresStorage) GetAllRecipes() ([]models.Recipe, error) {
	query := `
		SELECT id, name, ingredients, instructions, cooking_time, servings, category,
		       created_at, updated_at, created_by, updated_by
		FROM recipes
		ORDER BY created_at DESC
	`

	rows, err := ps.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipes: %v", err)
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		err := rows.Scan(
			&recipe.ID, &recipe.Name, pq.Array(&recipe.Ingredients), &recipe.Instructions,
			&recipe.CookingTime, &recipe.Servings, &recipe.Category,
			&recipe.CreatedAt, &recipe.UpdatedAt, &recipe.CreatedBy, &recipe.UpdatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %v", err)
		}
		recipes = append(recipes, recipe)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating recipes: %v", err)
	}

	return recipes, nil
}

// GetRecipeByID retrieves a specific recipe by ID
func (ps *PostgresStorage) GetRecipeByID(id string) (*models.Recipe, error) {
	query := `
		SELECT id, name, ingredients, instructions, cooking_time, servings, category,
		       created_at, updated_at, created_by, updated_by
		FROM recipes
		WHERE id = $1
	`

	var recipe models.Recipe
	err := ps.db.QueryRow(query, id).Scan(
		&recipe.ID, &recipe.Name, pq.Array(&recipe.Ingredients), &recipe.Instructions,
		&recipe.CookingTime, &recipe.Servings, &recipe.Category,
		&recipe.CreatedAt, &recipe.UpdatedAt, &recipe.CreatedBy, &recipe.UpdatedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("recipe with ID %s not found", id)
		}
		return nil, fmt.Errorf("failed to get recipe: %v", err)
	}

	return &recipe, nil
}

// SaveRecipe adds a new recipe or updates an existing one
func (ps *PostgresStorage) SaveRecipe(recipe models.Recipe, userID *int) error {
	// Check if recipe exists
	existingRecipe, err := ps.GetRecipeByID(recipe.ID)
	if err != nil && err.Error() != fmt.Sprintf("recipe with ID %s not found", recipe.ID) {
		return err
	}

	if existingRecipe != nil {
		// Update existing recipe
		return ps.updateRecipe(recipe, userID)
	} else {
		// Create new recipe
		return ps.createRecipe(recipe, userID)
	}
}

// createRecipe creates a new recipe
func (ps *PostgresStorage) createRecipe(recipe models.Recipe, userID *int) error {
	// Generate new UUID if not provided
	if recipe.ID == "" {
		recipe.ID = uuid.New().String()
	}

	query := `
		INSERT INTO recipes (id, name, ingredients, instructions, cooking_time, servings, category, created_by, updated_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING created_at, updated_at
	`

	err := ps.db.QueryRow(
		query,
		recipe.ID, recipe.Name, pq.Array(recipe.Ingredients), recipe.Instructions,
		recipe.CookingTime, recipe.Servings, recipe.Category, userID, userID,
	).Scan(&recipe.CreatedAt, &recipe.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create recipe: %v", err)
	}

	return nil
}

// updateRecipe updates an existing recipe
func (ps *PostgresStorage) updateRecipe(recipe models.Recipe, userID *int) error {
	query := `
		UPDATE recipes 
		SET name = $2, ingredients = $3, instructions = $4, cooking_time = $5, 
		    servings = $6, category = $7, updated_by = $8, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`

	err := ps.db.QueryRow(
		query,
		recipe.ID, recipe.Name, pq.Array(recipe.Ingredients), recipe.Instructions,
		recipe.CookingTime, recipe.Servings, recipe.Category, userID,
	).Scan(&recipe.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to update recipe: %v", err)
	}

	return nil
}

// DeleteRecipe removes a recipe by ID
func (ps *PostgresStorage) DeleteRecipe(id string) error {
	query := `DELETE FROM recipes WHERE id = $1`

	result, err := ps.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete recipe: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("recipe with ID %s not found", id)
	}

	return nil
}

// GetRecipesByCategory retrieves recipes by category
func (ps *PostgresStorage) GetRecipesByCategory(category string) ([]models.Recipe, error) {
	query := `
		SELECT id, name, ingredients, instructions, cooking_time, servings, category,
		       created_at, updated_at, created_by, updated_by
		FROM recipes
		WHERE category = $1
		ORDER BY created_at DESC
	`

	rows, err := ps.db.Query(query, category)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipes by category: %v", err)
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		err := rows.Scan(
			&recipe.ID, &recipe.Name, pq.Array(&recipe.Ingredients), &recipe.Instructions,
			&recipe.CookingTime, &recipe.Servings, &recipe.Category,
			&recipe.CreatedAt, &recipe.UpdatedAt, &recipe.CreatedBy, &recipe.UpdatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %v", err)
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}

// SearchRecipes searches recipes by name or ingredients
func (ps *PostgresStorage) SearchRecipes(searchTerm string) ([]models.Recipe, error) {
	query := `
		SELECT id, name, ingredients, instructions, cooking_time, servings, category,
		       created_at, updated_at, created_by, updated_by
		FROM recipes
		WHERE name ILIKE $1 OR $2 = ANY(ingredients)
		ORDER BY created_at DESC
	`

	searchPattern := "%" + searchTerm + "%"
	rows, err := ps.db.Query(query, searchPattern, searchTerm)
	if err != nil {
		return nil, fmt.Errorf("failed to search recipes: %v", err)
	}
	defer rows.Close()

	var recipes []models.Recipe
	for rows.Next() {
		var recipe models.Recipe
		err := rows.Scan(
			&recipe.ID, &recipe.Name, pq.Array(&recipe.Ingredients), &recipe.Instructions,
			&recipe.CookingTime, &recipe.Servings, &recipe.Category,
			&recipe.CreatedAt, &recipe.UpdatedAt, &recipe.CreatedBy, &recipe.UpdatedBy,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan recipe: %v", err)
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil
}
