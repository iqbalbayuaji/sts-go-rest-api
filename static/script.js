// API Base URL
const API_BASE = '/api/recipes';

// Global variables
let currentEditingId = null;
let recipes = [];

// DOM Elements
const recipeForm = document.getElementById('recipe-form');
const recipesContainer = document.getElementById('recipes-container');
const loadingElement = document.getElementById('loading');
const submitBtn = document.getElementById('submit-btn');
const cancelBtn = document.getElementById('cancel-btn');
const formTitle = document.getElementById('form-title');
const refreshBtn = document.getElementById('refresh-btn');

// Initialize app
document.addEventListener('DOMContentLoaded', function() {
    loadRecipes();
    setupEventListeners();
});

// Setup event listeners
function setupEventListeners() {
    recipeForm.addEventListener('submit', handleFormSubmit);
    cancelBtn.addEventListener('click', cancelEdit);
    refreshBtn.addEventListener('click', loadRecipes);
}

// Load all recipes
async function loadRecipes() {
    try {
        showLoading(true);
        const response = await fetch(API_BASE);
        const data = await response.json();
        
        if (data.success) {
            recipes = data.data || [];
            displayRecipes();
        } else {
            throw new Error(data.error || 'Failed to load recipes');
        }
    } catch (error) {
        console.error('Error loading recipes:', error);
        showToast('Failed to load recipes: ' + error.message, 'error');
        recipesContainer.innerHTML = '<div class="empty-state"><h3>Failed to load recipes</h3><p>Please try refreshing the page.</p></div>';
    } finally {
        showLoading(false);
    }
}

// Display recipes in the UI
function displayRecipes() {
    if (recipes.length === 0) {
        recipesContainer.innerHTML = `
            <div class="empty-state">
                <h3>No recipes yet</h3>
                <p>Add your first recipe using the form on the left!</p>
            </div>
        `;
        return;
    }

    recipesContainer.innerHTML = recipes.map(recipe => `
        <div class="recipe-card" data-id="${recipe.id}">
            <div class="recipe-header">
                <div>
                    <h3 class="recipe-title">${escapeHtml(recipe.name)}</h3>
                    <div class="recipe-meta">
                        <span>⏱️ ${escapeHtml(recipe.cooking_time)}</span>
                        <span>👥 ${recipe.servings} servings</span>
                        <span class="recipe-category">${escapeHtml(recipe.category)}</span>
                    </div>
                </div>
            </div>
            
            <div class="recipe-ingredients">
                <h4>Ingredients:</h4>
                <ul class="ingredients-list">
                    ${recipe.ingredients.map(ingredient => `<li>${escapeHtml(ingredient)}</li>`).join('')}
                </ul>
            </div>
            
            <div class="recipe-instructions">
                <h4>Instructions:</h4>
                <p>${escapeHtml(recipe.instructions)}</p>
            </div>
            
            <div class="recipe-actions">
                <button class="btn-edit" onclick="editRecipe('${recipe.id}')">
                    ✏️ Edit
                </button>
                <button class="btn-danger" onclick="deleteRecipe('${recipe.id}')">
                    🗑️ Delete
                </button>
            </div>
        </div>
    `).join('');
}

// Handle form submission
async function handleFormSubmit(e) {
    e.preventDefault();
    
    const formData = new FormData(recipeForm);
    const recipeData = {
        name: formData.get('name').trim(),
        ingredients: formData.get('ingredients').split('\n').map(i => i.trim()).filter(i => i),
        instructions: formData.get('instructions').trim(),
        cooking_time: formData.get('cooking_time').trim(),
        servings: parseInt(formData.get('servings')),
        category: formData.get('category')
    };

    // Validation
    if (!recipeData.name || !recipeData.instructions || !recipeData.cooking_time || 
        !recipeData.category || recipeData.ingredients.length === 0 || recipeData.servings < 1) {
        showToast('Please fill in all required fields', 'error');
        return;
    }

    try {
        submitBtn.disabled = true;
        submitBtn.textContent = currentEditingId ? 'Updating...' : 'Adding...';

        let response;
        if (currentEditingId) {
            // Update existing recipe
            recipeData.id = currentEditingId;
            response = await fetch(API_BASE, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(recipeData)
            });
        } else {
            // Create new recipe
            response = await fetch(API_BASE, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(recipeData)
            });
        }

        const data = await response.json();
        
        if (data.success) {
            showToast(data.message, 'success');
            recipeForm.reset();
            cancelEdit();
            loadRecipes();
        } else {
            throw new Error(data.error || 'Failed to save recipe');
        }
    } catch (error) {
        console.error('Error saving recipe:', error);
        showToast('Failed to save recipe: ' + error.message, 'error');
    } finally {
        submitBtn.disabled = false;
        submitBtn.textContent = currentEditingId ? 'Update Recipe' : 'Add Recipe';
    }
}

// Edit recipe
function editRecipe(id) {
    const recipe = recipes.find(r => r.id === id);
    if (!recipe) {
        showToast('Recipe not found', 'error');
        return;
    }

    currentEditingId = id;
    
    // Populate form
    document.getElementById('name').value = recipe.name;
    document.getElementById('ingredients').value = recipe.ingredients.join('\n');
    document.getElementById('instructions').value = recipe.instructions;
    document.getElementById('cooking_time').value = recipe.cooking_time;
    document.getElementById('servings').value = recipe.servings;
    document.getElementById('category').value = recipe.category;

    // Update UI
    formTitle.textContent = 'Edit Recipe';
    submitBtn.textContent = 'Update Recipe';
    cancelBtn.style.display = 'inline-block';

    // Scroll to form
    document.querySelector('.form-section').scrollIntoView({ behavior: 'smooth' });
}

// Cancel edit
function cancelEdit() {
    currentEditingId = null;
    recipeForm.reset();
    formTitle.textContent = 'Add New Recipe';
    submitBtn.textContent = 'Add Recipe';
    cancelBtn.style.display = 'none';
}

// Delete recipe
async function deleteRecipe(id) {
    if (!confirm('Are you sure you want to delete this recipe? This action cannot be undone.')) {
        return;
    }

    try {
        const response = await fetch(`${API_BASE}/${id}`, {
            method: 'DELETE'
        });

        const data = await response.json();
        
        if (data.success) {
            showToast(data.message, 'success');
            
            // If we're editing this recipe, cancel edit
            if (currentEditingId === id) {
                cancelEdit();
            }
            
            loadRecipes();
        } else {
            throw new Error(data.error || 'Failed to delete recipe');
        }
    } catch (error) {
        console.error('Error deleting recipe:', error);
        showToast('Failed to delete recipe: ' + error.message, 'error');
    }
}

// Show/hide loading
function showLoading(show) {
    loadingElement.style.display = show ? 'block' : 'none';
}

// Show toast notification
function showToast(message, type = 'success') {
    const toast = document.getElementById('toast');
    toast.textContent = message;
    toast.className = `toast ${type}`;
    toast.classList.add('show');

    setTimeout(() => {
        toast.classList.remove('show');
    }, 3000);
}

// Escape HTML to prevent XSS
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Format date for display
function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
}
