package dbinfo

import (
	"context"
	"database/sql"

	"github.com/gocql/gocql"

	"github.com/globalsign/mgo"
	"github.com/gomodule/redigo/redis"
)

type Store interface {
	GetVersion(context.Context) (*DBVersion, error)
	GetActiveClient(context.Context) (*DBActiveClient, error)
	GetHealth(context.Context) (*DBHealth, error)
}

type Conn struct {
	Session    *mgo.Session
	DB         *sql.DB
	Con        redis.Conn
	CQLSession *gocql.Session
}

type DBVersion struct {
	Version string `json:"version,omitempty"`
}

type DBActiveClient struct {
	ActiveClient int `json:"active_client,omitempty"`
}

type DBHealth struct {
	PsqlHealth      PsqlHealth      `json:"psql_health,omitempty"`
	RedisHealth     RedisHealth     `json:"redis_health,omitempty"`
	MongoHealth     MongoHealth     `json:"mongo_health,omitempty"`
	CassandraHealth CassandraHealth `json:"cassandra_health,omitempty"`
}

type PsqlHealth struct {
	DBInformation    DBInformation      `json:"db_information,omitempty"`
	TableInformation []TableInformation `json:"table_information,omitempty"`
}

type DBInformation struct {
	DBName string `json:"db_name,omitempty"`
	DBSize string `json:"db_size,omitempty"`
}

type TableInformation struct {
	SchemaName string `json:"schema_name,omitempty"`
	TableName  string `json:"table_name,omitempty"`
	TableSize  string `json:"table_size,omitempty"`
	IndexSize  string `json:"index_size,omitempty"`
}

type RedisHealth struct {
	AvailableKey int         `json:"available_key,omitempty"`
	MemoryUsage  string      `json:"memory_usage,omitempty"`
	ExpiredKeys  string      `json:"expired_key,omitempty"`
	EvictedKeys  string      `json:"evicted_key,omitempty"`
	SlowlogCount int         `json:"slow_count,omitempty"`
	MemoryStats  MemoryStats `json:"memory_stats,omitempty"`
}

type MemoryStats struct {
	PeakAllocated    int64 `json:"peak_allocated,omitempty"`
	TotalAllowed     int64 `json:"total_allowed,omitempty"`
	StartupAllocated int64 `json:"startup_allocated,omitempty"`
	PeakPercentage   int64 `json:"peak_percentage,omitempty"`
	Fragmentation    int64 `json:"fragmentation,omitempty"`
}

type MongoHealth struct {
	DBName              string  `json:"dbname,omitempty"`
	AvailableCollection int     `json:"available_collection,omitempty"`
	StorageSize         float64 `json:"storage_size,omitempty"`
	Indexes             int     `json:"indexes,omitempty"`
	DataSize            float64 `json:"data_size,omitempty"`
}

type CassandraHealth struct {
	ID              string `json:"id"`
	GossipActive    string `json:"gossip"`
	ThriftActive    string `json:"thrift"`
	NativeTransport string `json:"native"`
	Load            string `json:"load"`
	GenerationNo    string `json:"gen_number"`
	Uptime          string `json:"uptime"`
}
