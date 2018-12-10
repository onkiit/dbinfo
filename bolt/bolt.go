package bolt

import (
	"context"
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/onkiit/dbinfo"
)

type boltdb struct {
	db *bolt.DB
}

func (b boltdb) GetVersion(ctx context.Context) (*dbinfo.DBVersion, error) {
	return nil, nil
}

func (b boltdb) GetActiveClient(ctx context.Context) (*dbinfo.DBActiveClient, error) {
	dbStats := b.db.Stats()

	res := &dbinfo.DBActiveClient{
		ActiveClient: dbStats.TxN,
	}

	return res, nil
}

func (b boltdb) GetHealth(ctx context.Context) (*dbinfo.DBHealth, error) {
	tx, err := b.db.Begin(false)
	if err != nil {
		return nil, err
	}
	cursor := tx.Cursor()
	bucket := cursor.Bucket()
	stats := bucket.Stats()
	fmt.Printf("health status\n Number of bucket: %d\n Total Keys: %d\n", stats.BucketN, stats.KeyN)
	res := &dbinfo.DBHealth{
		BoltHealth: dbinfo.BoltHealth{
			NumberOfBucket: stats.BucketN,
			NumberOfKey:    stats.KeyN,
		},
	}
	return res, nil
}

func New(con *dbinfo.Conn) dbinfo.Store {
	return boltdb{db: con.BoltDB}
}
