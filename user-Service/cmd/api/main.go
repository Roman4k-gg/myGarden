package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"

	userv1 "github.com/roman4k-gg/myGarden/pkg/user_v1"
)

type server struct {
	userv1.UnimplementedUserServiceServer
}

func (s *server) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	log.Printf("Получен запрос на регистрацию! Email: %s, Имя: %s", req.Email, req.Name)
	
	return &userv1.RegisterResponse{
		UserId: "1", 
	}, nil
}

func (s *server) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	return &userv1.LoginResponse{AccessToken: "mock-jwt-token"}, nil
}

func (s *server) GetProfile(ctx context.Context, req *userv1.GetProfileRequest) (*userv1.GetProfileResponse, error) {
	return &userv1.GetProfileResponse{
		UserId: 1,
		Email:  "test@example.com",
		Name:   "Test User",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	userv1.RegisterUserServiceServer(s, &server{})

	log.Println("User Service запущен на порту :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
