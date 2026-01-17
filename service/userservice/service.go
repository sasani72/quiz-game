package userservice

import (
	"fmt"
	"quiz-game/entity"
	"quiz-game/pkg/phonenumber"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
}
type Service struct {
	Repository Repository
}

type RegisterRequest struct {
	Name        string
	PhoneNumber string
}

type RegisterResponse struct {
	User entity.User
}

func (s Service) Register(req RegisterRequest) (RegisterResponse, error) {
	// TODO - we should verify phone number by verification code
	// validate phone num
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("invalid phone number")
	}

	// check unique phone
	if isUnique, err := s.Repository.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
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
	// create new user in storage
	createdUser, err := s.Repository.Register(entity.User{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	})
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %v", err)
	}
	// return user
	return RegisterResponse{createdUser}, nil
}
