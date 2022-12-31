package psqldb

import (
	"database/sql"
	"fmt"
	"simpleblog/auth"
)

type UserStorage sql.DB

func (s *UserStorage) AddUser(u auth.StoredUser) error {
	db := (*sql.DB)(s)
	_, err := db.Exec("insert into users(username, password, role, timestamp) VALUES($1, $2, $3, $4)", u.Username, u.Password, u.Role, u.Timestamp)
	if err != nil {
		return fmt.Errorf("couldn not execute query: %w", err)
	}
	return nil
}

func (s *UserStorage) GetUserByUsername(phone string) (auth.StoredUser, error) {
	db := (*sql.DB)(s)

	var res auth.StoredUser
	err := db.QueryRow("select username, password, role, timestamp from users where username = $1", phone).
		Scan(&res.Username, &res.Password, &res.Role, &res.Timestamp)

	if err == sql.ErrNoRows {
		return res, auth.ErrUsernameNotFound
	}

	if err != nil {
		return res, fmt.Errorf("failed to query data: %w", err)
	}
	return res, nil
}
