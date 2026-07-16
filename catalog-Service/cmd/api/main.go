package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	catalogv1 "github.com/roman4k-gg/myGarden/pkg/catalog_v1"
	"github.com/roman4k-gg/myGarden/catalog-Service/internal/storage"
)

type server struct {
	catalogv1.UnimplementedCatalogServiceServer
	db *storage.Storage
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
	ctx := context.Background()
	
	connStr := "postgres://user:password@localhost:5432/mygarden_db?sslmode=disable"
	db, err := storage.NewStorage(ctx, connStr)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	catalogv1.RegisterCatalogServiceServer(s, &server{db: db})

	log.Println("Catalog Service listening on port 50052")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
