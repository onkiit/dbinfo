package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/onkiit/dbinfo"
)

type sqlitedb struct {
	db *sql.DB
}

func (s *sqlitedb) GetVersion(ctx context.Context) (*dbinfo.DBVersion, error) {
	var version string
	err := s.db.QueryRow("select sqlite_version() as version;").Scan(&version)
	if err != nil {
		return nil, err
	}

	res := &dbinfo.DBVersion{
		Version: fmt.Sprintf("SQLite version %s\n", version),
	}
	return res, nil
}

func (s *sqlitedb) GetActiveClient(ctx context.Context) (*dbinfo.DBActiveClient, error) {
	return nil, nil
}

func (s *sqlitedb) GetHealth(ctx context.Context) (*dbinfo.DBHealth, error) {
	var pageSize, pageCount int
	if err := s.db.QueryRow("select page_size as pageSize, page_count as pageCount from pragma_page_size, pragma_page_count;").Scan(&pageSize, &pageCount); err != nil {
		return nil, err
	}
	fmt.Printf("health_status: \n pragma_page_size: %d\n pragma_page_count: %d\n", pageSize, pageCount)

	res := &dbinfo.DBHealth{
		SQLiteHealth: dbinfo.SQLiteHealth{
			PragmaPageSize:  pageSize,
			PragmaPageCount: pageCount,
		},
	}
	return res, nil
}

func New(con *dbinfo.Conn) dbinfo.Store {
	return &sqlitedb{db: con.DB}
}
