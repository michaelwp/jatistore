package services

import (
	"errors"
	"os"
	"time"

	"jatistore/internal/models"
	"jatistore/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// UserService handles business logic for user operations
type UserService struct {
	userRepo *repository.UserRepository
}

// NewUserService creates a new UserService instance
func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// Register creates a new user account
func (s *UserService) Register(req *models.RegisterRequest) (*models.User, error) {
	// Check if username already exists
	existingUser, _ := s.userRepo.GetUserByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, _ = s.userRepo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	user := &models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
		IsActive: true,
	}

	err := s.userRepo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	// Don't return the password
	user.Password = ""
	return user, nil
}

// Login authenticates a user and returns a JWT token
func (s *UserService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	// Get user by username
	user, err := s.userRepo.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Verify password
	if !s.userRepo.CheckPassword(user, req.Password) {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := s.generateJWTToken(user)
	if err != nil {
		return nil, err
	}

	// Don't return the password
	user.Password = ""

	return &models.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Don't return the password
	user.Password = ""
	return user, nil
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAllUsers()
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(id uuid.UUID, req *models.UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Check if username is being changed and if it already exists
	if req.Username != user.Username {
		existingUser, _ := s.userRepo.GetUserByUsername(req.Username)
		if existingUser != nil {
			return nil, errors.New("username already exists")
		}
	}

	// Check if email is being changed and if it already exists
	if req.Email != user.Email {
		existingUser, _ := s.userRepo.GetUserByEmail(req.Email)
		if existingUser != nil {
			return nil, errors.New("email already exists")
		}
	}

	user.Username = req.Username
	user.Email = req.Email
	user.Role = req.Role
	user.IsActive = req.IsActive

	err = s.userRepo.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	// Don't return the password
	user.Password = ""
	return user, nil
}

// ChangePassword changes a user's password
func (s *UserService) ChangePassword(userID uuid.UUID, req *models.ChangePasswordRequest) error {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return err
	}

	// Verify current password
	if !s.userRepo.CheckPassword(user, req.CurrentPassword) {
		return errors.New("current password is incorrect")
	}

	// Update password
	return s.userRepo.UpdatePassword(userID, req.NewPassword)
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(id uuid.UUID) error {
	return s.userRepo.DeleteUser(id)
}

// generateJWTToken generates a JWT token for the user
func (s *UserService) generateJWTToken(user *models.User) (string, error) {
	// Get JWT secret from environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key" // Default secret for development
	}

	// Create claims
	claims := &models.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *UserService) ValidateToken(tokenString string) (*models.Claims, error) {
	// Get JWT secret from environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key" // Default secret for development
	}

	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token
	if claims, ok := token.Claims.(*models.Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
