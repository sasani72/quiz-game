package uservalidator

import (
	"fmt"
	"quiz-game/entity"
	"quiz-game/pkg/errmsg"
)

const (
	phoneNumberRegex = "^09[0-9]{9}$"
)

type Repository interface {
	IsPhoneNumberUnique(phoneNumber string) (bool, error)
	GetUserByPhoneNumber(phoneNumber string) (entity.User, error)
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{
		repo: repo,
	}
}

func (v Validator) checkPhoneNumberUniqueness(value interface{}) error {
	const op = "uservalidator.checkPhoneNumberUniqueness"
	phoneNumber := value.(string)
	if isUnique, err := v.repo.IsPhoneNumberUnique(phoneNumber); err != nil || !isUnique {
		if err != nil {
			return err
		}
		if !isUnique {
			return fmt.Errorf(errmsg.PhoneNumberNotUnique)
		}
	}

	return nil
}
