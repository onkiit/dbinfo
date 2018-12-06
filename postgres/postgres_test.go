package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/onkiit/dbinfo"

	_ "github.com/lib/pq"
)

func TestGetVersion(t *testing.T) {
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost/Coba")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	con := &dbinfo.Conn{
		DB: db,
	}
	store := New(con)

	ver, err := store.GetVersion(context.Background())
	if err != nil {
		t.Error(err)
	}

	t.Log(ver)
}

func TestGetActiveClient(t *testing.T) {
	db, err := sql.Open("postgres", "postgresql://postgres:postgres@localhost/Coba")
	if err != nil {
		t.Error(err)
	}
	defer db.Close()

	con := &dbinfo.Conn{
		DB: db,
	}
	store := New(con)

	cl, err := store.GetActiveClient(context.Background())
	if err != nil {
		t.Error(err)
	}

	t.Log(cl)
}
