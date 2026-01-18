package userservice

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"quiz-game/entity"
	"quiz-game/pkg/phonenumber"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
}
type Service struct {
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

func New(repo Repository) Service {
	return Service{repo: repo}
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
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	// check existence of phone num in db
	// get user by phone num
	// compare user pass with req.pass
	// return ok

	panic("implement me")
}

func getMD5Hash(text string) string {
	hasher := md5.Sum([]byte(text))
	return hex.EncodeToString(hasher[:])
}
