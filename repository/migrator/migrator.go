package migrator

import (
	"QuestionGame/repository/mysql"
	"database/sql"
	"fmt"

	migrate "github.com/rubenv/sql-migrate"
)

type Migrator struct {
	dialect   string
	dbConfig  mysql.MysqlConfig
	migration *migrate.FileMigrationSource
}

func New(dbConfig mysql.MysqlConfig) Migrator {
	migration := &migrate.FileMigrationSource{
		Dir: "./repository/mysql/migrations",
	}

	return Migrator{dialect: "mysql", dbConfig: dbConfig, migration: migration}
}

func (m Migrator) Up() {
	dbSourceName := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
		m.dbConfig.User,
		m.dbConfig.Password,
		m.dbConfig.Host,
		m.dbConfig.Port,
		m.dbConfig.DBName,
	)

	db, err := sql.Open(m.dialect, dbSourceName)
	if err != nil {
		panic(fmt.Errorf("can't open mysql db: %v", err))
	}

	n, err := migrate.Exec(db, m.dialect, m.migration, migrate.Up)
	if err != nil {
		panic(fmt.Errorf("can't apply migrations: %v\n", err))
	}

	fmt.Printf("Applied %d migrations!", n)
}

func (m Migrator) Down() {
	dbSourceName := fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
		m.dbConfig.User,
		m.dbConfig.Password,
		m.dbConfig.Host,
		m.dbConfig.Port,
		m.dbConfig.DBName,
	)

	db, err := sql.Open(m.dialect, dbSourceName)
	if err != nil {
		panic(fmt.Errorf("can't open mysql db: %v", err))
	}

	n, err := migrate.Exec(db, m.dialect, m.migration, migrate.Down)
	if err != nil {
		panic(fmt.Errorf("can't rollback migrations: %v\n", err))
	}

	fmt.Printf("Rollback %d migrations!", n)
}

func (m Migrator) Status() {
	// TODO: add status
}
