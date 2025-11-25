package providers

import (
	"database/sql"
	"errors"
	"strconv"

	"github.com/Kieran-McDonagh/go-api-template/cmd/internal/types"
	"github.com/mattn/go-sqlite3"
)

type UserProvider struct {
	DB *sql.DB
}

func NewUserProvider(DB *sql.DB) UserProvider {
	return UserProvider{
		DB: DB,
	}
}

func (U UserProvider) Create(newUser types.NewUser) (*string, error) {
	var insertedId int
	err := U.DB.QueryRow("INSERT INTO users(email, username, password, role) VALUES(?,?,?,?) RETURNING id;", newUser.Email, newUser.Username, newUser.Password, newUser.Role).Scan(&insertedId)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) {
			if sqliteErr.Code == sqlite3.ErrConstraint &&
				sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
				return nil, types.ErrNotUniqueEmail
			}
		}

		return nil, err
	}
	idStr := strconv.Itoa(insertedId)
	return &idStr, nil
}

func (U UserProvider) One(ID string) (*types.GetUserResponseBody, error) {
	var user types.GetUserResponseBody
	rows := U.DB.QueryRow("SELECT id, email, username, role FROM users WHERE id = ?", ID)

	if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (U UserProvider) OneByEmail(email string) (*types.User, error) {
	var user types.User
	rows := U.DB.QueryRow("SELECT id, email, username, password, role FROM users WHERE email = ?", email)

	if err := rows.Scan(&user.ID, &user.Email, &user.Username, &user.Password, &user.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
