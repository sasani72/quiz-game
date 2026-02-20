package uservalidator

import (
	"quiz-game/dto"
	"quiz-game/pkg/errmsg"
	"quiz-game/pkg/richerror"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

func (v Validator) ValidateRegisterRequest(req dto.RegisterRequest) (map[string]string, error) {
	const op = "uservalidator.ValidateRegisterRequest"

	if err := validation.ValidateStruct(&req,
		// TODO - add length limit to config
		validation.Field(&req.Name, validation.Required,
			validation.Length(3, 50)),
		validation.Field(&req.Password, validation.Required,
			validation.Match(regexp.MustCompile("^[a-zA-Z0-9!@#%^*]{10,}$"))),
		validation.Field(&req.PhoneNumber, validation.Required,
			validation.Match(regexp.MustCompile(phoneNumberRegex)).Error(errmsg.InvalidPhoneNumber),
			validation.By(v.checkPhoneNumberUniqueness)),
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
