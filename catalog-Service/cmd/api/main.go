package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/roman4k-gg/myGarden/catalog-Service/internal/config"
	"github.com/roman4k-gg/myGarden/catalog-Service/internal/storage"
	catalogv1 "github.com/roman4k-gg/myGarden/pkg/catalog_v1"

	"github.com/segmentio/kafka-go"
)

type server struct {
	catalogv1.UnimplementedCatalogServiceServer
	db          *storage.Storage
	kafkaWriter *kafka.Writer
}

func (s *server) GetPlant(ctx context.Context, req *catalogv1.GetPlantRequest) (*catalogv1.GetPlantResponse, error) {
	plant, err := s.db.GetPlant(ctx, req.PlantId)
	if err != nil {
		return nil, err
	}
	return &catalogv1.GetPlantResponse{Plant: plant}, nil
}

func (s *server) ListPlants(ctx context.Context, req *catalogv1.ListPlantsRequest) (*catalogv1.ListPlantsResponse, error) {
	plants, err := s.db.ListPlants(ctx)
	if err != nil {
		return nil, err
	}
	return &catalogv1.ListPlantsResponse{Plants: plants}, nil
}

func (s *server) AddFavorite(ctx context.Context, req *catalogv1.AddFavoriteRequest) (*catalogv1.AddFavoriteResponse, error) {
	err := s.db.AddFavorite(ctx, req.UserId, req.PlantId)
	if err != nil {
		return nil, err
	}

	msgStr := fmt.Sprintf("User %d added plant %d to favorites!", req.UserId, req.PlantId)
	s.kafkaWriter.WriteMessages(ctx, kafka.Message{
		Key:   []byte("favorite_added"),
		Value: []byte(msgStr),
	})

	return &catalogv1.AddFavoriteResponse{Success: true}, nil
}

func (s *server) GetFavorites(ctx context.Context, req *catalogv1.GetFavoritesRequest) (*catalogv1.GetFavoritesResponse, error) {
	plants, err := s.db.GetFavorites(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &catalogv1.GetFavoritesResponse{Favorites: plants}, nil
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	ctx := context.Background()

	db, err := storage.NewStorage(ctx, cfg.DatabaseURL, cfg.RedisAddr)
	if err != nil {
		log.Fatalf("failed to connect to storage: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	kw := &kafka.Writer{
		Addr:     kafka.TCP(cfg.KafkaBroker),
		Topic:    "favorites_notifications",
		Balancer: &kafka.LeastBytes{},
	}
	defer kw.Close()

	s := grpc.NewServer()
	catalogv1.RegisterCatalogServiceServer(s, &server{
		db:          db,
		kafkaWriter: kw,
	})

	log.Printf("Catalog Service listening on %s", cfg.GRPCPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
