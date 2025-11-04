package userservice

import (
	"QuestionGame/dto"
	"QuestionGame/entity"
	"QuestionGame/pkg/richerror"

	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Repository interface {
	Register(u entity.User) (entity.User, error)
	GetUserByPhoneNumber(phonenumber string) (entity.User, bool, error)
	GetUserByID(userID uint) (entity.User, error)
}

type AuthGenerator interface {
	CreateAccessToken(user entity.User) (string, error)
	CreateRefreshToken(user entity.User) (string, error)
}

type Service struct {
	auth AuthGenerator
	repo Repository
}

func New(auth AuthGenerator, repo Repository) Service {
	return Service{auth: auth, repo: repo}
}

func (s Service) Register(req dto.RegisterRequest) (dto.RegisterResponse, error) {
	// TODO - we should verify phone number by verification code

	bytePassword := []byte(req.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(bytePassword, 0)
	if err != nil {
		return dto.RegisterResponse{}, fmt.Errorf("there is problem in hashing password: %v\n", err)
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
		return dto.RegisterResponse{}, fmt.Errorf("unexpected error: %w", err)
	}

	// return created user
	return dto.RegisterResponse{
		User: dto.UserInfo{
			ID:          created_user.ID,
			PhoneNumber: created_user.PhoneNumber,
			Name:        created_user.Name,
		},
	}, nil
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
type LoginResponse struct {
	User   dto.UserInfo `json:"user"`
	Tokens Tokens       `json:"tokens"`
}

func (s Service) Login(req LoginRequest) (LoginResponse, error) {
	// TODO - it would be better to separate check existence and get user into two functions
	const op = "userservice.Login"

	// get user and check existence
	user, exists, err := s.repo.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil {
		return LoginResponse{}, richerror.New(op).SetError(err)
	}
	if !exists {
		return LoginResponse{},
			richerror.New(op).SetMessage("username or password is incorrect")
	}

	// hash request password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), 0)
	if err != nil {
		return LoginResponse{},
			richerror.New(op).SetError(err).SetKind(richerror.KindUnexpected).SetMessage("there is problem in hashing password")
	}

	// compare password and user's password
	if string(hashedPassword) != user.Password {
		return LoginResponse{},
			richerror.New(op).SetMessage("username or password is incorrect")
	}

	// create tokens
	accessToken, err := s.auth.CreateAccessToken(user)
	if err != nil {
		return LoginResponse{},
			richerror.New(op).SetError(err).SetMessage("unexpected error")
	}

	refreshToken, err := s.auth.CreateRefreshToken(user)
	if err != nil {
		return LoginResponse{},
			richerror.New(op).SetError(err).SetMessage("unexpected error")
	}

	// return ok
	return LoginResponse{
		User: dto.UserInfo{
			ID:          user.ID,
			PhoneNumber: user.PhoneNumber,
			Name:        user.Name,
		},
		Tokens: Tokens{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}, nil
}

type ProfileRequest struct {
	UserID uint
}

type ProfileResponse struct {
	Name string `json:"name"`
}

func (s Service) Profile(req ProfileRequest) (ProfileResponse, error) {
	const op = "userservice.Profile"

	user, err := s.repo.GetUserByID(req.UserID)
	if err != nil {
		return ProfileResponse{},
			richerror.New(op).SetError(err).SetMeta(map[string]any{"req": req})
	}

	return ProfileResponse{Name: user.Name}, nil
}
