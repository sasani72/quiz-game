package userservice

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"quiz-game/entity"
	"quiz-game/pkg/phonenumber"

	//"golang.org/x/crypto/bcrypt"
	jwt "github.com/golang-jwt/jwt/v4"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, bool, error)
	GetUserByID(userID uint) (entity.User, error)
}

type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
	CreateRefreshToken(user entity.User) (string, error)
}
type Service struct {
	auth AuthGenerator
	repo Repository
}

type RegisterRequest struct {
	Name        string `json:"name"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type RegisterResponse struct {
	User entity.User
}

func New(auth AuthGenerator, repo Repository) Service {
	return Service{auth: auth, repo: repo}
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO - we should verify phone number by verification code
	// validate phone num
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("invalid phone number")
	}

	// check unique phone
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return RegisterResponse{}, fmt.Errorf("unexpected error: %v", err)
		}
		if !isUnique {
			return RegisterResponse{}, fmt.Errorf("phone number %s is already used", req.PhoneNumber)
		}
	}
	// validate name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name should be at least 3 characters long")
	}

	// TODO - check pass with regex pattern
	//validate password
	if len(req.Password) < 8 {
		return RegisterResponse{}, fmt.Errorf("password should be at least 8 characters long")
	}

	// TODO - replace md5 with bcrypt
	//pass := []byte(req.Password)
	//bcrypt.GenerateFromPassword(pass, 0)
	// create new user in storage
	createdUser, err := s.repo.Register(entity.User{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    getMD5Hash(req.Password),
	})
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %v", err)
	}
	// return user
	return RegisterResponse{createdUser}, nil
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	// TODO - maybe its better to use two separate methods for user existence and get user by phone
	user, exist, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	if !exist || user.Password != getMD5Hash(req.Password) {
		return LoginResponse{}, fmt.Errorf("username or password is invalid")
	}

	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	return LoginResponse{accessToken, refreshToken}, nil
}

type ProfileRequest struct {
	UserID uint
}

type ProfileResponse struct {
	Name string `json:"name"`
}

// all request inputs for interactor/service should be sanitized.
func (s Service) Profile(req ProfileRequest) (ProfileResponse, error) {
	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		// TODO - we need rich error
		return ProfileResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	return ProfileResponse{user.Name}, nil
}

func getMD5Hash(text string) string {
	hasher := md5.Sum([]byte(text))
	return hex.EncodeToString(hasher[:])
}
