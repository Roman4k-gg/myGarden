package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/roman4k-gg/myGarden/user-Service/internal/config"
	"github.com/roman4k-gg/myGarden/user-Service/internal/storage"
	userv1 "github.com/roman4k-gg/myGarden/pkg/user_v1"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
)

type server struct {
	userv1.UnimplementedUserServiceServer
	db        *storage.Storage
	jwtSecret []byte
}

func (s *server) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	id, err := s.db.CreateUser(ctx, req.Email, string(hashedPassword), req.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}
	return &userv1.RegisterResponse{
		UserId: fmt.Sprintf("%d", id),
	}, nil
}

func (s *server) GetProfile(ctx context.Context, req *userv1.GetProfileRequest) (*userv1.GetProfileResponse, error) {
	return &userv1.GetProfileResponse{
		UserId: 1,
		Email:  "test@example.com",
		Name:   "Test User",
	}, nil
}

func (s *server) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	user, err := s.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token")
	}

	return &userv1.LoginResponse{
		AccessToken: tokenString,
	}, nil
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	ctx := context.Background()

	db, err := storage.NewStorage(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	userv1.RegisterUserServiceServer(s, &server{
		db:        db,
		jwtSecret: []byte(cfg.JWTSecret),
	})

	log.Printf("User Service listening on %s", cfg.GRPCPort)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
