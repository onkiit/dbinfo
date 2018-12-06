package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/onkiit/dbinfo"
)

func (p *psql) getTables(ctx context.Context) (map[string][]string, error) {
	rows, err := p.db.QueryContext(ctx, "select schemaname as schema, relname as table from pg_statio_all_tables where schemaname not in ('pg_catalog', 'pg_toast', 'information_schema')")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	tables := make(map[string][]string)
	for rows.Next() {
		var schema, table string
		err := rows.Scan(&schema, &table)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		tables[schema] = append(tables[schema], table)
	}

	if len(tables) < 1 {
		return nil, errors.New("Could not find table")
	}

	return tables, nil
}

func (p *psql) getTableSize(ctx context.Context) ([]dbinfo.TableInformation, error) {
	tables, err := p.getTables(ctx)
	if err != nil {
		return nil, err
	}

	var res []dbinfo.TableInformation
	for k, v := range tables {
		var info dbinfo.TableInformation
		var tableSize, indexSize string
		if len(v) < 1 {
			return nil, errors.New("Schema has no table")
		}

		for _, val := range v {
			qTable := fmt.Sprintf("SELECT (pg_total_relation_size('%s.%s')) as tableSize", k, val)
			qIndex := fmt.Sprintf("SELECT (pg_indexes_size('%s.%s')) as indexSize", k, val)
			err := p.db.QueryRow(qTable).Scan(&tableSize)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			err = p.db.QueryRow(qIndex).Scan(&indexSize)
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			info = dbinfo.TableInformation{
				SchemaName: k,
				TableName:  val,
				TableSize:  tableSize,
				IndexSize:  indexSize,
			}

			res = append(res, info)
		}
	}
	return res, nil
}
