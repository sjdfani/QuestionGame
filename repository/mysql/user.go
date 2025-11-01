package mysql

import (
	"QuestionGame/entity"
	"QuestionGame/pkg/errmsg"
	"QuestionGame/pkg/richerror"
	"database/sql"
	"time"
)

func (d *MysqlDB) IsPhoneNumberUnique(phonenumber string) (bool, error) {
	const op = "mysql.IsPhoneNumberUnique"

	row := d.db.QueryRow("SELECT * FROM users WHERE phonenumber = ?", phonenumber)

	_, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return true, nil
		}

		return false,
			richerror.New(op).SetError(err).SetKind(richerror.KindUnexpected).SetMessage(errmsg.ErrorMsgCantScanQueryScan)
	}
	return false, nil
}

func (d *MysqlDB) Register(u entity.User) (entity.User, error) {
	const op = "mysql.Register"

	res, err := d.db.Exec("INSERT INTO users (name, phonenumber, password) VALUES (?, ?, ?)", u.Name, u.PhoneNumber, u.Password)
	if err != nil {
		return entity.User{},
			richerror.New(op).SetError(err).SetKind(richerror.KindUnexpected).SetMessage(errmsg.ErrorMsgCantExc)
	}
	id, _ := res.LastInsertId()
	u.ID = uint(id)

	return u, nil
}

func (d *MysqlDB) GetUserByPhoneNumber(phonenumber string) (entity.User, bool, error) {
	const op = "mysql.GetUserByPhoneNumber"

	row := d.db.QueryRow("SELECT * FROM users WHERE phonenumber = ?", phonenumber)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, false, nil
		}

		return entity.User{}, false,
			richerror.New(op).SetError(err).SetKind(richerror.KindUnexpected).SetMessage(errmsg.ErrorMsgCantScanQueryScan)
	}

	return user, true, nil
}

func (d *MysqlDB) GetUserByID(userID uint) (entity.User, error) {
	const op = "mysql.GetUserByID"

	row := d.db.QueryRow("SELECT * FROM users WHERE id = ?", userID)

	user, err := scanUser(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return entity.User{},
				richerror.New(op).SetError(err).SetKind(richerror.KindNotFound).SetMessage(errmsg.ErrorMsgNotFound)
		}

		return entity.User{},
			richerror.New(op).SetError(err).SetKind(richerror.KindUnexpected).SetMessage(errmsg.ErrorMsgCantScanQueryScan)
	}

	return user, nil
}

func scanUser(row *sql.Row) (entity.User, error) {
	var createdAt time.Time
	var user entity.User

	err := row.Scan(&user.ID, &user.Name, &user.PhoneNumber, &createdAt, &user.Password)

	return user, err
}
