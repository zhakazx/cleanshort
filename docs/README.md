# CleanShort API Documentation

This folder contains comprehensive API documentation for the CleanShort URL shortener service.

## 📁 Files Overview

### `openapi.yaml`
- **OpenAPI 3.1 Specification** - Complete API specification in YAML format
- Contains all endpoints, request/response schemas, authentication details, and examples
- Can be imported into any OpenAPI-compatible tool (Swagger UI, Postman, Insomnia, etc.)

### `postman_collection.json`
- **Postman Collection** - Ready-to-use collection for API testing
- Includes all endpoints with example requests and automatic token management
- Contains test scripts for authentication flow
- Variables for easy environment switching

### `api-docs.html`
- **Interactive Documentation** - HTML page with embedded Swagger UI
- Provides a user-friendly interface to explore and test the API
- Can be served statically or opened directly in a browser
- Includes quick navigation and downloadable resources

## 🚀 How to Use

### Option 1: View Interactive Documentation
1. Open `api-docs.html` in your web browser
2. The page will load the OpenAPI specification and display it using Swagger UI
3. You can explore endpoints, view schemas, and test API calls directly

### Option 2: Import into Postman
1. Open Postman
2. Click "Import" and select `postman_collection.json`
3. Set up environment variables:
   - `base_url`: Your API base URL (default: `http://localhost:8080`)
   - `user_email`: Test user email
   - `user_password`: Test user password
4. Run the "Register User" and "Login User" requests to get started

### Option 3: Use OpenAPI Specification
1. Import `openapi.yaml` into your preferred API tool:
   - **Swagger Editor**: https://editor.swagger.io/
   - **Insomnia**: File → Import → OpenAPI
   - **VS Code**: Use OpenAPI extensions
   - **Postman**: Import → OpenAPI

## 🔧 Serving Documentation Locally

If you want to serve the documentation as part of your application:

### Option 1: Static File Server
```bash
# Serve the docs folder
cd docs
python -m http.server 8081
# Visit http://localhost:8081/api-docs.html
```

### Option 2: Add to Fiber Application
Add this to your `routes/routes.go`:

```go
// Serve API documentation
app.Static("/docs", "./docs")
app.Get("/docs", func(c *fiber.Ctx) error {
    return c.Redirect("/docs/api-docs.html")
})
```

Then visit: `http://localhost:8080/docs`

## 📋 API Testing Workflow

### Using Postman Collection:

1. **Setup**:
   - Import the collection
   - Set environment variables
   - Ensure your API server is running

2. **Authentication Flow**:
   - Run "Register User" (creates account)
   - Run "Login User" (gets tokens, automatically saved)
   - Tokens are automatically used for subsequent requests

3. **Test Links**:
   - Create links with "Create Link" requests
   - List, update, and delete links
   - Test public redirect functionality

4. **Error Testing**:
   - Try invalid requests to test error handling
   - Test rate limiting by making rapid requests

## 🔍 API Features Documented

- ✅ **Authentication**: Registration, login, token refresh, logout
- ✅ **Links Management**: CRUD operations with filtering and pagination
- ✅ **Public Redirect**: URL redirection with click tracking
- ✅ **Health Checks**: Application health and readiness
- ✅ **Error Handling**: Comprehensive error responses
- ✅ **Rate Limiting**: Request limits and headers
- ✅ **Security**: JWT authentication and validation

## 📖 Additional Resources

- **Main README**: `../README.md` - Project setup and overview
- **Source Code**: `../` - Complete implementation
- **Environment Config**: `../.env.example` - Configuration template

## 🤝 Contributing

When adding new endpoints or modifying existing ones:

1. Update `openapi.yaml` with new specifications
2. Add corresponding requests to `postman_collection.json`
3. Test all changes using the interactive documentation
4. Ensure examples and descriptions are accurate and helpful

## 📞 Support

For questions about the API or documentation:
- Check the interactive documentation for detailed examples
- Review the Postman collection for working requests
- Refer to the main project README for setup instructions