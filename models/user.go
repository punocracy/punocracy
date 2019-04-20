package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// NewUser creates a new user
func NewUser(db *sqlx.DB) *User {
	user := &User{}
	user.db = db
	user.table = "Users_T"
	user.hasID = true

	return user
}

// PermissionLevel defines the level of privileges a user has
type PermissionLevel int

const (
	// Administrator is the user with the highest privilege
	Administrator PermissionLevel = iota
	// Curator is the user that reviews user submitted phrases
	Curator
	// RegularUser is the user that can rate and submit phrases
	RegularUser
	// NonUser is self explanatory
	NonUser
)

type UserRow struct {
	ID        int64           `db:"id"`
	Username  string          `db:"username"`
	Email     string          `db:"email"`
	Password  string          `db:"password"`
	PermLevel PermissionLevel `db:"permLevel"`
}

type User struct {
	Base
}

func (u *User) userRowFromSqlResult(tx *sqlx.Tx, sqlResult sql.Result) (*UserRow, error) {
	userId, err := sqlResult.LastInsertId()
	if err != nil {
		return nil, err
	}

	return u.GetById(tx, userId)
}

// AllUsers returns all user rows.
func (u *User) AllUsers(tx *sqlx.Tx) ([]*UserRow, error) {
	users := []*UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v", u.table)
	err := u.db.Select(&users, query)

	return users, err
}

// GetById returns record by id.
func (u *User) GetById(tx *sqlx.Tx, id int64) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE id=?", u.table)
	err := u.db.Get(user, query, id)

	return user, err
}

// GetByEmail returns record by email.
func (u *User) GetByEmail(tx *sqlx.Tx, email string) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE email=?", u.table)
	err := u.db.Get(user, query, email)

	return user, err
}

// GetByUsername retrieves a user for the db by username
func (u *User) GetByUsername(tx *sqlx.Tx, username string) (*UserRow, error) {
	user := &UserRow{}
	query := fmt.Sprintf("SELECT * FROM %v WHERE username=?", u.table)
	err := u.db.Get(user, query, username)

	return user, err
}

// GetUserByUsernameAndPassword returns record by email but checks password first.
func (u *User) GetUserByUsernameAndPassword(tx *sqlx.Tx, username, password string) (*UserRow, error) {
	user, err := u.GetByUsername(tx, username)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}

	return user, err
}

// Signup create a new record of user.
func (u *User) Signup(tx *sqlx.Tx, username, email, password, passwordAgain string) (*UserRow, error) {
	if username == "" {
		return nil, errors.New("username cannot be blank")
	}
	if email == "" {
		return nil, errors.New("email cannot be blank")
	}
	if password == "" {
		return nil, errors.New("password cannot be blank")
	}
	if password != passwordAgain {
		return nil, errors.New("password is invalid")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
	if err != nil {
		return nil, err
	}

	data := make(map[string]interface{})
	data["username"] = username
	data["email"] = email
	data["passwordHash"] = hashedPassword
	data["permLevel"] = RegularUser

	sqlResult, err := u.InsertIntoTable(tx, data)
	if err != nil {
		return nil, err
	}

	return u.userRowFromSqlResult(tx, sqlResult)
}

// UpdateEmailAndPasswordById updates user email and password.
func (u *User) UpdateEmailAndPasswordById(tx *sqlx.Tx, userId int64, email, password, passwordAgain string) (*UserRow, error) {
	data := make(map[string]interface{})

	if email != "" {
		data["email"] = email
	}

	if password != "" && passwordAgain != "" && password == passwordAgain {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 5)
		if err != nil {
			return nil, err
		}

		data["password"] = hashedPassword
	}

	if len(data) > 0 {
		_, err := u.UpdateByID(tx, data, userId)
		if err != nil {
			return nil, err
		}
	}

	return u.GetById(tx, userId)
}
