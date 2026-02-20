package userservice

import (
	"fmt"
	"quiz-game/param"
	"quiz-game/pkg/richerror"
)

func (s Service) Login(req param.LoginRequest) (param.LoginResponse, error) {
	const op = "userservice.Login"
	// TODO - maybe its better to use two separate methods for user existence and get user by phone
	user, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return param.LoginResponse{}, richerror.New(op).WithErr(err).
			WithMeta(map[string]interface{}{"phone_number": req.PhoneNumber})
	}

	if user.Password != getMD5Hash(req.Password) {
		return param.LoginResponse{}, fmt.Errorf("username or password is invalid")
	}

	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return param.LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	return param.LoginResponse{
		User: param.UserInfo{
			ID:          user.ID,
			PhoneNumber: user.PhoneNumber,
			Name:        user.Name,
		},
		Tokens: param.Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}
