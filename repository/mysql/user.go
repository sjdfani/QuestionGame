package mysql

import (
	"QuestionGame/entity"
	"database/sql"
	"fmt"
	"time"
)

func (d *MysqlDB) IsPhoneNumberUnique(phonenumber string) (bool, error) {
	row := d.db.QueryRow("SELECT * FROM users WHERE phonenumber = ?", phonenumber)

	_, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false, fmt.Errorf("can not scan row: %w", err)
	}
	return false, nil
}

func (d *MysqlDB) Register(u entity.User) (entity.User, error) {
	res, err := d.db.Exec("INSERT INTO users (name, phonenumber, password) VALUES (?, ?, ?)", u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return entity.User{}, fmt.Errorf("can not execute command: %w", err)
	}
	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}

func (d *MysqlDB) GetUserByPhoneNumber(phonenumber string) (entity.User, bool, error) {
	row := d.db.QueryRow("SELECT * FROM users WHERE phonenumber = ?", phonenumber)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, false, nil
		}

		return entity.User{}, false, fmt.Errorf("can not scan row: %w", err)
	}

	return user, true, nil
}

func (d *MysqlDB) GetUserByID(userID uint) (entity.User, error) {
	row := d.db.QueryRow("SELECT * FROM users WHERE id = ?", userID)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, fmt.Errorf("record not found")
		}

		return entity.User{}, fmt.Errorf("can not scan row: %w", err)
	}

	return user, nil
}

func scanUser(row *sql.Row) (entity.User, error) {
	var createdAt time.Time
	var user entity.User

	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAt, &user.Password)

	return user, err
}
