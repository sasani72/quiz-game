package userservice

import (
	"fmt"
	"quiz-game/entity"
	"quiz-game/param"
)

func (s Service) Register(req param.RegisterRequest) (param.RegisterResponse, error) {
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
		return param.RegisterResponse{}, fmt.Errorf("unexpected error: %v", err)
	}
	// return user
	return param.RegisterResponse{param.UserInfo{
		ID:          createdUser.ID,
		Name:        createdUser.Name,
		PhoneNumber: createdUser.PhoneNumber,
	}}, nil
}
