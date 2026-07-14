package storage

import (
	"context"
	"fmt"
	"log"
	"github.com/jackc/pgx/v5/pgxpool"
)


type Storage struct {
	db *pgxpool.Pool
}

type User struct {
	ID int
	Email string
	PasswordHash string
	Name string
}

func NewStorage(ctx context.Context, connString string) (*Storage, error) {

	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Connected to database!")

	s := &Storage{db: db} 

	return s, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) CreateUser(ctx context.Context, email, passwordHash, name string) (int, error) {
	var id int
	query := `
		INSERT INTO users (email, password_hash, name) 
		VALUES ($1, $2, $3) 
		RETURNING id;
	`
	err := s.db.QueryRow(ctx, query, email, passwordHash, name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %w", err)
	}
	return id, nil
}

func (s *Storage) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `SELECT id, email, password_hash, name FROM users WHERE email = $1`
	err := s.db.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.Name)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}