package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error
	connStr := "postgres://username:password@localhost:5432/events_db?sslmode=disable"
	DB, err = sql.Open("postgres", connStr)

	if err != nil {
		panic("Could not connect to DB: " + err.Error())
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
	createTables()
}

func createTables() {
	createUserTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL
	)`
	_, err := DB.Exec(createUserTable)

	if err != nil {
		panic("Could not create user table: " + err.Error())
	}

	createEventsTable := `
	CREATE TABLE IF NOT EXISTS events (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		location TEXT NOT NULL,
		datetime TIMESTAMP NOT NULL,
		user_id INTEGER REFERENCES users(id)
	)`

	_, err = DB.Exec(createEventsTable)

	if err != nil {
		panic("Could not create events table: " + err.Error())
	}

	createRegistrationsTable := `
	CREATE TABLE IF NOT EXISTS registrations (
		id SERIAL PRIMARY KEY,
		event_id INTEGER REFERENCES events(id),
		user_id INTEGER REFERENCES users(id)
	)`
	_, err = DB.Exec(createRegistrationsTable)

	if err != nil {
		panic("Could not create registration table: " + err.Error())
	}
}
