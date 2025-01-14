package db

import (
	"context"
	"fmt"
	"os"
	_ "os"

	"github.com/jackc/pgx/v4/pgxpool"
)

// User represents the user data structure
type User struct {
	UserID    string
	FirstName string
	LastName  string
	City      string
	State     string
	Address1  string
	Address2  string
	Zip       string
}

// DBService provides database operations
type DBService struct {
	pool *pgxpool.Pool
}

// NewDBService creates a new database service instance
func NewDBService() (*DBService, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	return &DBService{pool: pool}, nil
}

// Close closes the database connection pool
func (s *DBService) Close() {
	s.pool.Close()
}

// CreateUser inserts a new user record and returns a success message
func (s *DBService) CreateUser(ctx context.Context, user User) (string, error) {
	query := `
        INSERT INTO users (userid, firstName, lastName, city, state, address1, address2, zip)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	result, err := s.pool.Exec(ctx, query,
		user.UserID,
		user.FirstName,
		user.LastName,
		user.City,
		user.State,
		user.Address1,
		user.Address2,
		user.Zip)

	if err != nil {
		return "", fmt.Errorf("error creating user: %v", err)
	}

	// Get number of rows affected
	rowsAffected := result.RowsAffected()
	message := fmt.Sprintf("Successfully created user with ID %s. Rows affected: %d", user.UserID, rowsAffected)

	return message, nil
}

// GetUser retrieves a user by ID
func (s *DBService) GetUser(ctx context.Context, userID string) (*User, error) {
	query := `
		SELECT userid, firstName, lastName, city, state, address1, address2, zip
		FROM users
		WHERE userid = $1`

	var user User
	err := s.pool.QueryRow(ctx, query, userID).Scan(
		&user.UserID,
		&user.FirstName,
		&user.LastName,
		&user.City,
		&user.State,
		&user.Address1,
		&user.Address2,
		&user.Zip,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	return &user, nil
}

// UpdateUser updates an existing user record
func (s *DBService) UpdateUser(ctx context.Context, user User) (string, error) {
	query := `
		UPDATE users
		SET firstName = $2, lastName = $3, city = $4, state = $5, 
			address1 = $6, address2 = $7, zip = $8
		WHERE userid = $1`

	result, err := s.pool.Exec(ctx, query,
		user.UserID,
		user.FirstName,
		user.LastName,
		user.City,
		user.State,
		user.Address1,
		user.Address2,
		user.Zip)

	if err != nil {
		return "", fmt.Errorf("error creating user: %v", err)
	}

	// Get number of rows affected
	rowsAffected := result.RowsAffected()
	message := fmt.Sprintf("Successfully created user with ID %s. Rows affected: %d", user.UserID, rowsAffected)

	return message, nil
}

// DeleteUser deletes a user by ID
func (s *DBService) DeleteUser(ctx context.Context, userID string) error {
	query := `DELETE FROM users WHERE userid = $1`

	_, err := s.pool.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %v", err)
	}

	return nil
}

// ListUsers retrieves all users
func (s *DBService) ListUsers(ctx context.Context) ([]User, error) {
	query := `
		SELECT userid, firstName, lastName, city, state, address1, address2, zip
		FROM users`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %v", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.UserID,
			&user.FirstName,
			&user.LastName,
			&user.City,
			&user.State,
			&user.Address1,
			&user.Address2,
			&user.Zip,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %v", err)
	}

	return users, nil
}
