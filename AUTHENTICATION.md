# JatiStore Authentication System

This document describes the authentication system implemented in JatiStore using JWT (JSON Web Tokens).

## Overview

The authentication system provides secure access to all API endpoints through JWT token-based authentication. All features can only be accessed by authenticated users.

## Features

- **User Registration**: Create new user accounts with role-based access
- **User Login**: Authenticate users and receive JWT tokens
- **Password Security**: Passwords are hashed using bcrypt
- **Role-Based Access Control**: Three user roles (admin, user, cashier)
- **Token Validation**: JWT tokens with 24-hour expiration
- **Protected Routes**: All API endpoints require authentication
- **Admin-Only Routes**: User management endpoints restricted to admin role

## User Roles

- **admin**: Full access to all features including user management
- **user**: Standard access to POS features
- **cashier**: Access to order processing and basic features

## API Endpoints

### Public Endpoints (No Authentication Required)

#### Register User
```
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "string",
  "email": "string",
  "password": "string",
  "role": "admin|user|cashier"
}
```

#### Login
```
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "string",
  "password": "string"
}
```

### Protected Endpoints (Authentication Required)

All endpoints below require a valid JWT token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

#### User Profile Management
- `GET /api/v1/auth/profile` - Get current user profile
- `PUT /api/v1/auth/profile` - Update current user profile
- `POST /api/v1/auth/change-password` - Change current user password

#### Admin-Only Endpoints
- `GET /api/v1/auth/users` - Get all users
- `GET /api/v1/auth/users/{id}` - Get user by ID
- `PUT /api/v1/auth/users/{id}` - Update user
- `DELETE /api/v1/auth/users/{id}` - Delete user

#### All Other API Endpoints
All existing endpoints (products, categories, inventory, customers, orders) now require authentication.

## Environment Configuration

Add the following to your `.env` file:

```env
# JWT Configuration
JWT_SECRET=your-secret-key-here
```

**Important**: Use a strong, unique secret key in production. The default key is only for development.

## Database Schema

The authentication system adds a `users` table to the database:

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user', 'cashier')),
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

## Security Features

1. **Password Hashing**: All passwords are hashed using bcrypt with default cost
2. **JWT Tokens**: Secure token-based authentication with 24-hour expiration
3. **Role Validation**: Server-side role validation for all protected routes
4. **Input Validation**: Comprehensive validation for all user inputs
5. **Account Status**: Users can be deactivated without deletion
6. **Unique Constraints**: Username and email must be unique

## 🔑 Password Policy

User passwords must meet the following requirements:
- Minimum 8 characters
- At least 1 numeric character
- At least 1 symbol
- At least 1 uppercase letter

If the password does not meet these requirements, registration or password change will fail.

## 🔒 Password Hashing Configuration

- **SALT**: A secret string from the environment, prepended to the password before hashing. Set `SALT` in your `.env` file.
- **ROUND**: Bcrypt cost (number of hashing rounds, default: 12). Set `ROUND` in your `.env` file.
- Passwords are hashed as: `bcrypt(SALT + password, cost=ROUND)`

Example `.env`:
```env
SALT=your-random-salt-string
ROUND=12
```

## Usage Examples

### 1. Register a New User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@jatistore.com",
    "password": "admin123",
    "role": "admin"
  }'
```

### 2. Login and Get Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

### 3. Access Protected Endpoint
```bash
curl -X GET http://localhost:8080/api/v1/products \
  -H "Authorization: Bearer <your_jwt_token>" \
  -H "Content-Type: application/json"
```

### 4. Get User Profile
```bash
curl -X GET http://localhost:8080/api/v1/auth/profile \
  -H "Authorization: Bearer <your_jwt_token>" \
  -H "Content-Type: application/json"
```

## Testing

Use the provided test script to verify the authentication system:

```bash
./test_auth.sh
```

This script will:
1. Register a new admin user
2. Login and get a JWT token
3. Test protected routes
4. Verify unauthorized access is blocked
5. Confirm authorized access works

## Error Responses

The API returns consistent error responses:

```json
{
  "success": false,
  "error": "Error message"
}
```

Common error scenarios:
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: Insufficient permissions for the requested role
- `400 Bad Request`: Invalid input data
- `409 Conflict`: Username or email already exists

## Migration from Previous Version

If you're upgrading from a version without authentication:

1. The database migration will automatically create the users table
2. All existing endpoints now require authentication
3. You'll need to register at least one admin user to access the system
4. Update your client applications to include JWT tokens in requests

## Best Practices

1. **Token Storage**: Store JWT tokens securely (e.g., in HTTP-only cookies or secure storage)
2. **Token Refresh**: Implement token refresh logic for long-running sessions
3. **Password Policy**: Enforce strong password requirements in your client application
4. **Rate Limiting**: Consider implementing rate limiting for login attempts
5. **Logging**: Monitor authentication attempts and failures
6. **HTTPS**: Always use HTTPS in production environments

## Troubleshooting

### Common Issues

1. **"Authorization header is required"**: Make sure to include the Authorization header with the Bearer token
2. **"Invalid or expired token"**: Token has expired or is malformed. Re-login to get a new token
3. **"Insufficient permissions"**: User role doesn't have access to the requested endpoint
4. **"Account is deactivated"**: User account has been deactivated by an admin

### Debug Mode

To enable debug logging, set the environment variable:
```env
LOG_LEVEL=debug
``` 