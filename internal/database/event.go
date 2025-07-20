package database

import (
	"context"
	"database/sql"
	"time"
)

type EventModel struct {
	DB *sql.DB
}

type Event struct {
	ID          int    `json:"id"`
	OwnerId     int    `json:"owner_id"`
	Name        string `json:"name" binding:"required,min=3,max=100"`
	Description string `json:"description" binding:"required,min=10,max=500"`
	Date        string `json:"date" binding:"required,datetime=2006-01-02"`
	Location    string `json:"location" binding:"required,min=3,max=100"`
}

func (s *EventModel) Insert(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO events (owner_id, name, description, date, location)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	return s.DB.QueryRowContext(ctx, query,
		event.OwnerId,
		event.Name,
		event.Description,
		event.Date,
		event.Location,
	).Scan(&event.ID)
}

func (s *EventModel) GetAll() ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "SELECT * FROM events"
	row, err := s.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	events := []*Event{}
	for row.Next() {
		var event Event
		err = row.Scan(&event.ID, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	if err := row.Err(); err != nil {
		return nil, err
	}
	return events, nil
}

func (s *EventModel) GetByID(id int) (*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, owner_id, name, description, date, location FROM events WHERE id = $1`

	event := Event{}
	err := s.DB.QueryRowContext(ctx, query, id).Scan(&event.ID, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (s *EventModel) Update(event *Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := "UPDATE events SET  name = $1, description = $2, date = $3, location = $4 WHERE id = $5"

	_, err := s.DB.ExecContext(ctx, query, event.Name, event.Description, event.Date, event.Location, event.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *EventModel) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM events WHERE id = $1`
	res, err := s.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (s *EventModel) GetByAttendeeId(attendeeId int) ([]*Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT e.id, e.owner_id, e.name, e.description, e.date, e.location FROM events e JOIN attendees a ON a.event_id = e.id WHERE a.id = $1`
	rows, err := s.DB.QueryContext(ctx, query, attendeeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*Event
	for rows.Next() {
		var event Event
		if err := rows.Scan(&event.ID, &event.OwnerId, &event.Name, &event.Description, &event.Date, &event.Location); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return events, nil
}
