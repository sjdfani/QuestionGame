package userservice

import (
	"QuestionGame/entity"
	"QuestionGame/pkg/phonenumber"

	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	IsPhoneNumberUnique(phonenumber string) (bool, error)
	Register(u entity.User) (entity.User, error)
	GetUserByPhoneNumber(phonenumber string) (entity.User, bool, error)
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

	// validate phone number
	if !phonenumber.IsValid(req.PhoneNumber) {
		return RegisterResponse{}, fmt.Errorf("phone number is not valid")
	}

	// check uniqueness phone number
	if isUnique, err := s.repo.IsPhoneNumberUnique(req.PhoneNumber); err != nil || !isUnique {
		if err != nil {
			return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
		}

		if !isUnique {
			return RegisterResponse{}, fmt.Errorf("phone number is not unique")
		}
	}

	// validate name
	if len(req.Name) < 3 {
		return RegisterResponse{}, fmt.Errorf("name length should be greater than 3")
	}

	// TODO: Check the password with regex pattern
	// validate password
	if len(req.Password) < 8 {
		return RegisterResponse{}, fmt.Errorf("password length should be greater than 8")
	}
	bytePassword := []byte(req.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, 0)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("there is problem in hashing password: %v\n", err)
	}

	user := entity.User{
		ID:          0,
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		Password:    string(hashedPassword),
	}

	// create new user in storage
	created_user, err := s.repo.Register(user)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	// return created user
	return RegisterResponse{User: created_user}, nil
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type LoginResponse struct {
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	// TODO - it would be better to separate check existence and get user into two functions

	// get user and check existence
	user, exists, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unexpected error: %w", err)
	}
	if !exists {
		return LoginResponse{}, fmt.Errorf("username or password is incorrect")
	}

	// hash request password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 0)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("there is problem in hashing password: %v\n", err)
	}

	// compare password and user's password
	if string(hashedPassword) != user.Password {
		return LoginResponse{}, fmt.Errorf("username or password is incorrect")
	}

	// return ok
	return LoginResponse{}, nil
}
