package database

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"go_proj/database/entities"
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

// StoreMessage inserts a new message log
func (s *DBService) InsertMessage(ctx context.Context, msg entities.MessageInput) (int64, error) {
	query := `
        INSERT INTO messages 
        (from_user, to_user, body, message_type, group_id, group_name, group_members)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id`

	var messageID int64
	err := s.pool.QueryRow(
		ctx,
		query,
		msg.FromUser,
		msg.ToUser,
		msg.Body,
		msg.MessageType,
		msg.GroupID,
		msg.GroupName,
		pq.Array(msg.GroupMembers), // Using pq.Array for the string array
	).Scan(&messageID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert message: %v", err)
	}

	return messageID, nil
}
