package config

import (
	"QuestionGame/repository/mysql"
	"QuestionGame/service/authservice"
)

type HTTPServer struct {
	Port int
}

type Config struct {
	HTTPServer HTTPServer
	Auth       authservice.Config
	Mysql      mysql.MysqlConfig
}
