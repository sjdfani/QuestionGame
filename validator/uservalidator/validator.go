package uservalidator

import (
	"QuestionGame/dto"
	"QuestionGame/pkg/errmsg"
	"QuestionGame/pkg/richerror"
	"fmt"
	"regexp"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	_ "github.com/go-ozzo/ozzo-validation/v4/is"
)

type Repository interface {
	IsPhoneNumberUnique(phonenumber string) (bool, error)
}

type Validator struct {
	repo Repository
}

func New(repo Repository) Validator {
	return Validator{repo: repo}
}

func (v Validator) ValidateRegisterRequest(req dto.RegisterRequest) (error, map[string]string) {
	const op = "uservalidator.ValidateRegisterRequest"

	if err := validation.ValidateStruct(&req,
		validation.Field(&req.Name, validation.Required, validation.Length(3, 50)),
		validation.Field(&req.Password, validation.Required, validation.Match(regexp.MustCompile("^[a-zA-Z]{10,}$"))),
		validation.Field(
			&req.PhoneNumber,
			validation.Required,
			validation.Match(regexp.MustCompile("^09[0-9]{9}$")),
			validation.By(v.checkPhoneNumberUniqueness),
		),
	); err != nil {

		fieldErrors := make(map[string]string)
		errV, ok := err.(validation.Errors)
		if ok {
			for key, value := range errV {
				fieldErrors[key] = value.Error()
			}
		}

		return richerror.New(op).
			SetMessage(errmsg.ErrorMsgInvalidInput).
			SetKind(richerror.KindInvalid).
			SetError(err).
			SetMeta(map[string]any{"req": req}), fieldErrors
	}

	return nil, nil
}

func (v Validator) checkPhoneNumberUniqueness(value interface{}) error {
	phoneNumber := value.(string)

	if isUnique, err := v.repo.IsPhoneNumberUnique(phoneNumber); err != nil || !isUnique {
		if err != nil {
			return err
		}

		if !isUnique {
			return fmt.Errorf(errmsg.ErrorMsgPhoneNumberIsNotUnique)
		}
	}

	return nil
}
