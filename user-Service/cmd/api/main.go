package main

import (
	"context"
	"log"
	"net"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"github.com/roman4k-gg/myGarden/user-Service/internal/storage"
	userv1 "github.com/roman4k-gg/myGarden/pkg/user_v1"
	"time"
	"github.com/golang-jwt/jwt/v5"
)

type server struct {
	userv1.UnimplementedUserServiceServer
	db *storage.Storage
}

func (s *server) Register(ctx context.Context, req *userv1.RegisterRequest) (*userv1.RegisterResponse, error) {
	log.Printf("Получен запрос на регистрацию! Email: %s, Имя: %s", req.Email, req.Name)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("ошибка при хэшировании пароля: %w", err)
	}
	id, err := s.db.CreateUser(ctx, req.Email, string(hashedPassword), req.Name)
	if err != nil {
		return nil, fmt.Errorf("ошибка при сохранении юзера в БД: %w", err)
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

var jwtSecretKey = []byte("my-super-secret-key-change-it-later")

func (s *server) Login(ctx context.Context, req *userv1.LoginRequest) (*userv1.LoginResponse, error) {
	user, err := s.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("неверный email или пароль")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка при генерации токена")
	}

	return &userv1.LoginResponse{
		AccessToken: tokenString,
	}, nil
}

func main() {

	ctx := context.Background()
	
	connStr := "postgres://user:password@localhost:5432/mygarden_db?sslmode=disable"

	db, err := storage.NewStorage(ctx, connStr)

	if err != nil {
		log.Fatalf("ошибка БД: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	userv1.RegisterUserServiceServer(s, &server{db: db})

	log.Println("User Service запущен на порту :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
