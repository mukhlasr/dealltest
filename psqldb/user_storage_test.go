package psqldb_test

import (
	"simpleblog/auth"
	"simpleblog/psqldb"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestAddUser(t *testing.T) {
	db := SetupDB(t)
	t.Cleanup(func() {
		CleanupTables(t, db)
	})

	s := (*psqldb.UserStorage)(db)

	expectedUser := auth.StoredUser{
		Username:  "foo",
		Password:  "12345",
		Role:      "user",
		Timestamp: time.Now().UTC().Truncate(time.Millisecond),
	}

	if err := s.AddUser(expectedUser); err != nil {
		t.Fatal("unexpected error:", err)
	}

	var res auth.StoredUser
	if err := db.QueryRow("select username, password, role, timestamp from users where username = $1", expectedUser.Username).
		Scan(&res.Username, &res.Password, &res.Role, &res.Timestamp); err != nil {
		t.Fatal("unexpected error on querying the result:", err)
	}
	res.Timestamp = res.Timestamp.UTC()
	assert.Equal(t, expectedUser, res)
}

func TestGetUserByPhoneNumber(t *testing.T) {
	db := SetupDB(t)
	t.Cleanup(func() {
		CleanupTables(t, db)
	})

	s := (*psqldb.UserStorage)(db)
	expectedUser := auth.StoredUser{
		Username:  "foo",
		Password:  "12345",
		Role:      "user",
		Timestamp: time.Now().UTC().Truncate(time.Millisecond),
	}
	_, err := db.Exec("insert into users(username, password, role, timestamp) values($1, $2, $3, $4)",
		expectedUser.Username,
		expectedUser.Password,
		expectedUser.Role,
		expectedUser.Timestamp,
	)
	if err != nil {
		t.Fatal("failed to insert data for testing:", err)
	}
	t.Run("success", func(t *testing.T) {
		res, err := s.GetUserByUsername(expectedUser.Username)
		if err != nil {
			t.Fatal("unexpected error:", err)
		}
		res.Timestamp = res.Timestamp.UTC()
		assert.Equal(t, expectedUser, res)
	})
	t.Run("unknown phone number", func(t *testing.T) {
		_, err := s.GetUserByUsername("unknownuser")
		assert.ErrorIs(t, err, auth.ErrUsernameNotFound)
	})
}
