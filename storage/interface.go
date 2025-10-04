package storage

import "recipe-api/models"

// RecipeStorage defines the interface for recipe storage operations
type RecipeStorage interface {
	GetAllRecipes() ([]models.Recipe, error)
	GetRecipeByID(id string) (*models.Recipe, error)
	SaveRecipe(recipe models.Recipe, userID *int) error
	DeleteRecipe(id string) error
	GetRecipesByCategory(category string) ([]models.Recipe, error)
	SearchRecipes(searchTerm string) ([]models.Recipe, error)
}

// UserStorage defines the interface for user storage operations
type UserStorage interface {
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateUser(user models.User) error
	UpdateUser(user models.User) error
	DeleteUser(id int) error
	ValidateCredentials(username, password string) (*models.User, error)
}
