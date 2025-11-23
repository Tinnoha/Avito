package repository

import (
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func NewDatabase() *sqlx.DB {
	connstr := "user=postgres password=0000 dbname=postgres sslmode=disable"

	db, err := sqlx.Connect("pgx", connstr)
	if err != nil {
		fmt.Println("Error: connect to data base:", err)
		return nil
	}

	err = TablesCreate(db)

	if err != nil {
		fmt.Println("Error: init tables in data base:", err)
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
		// Продолжаем выполнение даже если удаление не удалось
	}

	// 1. Сначала создаем teams (без внешних ключей)
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS teams(
        id SERIAL PRIMARY KEY,
        team_name VARCHAR(100) UNIQUE NOT NULL
    )`)
	if err != nil {
		fmt.Println("Error: create teams table:", err)
		return err
	}

	// 2. Затем создаем users (ссылается на teams)
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

	// 3. Затем создаем pull_requests (ссылается на users)
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

	// 4. В конце создаем pull_reviewer (ссылается на users и pull_requests)
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
