package httpserver

import (
	"QuestionGame/config"
	"QuestionGame/service/authservice"
	"QuestionGame/service/userservice"
	"QuestionGame/validator/uservalidator"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	config        config.Config
	authSvc       authservice.Service
	userSvc       userservice.Service
	userValidator uservalidator.Validator
}

func New(config config.Config, authSvc authservice.Service, userSvc userservice.Service, userValidator uservalidator.Validator) Server {
	return Server{
		config:        config,
		authSvc:       authSvc,
		userSvc:       userSvc,
		userValidator: userValidator,
	}
}

func (s Server) Serve() {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	usersGroup := e.Group("/users")
	usersGroup.POST("/users/register", s.userRegisterHandler)
	usersGroup.POST("/users/login", s.userLoginHandler)
	usersGroup.GET("/users/profile", s.userProfileHandler)

	e.Logger.Fatal(e.Start(fmt.Sprintf("%d", s.config.HTTPServer.Port)))
}
