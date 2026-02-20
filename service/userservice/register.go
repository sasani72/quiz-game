package userservice

import (
	"fmt"
	"quiz-game/dto"
	"quiz-game/entity"
)

func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
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
		return dto.RegisterResponse{}, fmt.Errorf("unexpected error: %v", err)
	}
	// return user
	return dto.RegisterResponse{dto.UserInfo{
		ID:          createdUser.ID,
		Name:        createdUser.Name,
		PhoneNumber: createdUser.PhoneNumber,
	}}, nil
}
