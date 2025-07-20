package repository

import (
	"database/sql"
	"errors"
	"os"
	"regexp"
	"strconv"
	"time"

	"jatistore/internal/database"
	"jatistore/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRepository handles database operations for users
type UserRepository struct {
	db *database.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func validatePasswordRules(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters")
	}
	if match, _ := regexp.MatchString(`[0-9]`, password); !match {
		return errors.New("password must contain at least one numeric character")
	}
	if match, _ := regexp.MatchString(`[A-Z]`, password); !match {
		return errors.New("password must contain at least one uppercase letter")
	}
	if match, _ := regexp.MatchString(`[^a-zA-Z0-9]`, password); !match {
		return errors.New("password must contain at least one symbol")
	}
	return nil
}

func getBcryptCost() int {
	costStr := os.Getenv("ROUND")
	if costStr == "" {
		return 12
	}
	cost, err := strconv.Atoi(costStr)
	if err != nil || cost < bcrypt.MinCost || cost > bcrypt.MaxCost {
		return 12
	}
	return cost
}

// CreateUser creates a new user in the database
func (r *UserRepository) CreateUser(user *models.User) error {
	if err := validatePasswordRules(user.Password); err != nil {
		return err
	}
	// Use SALT from env
	salt := os.Getenv("SALT")
	passwordWithSalt := salt + user.Password
	cost := getBcryptCost()
	// Hash the password with bcrypt cost from env
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithSalt), cost)
	if err != nil {
		return err
	}

	user.ID = uuid.New()
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, username, email, password, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err = r.db.Exec(query, user.ID, user.Username, user.Email, user.Password, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt)
	return err
}

// GetUserByID retrieves a user by ID
func (r *UserRepository) GetUserByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password, role, is_active, created_at, updated_at
		FROM users WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// GetUserByUsername retrieves a user by username
func (r *UserRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password, role, is_active, created_at, updated_at
		FROM users WHERE username = $1
	`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// GetUserByEmail retrieves a user by email
func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, username, email, password, role, is_active, created_at, updated_at
		FROM users WHERE email = $1
	`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// GetAllUsers retrieves all users from the database
func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	query := `
		SELECT id, username, email, role, is_active, created_at, updated_at
		FROM users ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email,
			&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// UpdateUser updates a user in the database
func (r *UserRepository) UpdateUser(user *models.User) error {
	user.UpdatedAt = time.Now()
	query := `
		UPDATE users 
		SET username = $1, email = $2, role = $3, is_active = $4, updated_at = $5
		WHERE id = $6
	`

	result, err := r.db.Exec(query, user.Username, user.Email, user.Role, user.IsActive, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// UpdatePassword updates a user's password
func (r *UserRepository) UpdatePassword(userID uuid.UUID, newPassword string) error {
	if err := validatePasswordRules(newPassword); err != nil {
		return err
	}
	// Use SALT from env
	salt := os.Getenv("SALT")
	passwordWithSalt := salt + newPassword
	cost := getBcryptCost()
	// Hash the password with bcrypt cost from env
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordWithSalt), cost)
	if err != nil {
		return err
	}

	query := `UPDATE users SET password = $1, updated_at = $2 WHERE id = $3`
	result, err := r.db.Exec(query, string(hashedPassword), time.Now(), userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// DeleteUser deletes a user from the database
func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// CheckPassword verifies if the provided password matches the user's password
func (r *UserRepository) CheckPassword(user *models.User, password string) bool {
	salt := os.Getenv("SALT")
	passwordWithSalt := salt + password
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwordWithSalt)) == nil
}
