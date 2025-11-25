package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (DB *sql.DB) {
	DB, err := sql.Open("sqlite3", "./internal/database/app.db") // Open a connection to the SQLite database file named app.db
	if err != nil {
		log.Fatal(err) // Log an error and stop the program if the database can't be opened
	}

	// SQL statement to create the users table if it doesn't exist
	usersTable := `
		CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		email TEXT NOT NULL UNIQUE,
		username TEXT NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL
 );`

	_, err = DB.Exec(usersTable)
	if err != nil {
		log.Fatalf("Error creating table: %q: %s\n", err, usersTable) // Log an error if table creation fails
	}

	return DB
}
