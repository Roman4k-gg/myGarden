package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	catalogv1 "github.com/roman4k-gg/myGarden/pkg/catalog_v1"
)

type Storage struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewStorage(ctx context.Context, connString, redisAddr string) (*Storage, error) {
	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err := db.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %w", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("unable to connect to redis: %w", err)
	}

	log.Println("Catalog Service DB & Redis connected")
	return &Storage{db: db, rdb: rdb}, nil
}

func (s *Storage) Close() {
	s.rdb.Close()
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
	cacheKey := "catalog:plants:all"

	val, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var plants []*catalogv1.Plant
		if err := json.Unmarshal([]byte(val), &plants); err == nil {
			log.Println("Plants loaded from Redis cache")
			return plants, nil
		}
	}

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

	if len(plants) > 0 {
		if data, err := json.Marshal(plants); err == nil {
			s.rdb.Set(ctx, cacheKey, data, 10*time.Minute)
		}
	}

	log.Println("Plants loaded from Postgres")
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
