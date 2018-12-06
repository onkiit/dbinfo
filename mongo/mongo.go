package mongo

import (
	"context"
	"fmt"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/onkiit/dbinfo"
)

type mongodb struct {
	session *mgo.Session
}

func (m *mongodb) GetVersion(ctx context.Context) (*dbinfo.DBVersion, error) {
	info, err := m.session.BuildInfo()
	if err != nil {
		return nil, err
	}

	res := &dbinfo.DBVersion{
		Version: fmt.Sprintf("Mongo version %s", info.Version),
	}

	return res, nil
}

func (m *mongodb) GetActiveClient(ctx context.Context) (*dbinfo.DBActiveClient, error) {
	var b bson.M
	if err := m.session.DB("test").Run("serverStatus", &b); err != nil {
		return nil, err
	}

	client := b["globalLock"].(bson.M)["activeClients"].(bson.M)["total"]

	res := &dbinfo.DBActiveClient{
		ActiveClient: client.(int),
	}
	return res, nil
}

func (m *mongodb) GetHealth(ctx context.Context) (*dbinfo.DBHealth, error) {
	r := bson.M{}
	if err := m.session.DB("test").Run("dbstats", &r); err != nil {
		return nil, err
	}

	res := &dbinfo.DBHealth{
		MongoHealth: dbinfo.MongoHealth{
			DBName:              r["db"].(string),
			AvailableCollection: r["collections"].(int),
			StorageSize:         r["storageSize"].(float64),
			Indexes:             r["indexes"].(int),
			DataSize:            r["dataSize"].(float64),
		},
	}
	return res, nil
}

func New(con *dbinfo.Conn) dbinfo.Store {
	return &mongodb{session: con.Session}
}
