package database

import (
	"context"
	"database/sql"
	"time"
)

type AttendeeModel struct {
	DB *sql.DB
}

type Attendee struct {
	ID      int `json:"id"`
	EventID int `json:"event_id"`
	UserID  int `json:"user_id"`
}

func (s *AttendeeModel) Insert(attendee *Attendee) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO attendees (event_id, user_id)
		VALUES ($1, $2) RETURNING id`

	err := s.DB.QueryRowContext(ctx, query, attendee.EventID, attendee.UserID).Scan(&attendee.ID)
	if err != nil {
		return err
	}
	return nil
}

func (s *AttendeeModel) Get(id int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `SELECT id, event_id, user_id FROM attendees WHERE id = ?`
	row := s.DB.QueryRowContext(ctx, query, id)
	var attendee Attendee
	if err := row.Scan(&attendee.ID, &attendee.EventID, &attendee.UserID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &attendee, nil
}

func (s *AttendeeModel) GetByEventAndUserId(eventID, userID int) (*Attendee, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, event_id, user_id FROM attendees WHERE event_id = ? AND user_id = ?`
	row := s.DB.QueryRowContext(ctx, query, eventID, userID)

	var attendee Attendee
	if err := row.Scan(&attendee.ID, &attendee.EventID, &attendee.UserID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Attendee not found
		}
		return nil, err // Other error
	}
	return &attendee, nil
}

func (s *AttendeeModel) GetAttendeesByEvent(id int) ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT u.id, u.email, u.name FROM users u JOIN attendees a ON a.user_id = u.id WHERE a.event_id = $1`
	rows, err := s.DB.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Email, &user.Username); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *AttendeeModel) Delete(attendeeID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM attendees WHERE id = ?`
	result, err := s.DB.ExecContext(ctx, query, attendeeID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows // No attendee found to delete
	}
	return nil
}
