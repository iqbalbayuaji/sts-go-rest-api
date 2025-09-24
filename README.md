# Recipe Notes - Go REST API Application

A simple client-server application for managing recipe notes using Go's `net/http` library with a web interface.

## Features

- **CRUD Operations**: Create, Read, Update, and Delete recipes
- **REST API**: Clean RESTful endpoints
- **Web Interface**: User-friendly HTML interface with JavaScript
- **JSON Storage**: Data persistence using JSON files
- **Responsive Design**: Works on desktop and mobile devices

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/recipes` | Get all recipes |
| POST | `/api/recipes` | Create a new recipe |
| PUT | `/api/recipes` | Update an existing recipe |
| DELETE | `/api/recipes/{id}` | Delete a recipe by ID |

## Project Structure

```
recipe-api/
├── main.go              # Server entry point
├── go.mod               # Go module file
├── handlers/            # HTTP request handlers
│   └── recipe_handler.go
├── models/              # Data structures
│   └── recipe.go
├── storage/             # Data persistence layer
│   └── json_storage.go
├── static/              # Web interface files
│   ├── index.html       # Main HTML page
│   ├── styles.css       # CSS styling
│   └── script.js        # JavaScript functionality
├── data/                # JSON data storage
│   └── recipes.json     # Recipe data file
└── README.md           # This file
```

## Recipe Data Structure

```json
{
  "id": "unique-uuid",
  "name": "Recipe Name",
  "ingredients": ["ingredient1", "ingredient2"],
  "instructions": "Step-by-step cooking instructions",
  "cooking_time": "30 minutes",
  "servings": 4,
  "category": "main course",
  "created_at": "2023-01-01T12:00:00Z",
  "updated_at": "2023-01-01T12:00:00Z"
}
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git (optional)

### Installation

1. **Clone or download the project**:
   ```bash
   git clone <repository-url>
   cd recipe-api
   ```

2. **Install dependencies**:
   ```bash
   go mod tidy
   ```

3. **Run the application**:
   ```bash
   go run main.go
   ```

4. **Access the application**:
   - Web Interface: http://localhost:8080
   - API Base URL: http://localhost:8080/api/recipes

## Usage

### Web Interface

1. Open your browser and go to `http://localhost:8080`
2. Use the form on the left to add new recipes
3. View all recipes on the right side
4. Click "Edit" to modify a recipe
5. Click "Delete" to remove a recipe
6. Click "Refresh" to reload the recipe list

### API Usage Examples

#### Get All Recipes
```bash
curl -X GET http://localhost:8080/api/recipes
```

#### Create a New Recipe
```bash
curl -X POST http://localhost:8080/api/recipes \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Chocolate Cake",
    "ingredients": ["2 cups flour", "1 cup sugar", "3 eggs"],
    "instructions": "1. Preheat oven to 350°F\n2. Mix ingredients\n3. Bake for 30 minutes",
    "cooking_time": "45 minutes",
    "servings": 8,
    "category": "dessert"
  }'
```

#### Update a Recipe
```bash
curl -X PUT http://localhost:8080/api/recipes \
  -H "Content-Type: application/json" \
  -d '{
    "id": "recipe-uuid-here",
    "name": "Updated Recipe Name",
    "ingredients": ["updated ingredients"],
    "instructions": "Updated instructions",
    "cooking_time": "30 minutes",
    "servings": 4,
    "category": "main course"
  }'
```

#### Delete a Recipe
```bash
curl -X DELETE http://localhost:8080/api/recipes/recipe-uuid-here
```

## API Response Format

### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... }
}
```

### Error Response
```json
{
  "success": false,
  "error": "Error message description"
}
```

## Features Explained

### Server Features
- **CORS Support**: Allows cross-origin requests
- **JSON Storage**: Persistent data storage in JSON format
- **UUID Generation**: Automatic unique ID generation for recipes
- **Input Validation**: Server-side validation for all recipe fields
- **Error Handling**: Comprehensive error handling with appropriate HTTP status codes
- **Thread Safety**: Concurrent access protection using mutex locks

### Client Features
- **Responsive Design**: Works on all screen sizes
- **Real-time Updates**: Automatic UI updates after operations
- **Form Validation**: Client-side validation before sending requests
- **Loading States**: Visual feedback during API calls
- **Toast Notifications**: Success and error messages
- **Edit Mode**: In-place editing of existing recipes

## Development

### Adding New Features

1. **New API Endpoints**: Add handlers in `handlers/recipe_handler.go`
2. **Data Models**: Extend structures in `models/recipe.go`
3. **Storage Operations**: Add methods in `storage/json_storage.go`
4. **UI Features**: Update HTML, CSS, and JavaScript in `static/` folder

### Testing

Test the API endpoints using tools like:
- **curl** (command line)
- **Postman** (GUI)
- **Thunder Client** (VS Code extension)
- **Browser Developer Tools** (for web interface)

## Troubleshooting

### Common Issues

1. **Port 8080 already in use**:
   - Change the port in `main.go`: `http.ListenAndServe(":8081", nil)`

2. **Permission denied for data directory**:
   - Ensure the application has write permissions to the `data/` directory

3. **Module not found errors**:
   - Run `go mod tidy` to download dependencies

4. **CORS errors in browser**:
   - The application includes CORS headers, but ensure you're accessing via `http://localhost:8080`

## Learning Objectives

This project demonstrates:
- Go web server development with `net/http`
- RESTful API design principles
- JSON data handling in Go
- File-based data persistence
- Frontend-backend integration
- Error handling and validation
- Concurrent programming with mutexes
- Modern web interface development

## Next Steps

Consider extending the application with:
- Database integration (PostgreSQL, MySQL)
- User authentication and authorization
- Image upload for recipes
- Recipe search and filtering
- Recipe categories and tags
- Recipe sharing functionality
- Mobile app development
- Docker containerization
