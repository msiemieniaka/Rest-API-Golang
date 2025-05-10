package models

import (
	"errors"
	"fmt"
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
	query := `
		INSERT INTO events(name, description, location, datetime, user_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	err := db.DB.QueryRow(query, e.Name, e.Description, e.Location, e.DateTime, e.UserID).Scan(&e.ID)
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
	query := "SELECT * FROM events WHERE id = $1"
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
    UPDATE events 
    SET name = $1, description = $2, location = $3, datetime = $4
    WHERE id = $5 AND user_id = $6
    RETURNING id`

	result := db.DB.QueryRow(query,
		event.Name,
		event.Description,
		event.Location,
		event.DateTime,
		event.ID,
		event.UserID)

	return result.Scan(&event.ID)
}

func (event Event) Delete() error {
	query := `
        DELETE FROM events 
        WHERE id = $1 AND user_id = $2 
        RETURNING id`

	result := db.DB.QueryRow(query, event.ID, event.UserID)

	var id int64
	if err := result.Scan(&id); err != nil {
		return errors.New("event not found or not authorized to delete")
	}
	return nil
}

func (e Event) IsUserRegistered(userID int64) (bool, error) {
	// More explicit query
	query := `
        SELECT COUNT(*) 
        FROM registrations 
        WHERE event_id = $1 
        AND user_id = $2
        AND EXISTS (SELECT 1 FROM events WHERE id = $1)`

	var count int
	err := db.DB.QueryRow(query, e.ID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("database error checking registration: %v", err)
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
	// First check if the event exists
	query := "SELECT id FROM events WHERE id = $1"
	var eventID int64
	err := db.DB.QueryRow(query, event.ID).Scan(&eventID)
	if err != nil {
		return errors.New("event not found")
	}

	// Then try to delete the registration
	query = "DELETE FROM registrations WHERE event_id = $1 AND user_id = $2"
	result, err := db.DB.Exec(query, event.ID, userID)
	if err != nil {
		return fmt.Errorf("error deleting registration: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return errors.New("no active registration found")
	}

	return nil
}
