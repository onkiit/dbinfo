package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/onkiit/dbinfo"
)

type psql struct {
	db *sql.DB
}

func (p *psql) GetVersion(ctx context.Context) (*dbinfo.DBVersion, error) {
	var v dbinfo.DBVersion
	err := p.db.QueryRowContext(ctx, "SELECT VERSION()").Scan(&v.Version)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (p *psql) GetActiveClient(ctx context.Context) (*dbinfo.DBActiveClient, error) {
	var c dbinfo.DBActiveClient
	err := p.db.QueryRowContext(ctx, "SELECT count(0) as active_client FROM pg_stat_activity where state='active' ").Scan(&c.ActiveClient)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &c, nil
}

func (p *psql) GetHealth(ctx context.Context) (*dbinfo.DBHealth, error) {
	rows, err := p.db.QueryContext(ctx, "select datname, pg_database_size(datname) as size from pg_database where datname = (select current_database()) order by pg_database_size(datname) desc;")
	if err != nil {
		return nil, err
	}

	var datname, size string
	for rows.Next() {
		err := rows.Scan(&datname, &size)
		if err != nil {
			return nil, err
		}
	}

	tableInfo, err := p.getTableSize(ctx)
	if err != nil {
		return nil, err
	}

	res := &dbinfo.DBHealth{
		PsqlHealth: dbinfo.PsqlHealth{
			DBInformation: dbinfo.DBInformation{
				DBName: datname,
				DBSize: size,
			},
			TableInformation: tableInfo,
		},
	}

	return res, nil
}

func New(con *dbinfo.Conn) dbinfo.Store {
	return &psql{db: con.DB}
}
