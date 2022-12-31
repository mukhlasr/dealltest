package psqldb_test

import (
	"database/sql"
	"simpleblog/psqldb"
	"testing"
)

func SetupDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("postgres", "")
	if err != nil {
		t.Fatal("could not connect to the database:", err)
	}
	_, err = db.Exec(psqldb.DBSchema)
	if err != nil {
		t.Fatal("could not setup the database:", err)
	}
	return db
}

func CleanupTables(t *testing.T, db *sql.DB) {
	t.Helper()
	_, err := db.Exec("truncate table users, posts")
	if err != nil {
		t.Fatal("failed to truncate table users", err)
	}
}
