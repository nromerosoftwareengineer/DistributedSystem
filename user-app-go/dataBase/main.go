package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {

	// Connection URL
	databaseURL := "postgresql://noelromero:Lufkin@localhost:5432/noelromero"
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, databaseURL)
	fmt.Fprintf(os.Stderr, "database info: %v\n %v\n", err, os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	rows, err := conn.Query(ctx, `
		SELECT userid, firstName, lastName, city, state, address1, address2, zip 
		FROM users
	`)
	if err != nil {
		log.Fatalf("Error querying database: %v\n", err)
	}
	defer rows.Close()

	// Iterate through the results
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
			log.Fatalf("Error scanning row: %v\n", err)
		}
		users = append(users, user)
	}

	// Check for any errors during iteration
	if err := rows.Err(); err != nil {
		log.Fatalf("Error iterating rows: %v\n", err)
	}

	fmt.Println("Retrieved Users:")
	for _, user := range users {
		fmt.Printf("\nUser ID: %s\n", user.UserID)
		fmt.Printf("Name: %s %s\n", user.FirstName, user.LastName)
		fmt.Printf("Address: %s, %s\n", user.Address1, user.Address2)
		fmt.Printf("Location: %s, %s %s\n", user.City, user.State, user.Zip)
		fmt.Println("------------------------")
	}
}
