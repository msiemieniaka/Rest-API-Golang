package models

import (
	"errors"
	"rest-api/app/db"
	"time"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int64
}

var events = []Event{}

func (e *Event) Save() error {
	query :=
		`INSERT INTO events(name, description, location, dateTime, user_id)
     VALUES (?,?,?,?,?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	e.ID = id
	return err
}

func GetAllEvent() ([]Event, error) {
	query := "SELECT * FROM events"
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := "SELECT * FROM events WHERE id = ?"
	row := db.DB.QueryRow(query, id)

	var event Event
	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &event.DateTime, &event.UserID)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func (event Event) Update() error {
	query := `
	UPDATE EVENTS
SET name = ?, description = ?, location = ?, dateTime = ?
WHERE id = ?
`
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
	return err
}

func (event Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID)
	return err
}

// Add this new method to check if user is already registered
func (e Event) IsUserRegistered(userID int64) (bool, error) {
	query := "SELECT COUNT(*) FROM registrations WHERE event_id = $1 AND user_id = $2"
	var count int
	err := db.DB.QueryRow(query, e.ID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Modify the Register method to check for existing registration
func (e Event) Register(userID int64) error {
	// First check if user is already registered
	isRegistered, err := e.IsUserRegistered(userID)
	if err != nil {
		return err
	}
	if isRegistered {
		return errors.New("user is already registered for this event")
	}

	// If not registered, proceed with registration
	query := "INSERT INTO registrations(event_id, user_id) VALUES ($1,$2)"
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(e.ID, userID)
	return err
}

// Add method to get all registered users for an event
func (e Event) GetRegisteredUsers() ([]int64, error) {
	query := "SELECT user_id FROM registrations WHERE event_id = $1"
	rows, err := db.DB.Query(query, e.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []int64
	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}
	return userIDs, nil
}

func (event Event) CancelRegistration(userID int64) error {
	query := "DELETE FROM registrations WHERE event_id = ? AND user_id = ?"
	stmt, err := db.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID, userID)
	return err

}
