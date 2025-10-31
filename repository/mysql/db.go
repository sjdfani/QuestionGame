package mysql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type MysqlDB struct {
	db *sql.DB
}

func New(cfg MysqlConfig) *MysqlDB {
	dbSourceName := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)

	db, err := sql.Open("mysql", dbSourceName)
	if err != nil {
		panic(fmt.Errorf("can't open mysql db: %v", err))
	}

	_, err = db.Exec("CREATE TABLE users (id INT PRIMARY KEY AUTO_INCREMENT, name VARCHAR(255) NOT NULL,phonenumber VARCHAR(255) NOT NULL UNIQUE,created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP);")
	if err != nil {
		fmt.Printf("Table creation skipped: %v\n", err)
	} else {
		fmt.Println("Table 'users' created successfully.")
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MysqlDB{db: db}
}
