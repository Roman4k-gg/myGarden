package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	catalogv1 "github.com/roman4k-gg/myGarden/pkg/catalog_v1"
)

type Storage struct {
	db *pgxpool.Pool
}

func NewStorage(ctx context.Context, connString string) (*Storage, error) {
	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	log.Println("Catalog Service DB connected")
	return &Storage{db: db}, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) GetPlant(ctx context.Context, plantID int32) (*catalogv1.Plant, error) {
	query := `SELECT id, name, description, watering_interval_days FROM plants WHERE id = $1`
	p := &catalogv1.Plant{}
	err := s.db.QueryRow(ctx, query, plantID).Scan(&p.Id, &p.Name, &p.Description, &p.WateringIntervalDays)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Storage) ListPlants(ctx context.Context) ([]*catalogv1.Plant, error) {
	query := `SELECT id, name, description, watering_interval_days FROM plants`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plants []*catalogv1.Plant
	for rows.Next() {
		p := &catalogv1.Plant{}
		err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.WateringIntervalDays)
		if err != nil {
			return nil, err
		}
		plants = append(plants, p)
	}
	return plants, nil
}

func (s *Storage) AddFavorite(ctx context.Context, userID, plantID int32) error {
	query := `INSERT INTO user_favorites (user_id, plant_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`
	_, err := s.db.Exec(ctx, query, userID, plantID)
	return err
}

func (s *Storage) GetFavorites(ctx context.Context, userID int32) ([]*catalogv1.Plant, error) {
	query := `
		SELECT p.id, p.name, p.description, p.watering_interval_days 
		FROM plants p
		JOIN user_favorites uf ON p.id = uf.plant_id
		WHERE uf.user_id = $1
	`
	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plants []*catalogv1.Plant
	for rows.Next() {
		p := &catalogv1.Plant{}
		err := rows.Scan(&p.Id, &p.Name, &p.Description, &p.WateringIntervalDays)
		if err != nil {
			return nil, err
		}
		plants = append(plants, p)
	}
	return plants, nil
}
