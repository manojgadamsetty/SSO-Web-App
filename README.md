# SSO Web Application

A modern, secure Single Sign-On (SSO) web application built with Go, featuring JWT authentication, OAuth2 integration, and a responsive web interface.

## Features

- **JWT Authentication**: Secure token-based authentication
- **OAuth2 Integration**: Sign in with Google and GitHub
- **User Management**: Complete profile management system
- **Security**: Password hashing, secure sessions, and input validation
- **Responsive Design**: Mobile-friendly interface with Bootstrap
- **Database**: SQLite with GORM ORM
- **Modern UI**: Clean, professional interface with Font Awesome icons

## Quick Start

### Prerequisites

- Go 1.21 or later
- Git

### Installation

1. **Clone the repository:**
   ```bash
   git clone <repository-url>
   cd sso-web-app
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   ```

3. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application:**
   ```bash
   go run cmd/server/main.go
   ```

5. **Open your browser:**
   Navigate to `http://localhost:8080`

## Configuration

### Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
# Server Configuration
PORT=8080

# Database Configuration
DATABASE_URL=sso_app.db

# JWT Configuration
JWT_SECRET=your-very-secure-secret-key

# Google OAuth Configuration
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback

# GitHub OAuth Configuration
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
GITHUB_REDIRECT_URL=http://localhost:8080/auth/github/callback
```

### OAuth Setup

#### Google OAuth
1. Go to the [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google+ API
4. Create OAuth 2.0 credentials
5. Add authorized redirect URIs: `http://localhost:8080/auth/google/callback`

#### GitHub OAuth
1. Go to [GitHub Settings > Developer settings > OAuth Apps](https://github.com/settings/developers)
2. Click "New OAuth App"
3. Fill in the application details
4. Set Authorization callback URL: `http://localhost:8080/auth/github/callback`

## Project Structure

```
sso-web-app/
├── cmd/
│   └── server/
│       └── main.go           # Application entry point
├── internal/
│   ├── handlers/
│   │   └── auth.go          # HTTP handlers
│   ├── middleware/
│   │   └── auth.go          # Authentication middleware
│   ├── models/
│   │   └── user.go          # Data models
│   ├── repository/
│   │   └── user.go          # Data access layer
│   └── services/
│       ├── auth.go          # Authentication service
│       └── oauth.go         # OAuth service
├── templates/               # HTML templates
├── static/                  # Static assets
├── configs/                 # Configuration files
├── migrations/              # Database migrations
├── .env.example            # Environment variables template
├── go.mod                  # Go module file
└── README.md               # Project documentation
```

## API Endpoints

### Authentication
- `POST /login` - User login
- `POST /register` - User registration
- `GET /logout` - User logout

### OAuth
- `GET /auth/google` - Initiate Google OAuth
- `GET /auth/google/callback` - Google OAuth callback
- `GET /auth/github` - Initiate GitHub OAuth
- `GET /auth/github/callback` - GitHub OAuth callback

### Protected Routes
- `GET /dashboard` - User dashboard
- `GET /profile` - User profile
- `POST /profile` - Update profile

### API Endpoints
- `GET /api/v1/user` - Get current user
- `PUT /api/v1/user` - Update user

## Development

### Running in Development Mode

```bash
# Install air for hot reloading (optional)
go install github.com/cosmtrek/air@latest

# Run with hot reloading
air

# Or run normally
go run cmd/server/main.go
```

### Building for Production

```bash
# Build the application
go build -o bin/sso-app cmd/server/main.go

# Run the built application
./bin/sso-app
```

### Testing

```bash
# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

## Security Features

- **Password Hashing**: Uses bcrypt for secure password storage
- **JWT Tokens**: Stateless authentication with configurable expiration
- **CSRF Protection**: State parameter validation for OAuth flows
- **Input Validation**: Server-side validation of all user inputs
- **HTTP-Only Cookies**: Secure token storage
- **Rate Limiting**: Protection against brute force attacks (recommended)

## Database Schema

The application uses SQLite with GORM for database operations. The main entities include:

### User Model
```go
type User struct {
    ID          uint
    Email       string
    Password    string
    FirstName   string
    LastName    string
    IsActive    bool
    IsVerified  bool
    GoogleID    string
    GitHubID    string
    AvatarURL   string
    Bio         string
    Website     string
    Location    string
    // ... timestamps
}
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you have any questions or need help getting started, please open an issue in the repository.

## Roadmap

- [ ] Email verification system
- [ ] Password reset functionality
- [ ] Two-factor authentication (2FA)
- [ ] Admin dashboard
- [ ] User roles and permissions
- [ ] API rate limiting
- [ ] Docker containerization
- [ ] Kubernetes deployment manifests
- [ ] Additional OAuth providers (Facebook, Twitter, etc.)
- [ ] Session management improvements
