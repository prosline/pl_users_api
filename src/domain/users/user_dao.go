package users

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/prosline/pl_users_api/src/datasources/pg"
	"github.com/prosline/pl_users_api/src/utils/date"
	"github.com/prosline/pl_users_api/src/utils/errors"
	"strings"
)

const (
	// Email unique error violation message
	indexUniqueEmail = "users_email_key"
)

var (
	userDB          = make(map[int64]*User)
	queryInsertUser = "INSERT INTO users(first_name,last_name,email,date_created) VALUES($1,$2,$3,$4) RETURNING id,date_created;"
	querySelectUser = "SELECT id,first_name,last_name,email,date_created FROM users where id = $1;"
	queryUpdateUser = "UPDATE users SET first_name=$1, last_name=$2 , email=$3 WHERE Id=$4;"
	queryDeleteUser = "DELETE FROM users WHERE Id=$1;"
)

func (user *User) Save() *errors.RestErr {
	var tx *sqlx.Tx
	tx = pg.ClientDB.MustBegin()

	stmt, err := tx.Prepare(queryInsertUser)
	var lastInsertId int64
	lastInsertId = 0
	date_created := ""

	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()

	err = tx.QueryRow(queryInsertUser, user.FirstName, user.LastName, user.Email, date.GetTimeNow()).Scan(&lastInsertId, &date_created)
	if err != nil && strings.Contains(err.Error(), indexUniqueEmail) {
		return errors.NewBadRequestError(fmt.Sprintf("Email %s already exists", user.Email))
	}
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	if lastInsertId == 0 {
		return errors.NewInternalServerError(fmt.Sprintf("Error while retrieving user id: %s", err.Error()))
	}
	tx.Commit()
	user.Id = lastInsertId
	user.DateCreated = date_created
	return nil
}

//func (user *User) Get() *errors.RestErr {
//	if dberr := pg.ClientDB.Ping(); dberr != nil {
//		panic(dberr)
//	}
//	stmt, err := pg.ClientDB.Prepare(querySelectUser)
//	if err != nil {
//		fmt.Println("error when trying to prepare get user statement", err)
//		return errors.NewInternalServerError(err.Error())
//	}
//	defer stmt.Close()
//
//	result := stmt.QueryRow(user.Id)
//
//	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated); getErr != nil {
//		fmt.Println("error when trying to get user by id", getErr)
//		return errors.NewInternalServerError(getErr.Error())
//	}
//	return nil
//}
func (user *User) Get() *errors.RestErr {
	if dberr := pg.ClientDB.Ping(); dberr != nil {
		panic(dberr)
	}
	stmt, er := pg.ClientDB.Prepare(querySelectUser)
	defer stmt.Close()
	if er != nil {
		return errors.NewInternalServerError(er.Error())
	}
	err := stmt.QueryRow(user.Id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated)

	//	u := &User{}
	//	err := pg.ClientDB.Get(u, querySelectUser, user.Id)
	if err != nil {
		return errors.NewBadRequestError(fmt.Sprintf("User %s not found", user.Id))
	}
	return nil
}

func (user *User) Update() *errors.RestErr {
	if dberr := pg.ClientDB.Ping(); dberr != nil {
		panic(dberr)
	}
	stmt, er := pg.ClientDB.Prepare(queryUpdateUser)
	defer stmt.Close()
	if er != nil {
		return errors.NewInternalServerError(er.Error())
	}
	_, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)
	if err != nil {
		return errors.NewInternalServerError(er.Error())
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	if dberr := pg.ClientDB.Ping(); dberr != nil {
		panic(dberr)
	}
	stmt, er := pg.ClientDB.Prepare(queryDeleteUser)
	defer stmt.Close()
	if er != nil {
		return errors.NewInternalServerError(er.Error())
	}
	_, err := stmt.Exec(user.Id)
	if err != nil {
		return errors.NewInternalServerError(er.Error())
	}
	return nil
}
