package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var LocalPostgres *sql.DB

func ConnectPostgres() {
	connStr := "host=localhost port=5432 user=user password=password dbname=localdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	// Set connection pool settings
	db.SetMaxOpenConns(50)                 // Maximum number of open connections
	db.SetMaxIdleConns(25)                 // Maximum number of idle connections
	db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of connections

	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS test_table (
            id SERIAL PRIMARY KEY,
            data TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}

	LocalPostgres = db
	log.Println("Connected to PostgreSQL")
}

func InsertToPostgres(data string) error {
	_, err := LocalPostgres.Exec("INSERT INTO test_table (data) VALUES ($1)", data)
	return err
}

func ReadFromPostgres(id int) (string, error) {
	var data string
	err := LocalPostgres.QueryRow("SELECT data FROM test_table WHERE id = $1", id).Scan(&data)
	return data, err
}
