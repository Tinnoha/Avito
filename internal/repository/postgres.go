package repository

import (
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewDatabase() *sqlx.DB {

	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("DB_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		password = "0000"
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "postgres"
	}

	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sqlx.Connect("pgx", connstr)
	if err != nil {
		fmt.Println("Error connect to db:", err)
		return nil
	}

	err = TablesCreate(db)
	if err != nil {
		fmt.Println("Error create tables:", err)
		return nil
	}

	return db
}

func TablesCreate(db *sqlx.DB) error {
	err := db.Ping()
	if err != nil {
		fmt.Println("Error: connection to data base:", err)
		return err
	}
	_, err = db.Exec(`DROP TABLE IF EXISTS pull_reviewer, pull_requests, users, teams CASCADE`)
	if err != nil {
		fmt.Println("Warning: drop tables:", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS teams(
        id SERIAL PRIMARY KEY,
        team_name VARCHAR(100) UNIQUE NOT NULL
    )`)
	if err != nil {
		fmt.Println("Error: create teams table:", err)
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(
        user_id VARCHAR(100) PRIMARY KEY,
        username VARCHAR(100) UNIQUE NOT NULL,
        team_id INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
        is_active BOOLEAN DEFAULT true
    )`)
	if err != nil {
		fmt.Println("Error: create users table:", err)
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS pull_requests(
        pull_request_id VARCHAR(100) PRIMARY KEY, 
        pull_request_name VARCHAR(100) NOT NULL,
        author_id VARCHAR(100) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
        status VARCHAR(7) DEFAULT 'open',
        created_at TIMESTAMPTZ DEFAULT NOW(),
        merged_at TIMESTAMPTZ DEFAULT NULL
    )`)
	if err != nil {
		fmt.Println("Error: create requests table:", err)
		return err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS pull_reviewer(
        pull_request_id VARCHAR(100) NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
        reviewer_id VARCHAR(100) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
        PRIMARY KEY (pull_request_id, reviewer_id)
    )`)
	if err != nil {
		fmt.Println("Error: create reviewer table:", err)
		return err
	}

	fmt.Println("All tables created successfully")
	return nil
}
