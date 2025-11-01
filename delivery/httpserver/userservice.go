package httpserver

import (
	"QuestionGame/pkg/httpmsg"
	"QuestionGame/service/userservice"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s Server) userRegisterHandler(c echo.Context) error {
	var req userservice.RegisterRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	response, err := s.userSvc.Register(req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, response)
}

func (s Server) userLoginHandler(c echo.Context) error {
	var lReq userservice.LoginRequest
	if err := c.Bind(&lReq); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	response, err := s.userSvc.Login(lReq)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, response)
}

func (s Server) userProfileHandler(c echo.Context) error {
	authToken := c.Request().Header.Get("Authorization")

	claim, err := s.authSvc.ParseToken(authToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, `{"detail": "Token is not valid"}`)
	}

	response, err := s.userSvc.Profile(userservice.ProfileRequest{UserID: claim.UserID})
	if err != nil {
		message, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, message)
	}

	return c.JSON(http.StatusOK, response)
}
