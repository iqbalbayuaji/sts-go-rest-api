# Recipe Notes - Go REST API Application

A secure client-server application for managing recipe notes using Go's `net/http` library with authentication and a web interface.

## Features

- **🔐 Authentication System**: Token-based authentication with login/logout
- **🛡️ Secure API**: All recipe endpoints protected with Bearer token authentication
- **📝 CRUD Operations**: Create, Read, Update, and Delete recipes
- **🌐 REST API**: Clean RESTful endpoints with proper HTTP status codes
- **💻 Web Interface**: User-friendly HTML interface with JavaScript
- **📁 JSON Storage**: Data persistence using JSON files
- **⚙️ YAML Configuration**: User management via configuration file
- **📱 Responsive Design**: Works on desktop and mobile devices
- **🔄 Token Management**: Automatic token cleanup and expiration handling

## API Endpoints

### Authentication Endpoints (Public)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/login` | Login with username/password, returns Bearer token |
| POST | `/api/logout` | Logout and invalidate token |

### Recipe Endpoints (Protected - Requires Bearer Token)
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
├── config.yaml          # Authentication configuration (users, JWT secret)
├── config.yaml.example  # Sample configuration file
├── auth/                # Authentication services
│   └── auth_service.go  # Token management and validation
├── handlers/            # HTTP request handlers
│   ├── recipe_handler.go # Recipe CRUD operations
│   └── auth_handler.go   # Login/logout and middleware
├── models/              # Data structures
│   ├── recipe.go        # Recipe and API response models
│   └── config.go        # Configuration models
├── storage/             # Data persistence layer
│   └── json_storage.go  # JSON file operations
├── static/              # Web interface files
│   ├── index.html       # Main recipe management page
│   ├── login.html       # Login page
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

## Authentication Configuration

### Setup Configuration File

1. **Copy the example configuration**:
   ```bash
   cp config.yaml.example config.yaml
   ```

2. **Edit config.yaml** to customize users and security settings:
   ```yaml
   users:
     - username: "admin"
       password: "admin123"
     - username: "chef"
       password: "cooking456"
   
   # IMPORTANT: Change this secret in production!
   jwt_secret: "your-super-secret-jwt-key-change-this-in-production"
   
   # Token expiration time in hours
   token_expiry_hours: 24
   ```

### Default Demo Users

| Username | Password | Description |
|----------|----------|-------------|
| admin | admin123 | Administrator account |
| chef | cooking456 | Chef account |
| user1 | password123 | Regular user account |

**⚠️ Security Note**: Change all default passwords and JWT secret before production use!

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

3. **Setup configuration**:
   ```bash
   cp config.yaml.example config.yaml
   ```
   Edit `config.yaml` to customize users and JWT secret.

4. **Run the application**:
   ```bash
   go run main.go
   ```

5. **Access the application**:
   - Web Interface: http://localhost:8080 (redirects to login)
   - Login Page: http://localhost:8080/login.html
   - API Endpoints: http://localhost:8080/api/

## Usage

### Web Interface

1. **Login**: Open your browser and go to `http://localhost:8080`
   - You'll be redirected to the login page
   - Use demo credentials: `admin` / `admin123` or `chef` / `cooking456`
   - Click on demo credentials in the login form to auto-fill

2. **Recipe Management**: After successful login:
   - Use the form on the left to add new recipes
   - View all recipes on the right side
   - Click "Edit" to modify a recipe
   - Click "Delete" to remove a recipe
   - Click "Refresh" to reload the recipe list

3. **Logout**: Click the "Logout" button in the header to end your session

### API Usage Examples

#### Login to Get Token
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

Response:
```json
{
  "success": true,
  "message": "Login successful",
  "token": "e8d1f2a3b4c5d6e7f8g9h0i1j2k3l4m5n6o7p8q9r0s1t2u3v4w5x6y7z8a9b0c1d2"
}
```

#### Get All Recipes (Protected)
```bash
curl -X GET http://localhost:8080/api/recipes \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

#### Create a New Recipe (Protected)
```bash
curl -X POST http://localhost:8080/api/recipes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "name": "Chocolate Cake",
    "ingredients": ["2 cups flour", "1 cup sugar", "3 eggs"],
    "instructions": "1. Preheat oven to 350°F\n2. Mix ingredients\n3. Bake for 30 minutes",
    "cooking_time": "45 minutes",
    "servings": 8,
    "category": "dessert"
  }'
```

#### Update a Recipe (Protected)
```bash
curl -X PUT http://localhost:8080/api/recipes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
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

#### Delete a Recipe (Protected)
```bash
curl -X DELETE http://localhost:8080/api/recipes/recipe-uuid-here \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

#### Logout
```bash
curl -X POST http://localhost:8080/api/logout \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
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
