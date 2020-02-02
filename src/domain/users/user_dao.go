package users

// Compare hash password saved to the database to the one provided during login.
// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/prosline/pl_users_api/src/datasources/pg"
	"github.com/prosline/pl_users_api/src/utils/crypto"
	"github.com/prosline/pl_users_api/src/utils/date"
	"github.com/prosline/pl_users_api/src/utils/errors"
	"strings"
)

const (
	// Email unique error violation message
	indexUniqueEmail = "users_email_key"
)

var (
	userDB                  = make(map[int64]*User)
	queryInsertUser         = "INSERT INTO users(first_name,last_name,email,date_created,status,password) VALUES($1,$2,$3,$4,$5,$6) RETURNING id,date_created;"
	querySelectUser         = "SELECT id,first_name,last_name,email,date_created,status FROM users where id = $1;"
	queryUpdateUser         = "UPDATE users SET first_name=$1, last_name=$2 , email=$3, status=$4, password=$5 WHERE Id=$6;"
	queryDeleteUser         = "DELETE FROM users WHERE Id=$1;"
	querySelectUserByStatus = "SELECT id,first_name,last_name,email,date_created,status,password FROM users where status = $1;"
)

func (user *User) Save() *errors.RestErr {
	var tx *sqlx.Tx
	tx = pg.ClientDB.MustBegin()

	stmt, err := tx.Prepare(queryInsertUser)
	var lastInsertId int64
	lastInsertId = 0
	date_created := ""
	// hashedPassword, errCrypt := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	hashedPassword, errCrypt := crypto.HashPassword(user.Password)
	if errCrypt != nil {
		return errors.NewInternalServerError(err.Error())
	}
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer stmt.Close()
	err = tx.QueryRow(queryInsertUser, user.FirstName, user.LastName, user.Email, date.GetTimeNowDB(), user.Status, string(hashedPassword)).Scan(&lastInsertId, &date_created)
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
	user.Password = ""
	return nil
}

func (user *User) Get() *errors.RestErr {
	if dberr := pg.ClientDB.Ping(); dberr != nil {
		panic(dberr)
	}
	stmt, er := pg.ClientDB.Prepare(querySelectUser)
	defer stmt.Close()
	if er != nil {
		return errors.NewInternalServerError(er.Error())
	}
	err := stmt.QueryRow(user.Id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)

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
	//                    first_name=$1, last_name=$2 , email=$3, status=$4, password=$5 WHERE Id=$6
	_, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Status, user.Password, user.Id)
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
	if er != nil {
		return errors.NewInternalServerError(er.Error())
	}
	_, err := stmt.Exec(user.Id)
	if err != nil {
		return errors.NewInternalServerError(er.Error())
	}
	defer stmt.Close()
	return nil
}
func (user *User) FindUsersByStatus(status string) ([]User, *errors.RestErr) {
	if dberr := pg.ClientDB.Ping(); dberr != nil {
		panic(dberr)
	}
	stmt, er := pg.ClientDB.Prepare(querySelectUserByStatus)
	if er != nil {
		return nil, errors.NewInternalServerError(er.Error())
	}
	defer stmt.Close()
	rs, err := stmt.Query(status)
	if err != nil {
		return nil, errors.NewInternalServerError(er.Error())
	}
	defer rs.Close()
	result := make([]User, 0)
	for rs.Next() {
		var u User
		if err := rs.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated, &u.Status, &u.Password); err != nil {
			return nil, errors.NewInternalServerError("error when tying to get users")
		}
		result = append(result, u)
	}
	if len(result) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return result, nil

}
