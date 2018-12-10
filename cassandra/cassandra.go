package cassandra

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/gocql/gocql"

	"github.com/onkiit/dbinfo"
)

var health = []string{"ID", "Gossip", "Thrift", "Native", "Load", "Generation", "Uptime", "Heap"}

type cassandradb struct {
	session *gocql.Session
}

func (c cassandradb) GetVersion(ctx context.Context) (*dbinfo.DBVersion, error) {
	q := c.session.Query("select cql_version from system.local;").Iter()
	var cqlVersion string
	q.Scan(&cqlVersion)

	res := &dbinfo.DBVersion{
		Version: fmt.Sprintf("Cassandra version %s CQL version %s", q.Host().Version(), cqlVersion),
	}
	return res, nil
}

func (c cassandradb) GetActiveClient(ctx context.Context) (*dbinfo.DBActiveClient, error) {
	return nil, errors.New("Not Available Now")
}

func (c cassandradb) GetHealth(ctx context.Context) (*dbinfo.DBHealth, error) {
	output, err := exec.Command("nodetool", "info").CombinedOutput()
	if err != nil {
		return nil, err
	}

	var result []string
	for _, v := range health {
		str, err := getString(string(output), v)
		if err != nil {
			return nil, err
		}
		val := getValue(str)
		result = append(result, val)
	}

	res := &dbinfo.DBHealth{
		CassandraHealth: dbinfo.CassandraHealth{
			ID:              result[0],
			GossipActive:    result[1],
			ThriftActive:    result[2],
			NativeTransport: result[3],
			Load:            result[4],
			GenerationNo:    result[5],
			Uptime:          result[6],
		},
	}

	return res, nil
}

func New(con *dbinfo.Conn) dbinfo.Store {
	return cassandradb{session: con.CQLSession}
}
