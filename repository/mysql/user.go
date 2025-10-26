package mysql

import (
	"QuestionGame/entity"
	"database/sql"
	"fmt"
)

func (d *MysqlDB) IsPhoneNumberUnique(phonenumber string) (bool, error) {
	user := entity.User{}
	var createdAt []uint8

	row := d.db.QueryRow("SELECT * FROM users WHERE phonenumber = ?", phonenumber)
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password, &createdAt)
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
	user := entity.User{}
	var createdAt []uint8

	row := d.db.QueryRow("SELECT * FROM users WHERE phonenumber = ?", phonenumber)
	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &user.Password, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, false, nil
		}

		return entity.User{}, false, fmt.Errorf("can not scan row: %w", err)
	}
	return user, true, nil
}
