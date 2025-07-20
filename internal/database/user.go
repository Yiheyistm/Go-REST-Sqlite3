package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

func (s *UserModel) GenerateToken(id int, appJwtSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": id,
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte(appJwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *UserModel) Insert(user *User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `INSERT INTO users(email, password, name) VALUES ($1, $2, $3) RETURNING id`
	return s.DB.QueryRowContext(ctx, query, user.Email, user.Password, user.Username).Scan(&user.ID)
}

func (s *UserModel) Get(id int) (*User, error) {
	query := `SELECT id, name, email, password FROM users WHERE id = $1`
	return s.getUser(query, id)
}

func (s *UserModel) GetByEmail(email string) (*User, error) {
	query := `SELECT id, name, email, password FROM users WHERE email = $1`
	return s.getUser(query, email)
}

func (s *UserModel) getUser(query string, args ...any) (*User, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	row := s.DB.QueryRowContext(ctx, query, args...)
	var user User
	if err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	return &user, nil
}
