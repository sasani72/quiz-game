package uservalidator

import (
	"fmt"
	"quiz-game/dto"
	"quiz-game/pkg/errmsg"
	"quiz-game/pkg/richerror"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateLoginRequest(req dto.LoginRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateLoginRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Password, validation.Required),
		validation.Field(&req.PhoneNumber, validation.Required,
			validation.Match(regexp.MustCompile(phoneNumberRegex)).Error(errmsg.InvalidPhoneNumber),
			validation.By(v.doesPhoneNumberExist)),
	); err != nil {

		fieldErrors := make(map[string]string)
		errV, ok := err.(validation.Errors)
		if ok {
			for key, value := range errV {
				if value != nil {
					fieldErrors[key] = value.Error()
				}
			}
		}

		return fieldErrors, richerror.New(op).WithMessage(errmsg.InvalidInput).
			WithKind(richerror.KindInvalid).
			WithErr(err).
			WithMeta(map[string]interface{}{"req": req})
	}

	return nil, nil
}

func (v Validator) doesPhoneNumberExist(value interface{}) error {
	const op = "uservalidator.doesPhoneNumberExist"
	phoneNumber := value.(string)
	_, err := v.repo.GetUserByPhoneNumber(phoneNumber)
	if err != nil {
		return fmt.Errorf(errmsg.NotFound)
	}

	return nil
}
