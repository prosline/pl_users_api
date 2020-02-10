package users

// Compare hash password saved to the database to the one provided during login.
// err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/prosline/pl_logger/logger"
	"github.com/prosline/pl_users_api/src/datasources/pg"
	"github.com/prosline/pl_util/utils/crypto"
	"github.com/prosline/pl_util/utils/date"
	"github.com/prosline/pl_util/utils/rest_errors"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
)

const (
	// Email unique error violation message
	indexUniqueEmail = "users_email_key"
)

var (
	userDB                   = make(map[int64]*User)
	queryInsertUser          = "INSERT INTO users(first_name,last_name,email,date_created,status,password) VALUES($1,$2,$3,$4,$5,$6) RETURNING id,date_created;"
	querySelectUser          = "SELECT id,first_name,last_name,email,date_created,status FROM users where id = $1;"
	queryUpdateUser          = "UPDATE users SET first_name=$1, last_name=$2 , email=$3, status=$4, password=$5 WHERE id=$6;"
	queryDeleteUser          = "DELETE FROM users WHERE id=$1;"
	querySelectUserByStatus  = "SELECT id,first_name,last_name,email,date_created,status,password FROM users where status = $1;"
	queryFindByEmailPassword = "SELECT id, first_name, last_name, email, date_created,password, status FROM users WHERE email=$1;"
)

func (user *User) Save() *rest_errors.RestErr {
	var tx *sqlx.Tx
	tx = pg.ClientDB.MustBegin()

	stmt, err := tx.Prepare(queryInsertUser)
	var lastInsertId int64
	lastInsertId = 0
	date_created := ""
	// hashedPassword, errCrypt := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	hashedPassword, errCrypt := crypto.HashPassword(user.Password)
	if errCrypt != nil {
		return rest_errors.NewInternalServerError(err.Error(),errCrypt)
	}
	if err != nil {
		logger.Error("Error when try to 'INSERT' user statement", err)
		return rest_errors.NewInternalServerError(err.Error(),err)
	}
	defer stmt.Close()
	err = tx.QueryRow(queryInsertUser, user.FirstName, user.LastName, user.Email, date.GetTimeNowDB(), user.Status, string(hashedPassword)).Scan(&lastInsertId, &date_created)
	if err != nil && strings.Contains(err.Error(), indexUniqueEmail) {
		logger.Error("Error when try to INSERT user to the database", err)
		return rest_errors.NewBadRequestError(fmt.Sprintf("Email %s already exists", user.Email))
	}
	if err != nil {
		logger.Error("Error when try to INSERT the user to the database", err)
		return rest_errors.NewInternalServerError(err.Error(),err)
	}
	if lastInsertId == 0 {
		return rest_errors.NewInternalServerError(fmt.Sprintf("Error while retrieving user id: %s", err.Error()),err)
	}
	tx.Commit()
	user.Id = lastInsertId
	user.DateCreated = date_created
	user.Password = ""
	return nil
}

func (user *User) Get() *rest_errors.RestErr {
	if dberr := pg.ClientDB.Ping(); dberr != nil {
		panic(dberr)
	}
	stmt, er := pg.ClientDB.Prepare(querySelectUser)
	defer stmt.Close()
	if er != nil {
		logger.Error("Error when try to prepare 'SELECT' user statement", er)
		return rest_errors.NewInternalServerError(er.Error(),er)
	}
	err := stmt.QueryRow(user.Id).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)
	if err != nil {
		logger.Error("Error when try to execute the 'SELECT' user by Id", er)
		return rest_errors.NewInternalServerError(fmt.Sprintf("User %s not found", user.Id),err)
	}
	return nil
}

func (user *User) Update() *rest_errors.RestErr {
	if dberr := pg.ClientDB.Ping(); dberr != nil {
		panic(dberr)
	}
	stmt, er := pg.ClientDB.Prepare(queryUpdateUser)
	defer stmt.Close()
	if er != nil {
		logger.Error("Error when try to 'UPDATE' user statement", er)
		return rest_errors.NewInternalServerError(er.Error(),er)
	}
	_, err := stmt.Exec(user.FirstName, user.LastName, user.Email, user.Status, user.Password, user.Id)
	if err != nil {
		logger.Error("Error when try to execute the 'UPDATE' user to the database", err)
		return rest_errors.NewInternalServerError(er.Error(),er)
	}
	return nil
}

func (user *User) Delete() *rest_errors.RestErr {
	if dberr := pg.ClientDB.Ping(); dberr != nil {
		panic(dberr)
	}
	stmt, er := pg.ClientDB.Prepare(queryDeleteUser)
	if er != nil {
		logger.Error("Error when try to 'DELETE' user statement", er)
		return rest_errors.NewInternalServerError(er.Error(),er)
	}
	tmp, err := stmt.Exec(user.Id)
	if n, e := tmp.RowsAffected(); n == 0 && e == nil {
		logger.Error("User ID = "+strconv.Itoa(int(user.Id))+" not found during 'DELETE' process", e)
	}
	if err != nil {
		logger.Error("Error when try to execute the 'DELETE' user in the database", err)
		return rest_errors.NewInternalServerError(er.Error(),err)
	}
	defer stmt.Close()
	return nil
}

func (user *User) FindUsersByStatus(status string) ([]User, *rest_errors.RestErr) {
	if dberr := pg.ClientDB.Ping(); dberr != nil {
		panic(dberr)
	}
	stmt, er := pg.ClientDB.Prepare(querySelectUserByStatus)
	if er != nil {
		logger.Error("Error when try to 'FindUserByStatus' user statement", er)
		return nil, rest_errors.NewInternalServerError(er.Error(),er)
	}
	defer stmt.Close()
	rs, err := stmt.Query(status)
	if err != nil {
		logger.Error("Error when try to execute 'FindUserByStatus' user in the database", err)
		return nil, rest_errors.NewInternalServerError(er.Error(),err)
	}
	defer rs.Close()
	result := make([]User, 0)
	for rs.Next() {
		var u User
		if err := rs.Scan(&u.Id, &u.FirstName, &u.LastName, &u.Email, &u.DateCreated, &u.Status, &u.Password); err != nil {
			logger.Error("Error when try to execute the 'SELECT in the FindUserByStatus method' user statement", er)
			return nil, rest_errors.NewInternalServerError("error when tying to get users",err)
		}
		result = append(result, u)
	}
	if len(result) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("no users matching status %s", status))
	}
	return result, nil

}
func (user *User) FindByEmailAndPassword() *rest_errors.RestErr {
	var txtPassword string
	txtPassword = user.Password
	if strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
		return rest_errors.NewBadRequestError("Incorrect email and/or password")
	}
	if dberr := pg.ClientDB.Ping(); dberr != nil {
		panic(dberr)
	}
	user.GetByEmail()
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(txtPassword)); err != nil {
		return rest_errors.NewBadRequestError("Password not found")
	}
	return nil
}
func (user *User) GetByEmail() *rest_errors.RestErr {
	stmt, er := pg.ClientDB.Prepare(queryFindByEmailPassword)
	if er != nil {
		logger.Error("Error when try to 'FindUserByStatus' user statement", er)
		return rest_errors.NewInternalServerError(er.Error(),er)
	}
	defer stmt.Close()
	if err := stmt.QueryRow(strings.TrimSpace(user.Email)).Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.DateCreated,
		&user.Password,
		&user.Status); err != nil {
		logger.Error("Error when try to execute the 'DELETE' user in the database", err)
		return rest_errors.NewInternalServerError(er.Error(),err)
	}
	return nil
}
