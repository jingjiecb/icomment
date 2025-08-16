package dao

import (
	"database/sql"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func SetupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func RunSqlFile(t *testing.T, db *sql.DB, filePath string) {
	sqlBytes, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		t.Fatal(err)
	}
}
