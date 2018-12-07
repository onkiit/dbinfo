package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/onkiit/dbinfo"
)

type mysql struct {
	DB *sql.DB
}

func (m mysql) GetVersion(ctx context.Context) (*dbinfo.DBVersion, error) {
	rows, err := m.DB.Query("SHOW VARIABLES LIKE '%version%'")
	if err != nil {
		fmt.Println("query ", err)
		return nil, nil
	}

	var version, variable, value string
	version = "MySql "
	for rows.Next() {
		err := rows.Scan(&variable, &value)
		if err != nil {
			fmt.Println("scan", err)
			return nil, err
		}
		version += value + " "
	}

	res := &dbinfo.DBVersion{
		Version: version,
	}

	return res, nil
}

func (m mysql) GetActiveClient(ctx context.Context) (*dbinfo.DBActiveClient, error) {
	var count int
	err := m.DB.QueryRowContext(ctx, "SELECT COUNT(0) as count FROM INFORMATION_SCHEMA.PROCESSLIST").Scan(&count)
	if err != nil {
		return nil, err
	}

	res := &dbinfo.DBActiveClient{
		ActiveClient: count,
	}

	return res, nil
}

func (m mysql) GetHealth(ctx context.Context) (*dbinfo.DBHealth, error) {
	err := errors.New("Not available now")
	return nil, err
}

func New(con *dbinfo.Conn) dbinfo.Store {
	return mysql{DB: con.DB}
}
