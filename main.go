package main

import (
	"QuestionGame/config"
	"QuestionGame/delivery/httpserver"
	"QuestionGame/repository/mysql"
	"QuestionGame/service/authservice"
	"QuestionGame/service/userservice"
	"QuestionGame/validator/uservalidator"
	"time"
)

const (
	JWTSignKey           = "secret_sign_key"
	AccessTokenSubject   = "at"
	RefreshTokenSubject  = "rt"
	AccessTokenDuration  = time.Hour * 24
	RefreshTokenDuration = time.Hour * 24 * 7
)

func main() {
	cfg := config.Config{
		HTTPServer: config.HTTPServer{Port: 8080},
		Auth: authservice.Config{
			SignKey:               JWTSignKey,
			AccessSubject:         AccessTokenSubject,
			RefreshSubject:        RefreshTokenSubject,
			AccessExpirationTime:  AccessTokenDuration,
			RefreshExpirationTime: RefreshTokenDuration,
		},
		Mysql: mysql.MysqlConfig{
			Host:     "127.0.0.1",
			Port:     3308,
			User:     "root",
			Password: "root_password",
			DBName:   "question_game_db",
		},
	}

	// TODO: add commands for applying
	// mgr := migrator.New(cfg.Mysql)
	// mgr.Up()

	authSvc, userSvc, userValidator := setupServices(cfg)
	server := httpserver.New(cfg, authSvc, userSvc, userValidator)

	server.Serve()
}

func setupServices(cfg config.Config) (authservice.Service, userservice.Service, uservalidator.Validator) {
	authSvc := authservice.New(cfg.Auth)

	mysql := mysql.New(cfg.Mysql)
	userSvc := userservice.New(authSvc, mysql)

	userValidator := uservalidator.New(mysql)

	return authSvc, userSvc, userValidator
}
