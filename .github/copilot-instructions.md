# Copilot Instructions for SSO Web Application

<!-- Use this file to provide workspace-specific custom instructions to Copilot. For more details, visit https://code.visualstudio.com/docs/copilot/copilot-customization#_use-a-githubcopilotinstructionsmd-file -->

## Project Overview
This is a Single Sign-On (SSO) web application built with Go. The project follows clean architecture principles and includes:

- JWT-based authentication
- OAuth2 integration (Google, GitHub)
- User management and registration
- Protected routes with middleware
- Database integration with migrations
- RESTful API design
- Environment-based configuration

## Code Style Guidelines
- Follow Go conventions and best practices
- Use descriptive variable and function names
- Include proper error handling
- Write unit tests for all business logic
- Use dependency injection for better testability
- Follow the repository pattern for data access

## Project Structure
- `cmd/` - Application entry points
- `internal/` - Private application code
- `pkg/` - Public library code
- `configs/` - Configuration files
- `migrations/` - Database migration files
- `static/` - Static web assets
- `templates/` - HTML templates

## Security Considerations
- Always validate input data
- Use secure password hashing (bcrypt)
- Implement proper JWT token validation
- Use HTTPS in production
- Sanitize user inputs to prevent XSS
- Implement rate limiting for authentication endpoints
