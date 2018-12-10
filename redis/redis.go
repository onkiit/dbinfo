package redis

import (
	"context"
	"fmt"
	"strconv"

	"github.com/onkiit/dbinfo"

	redigo "github.com/gomodule/redigo/redis"
)

type redis struct {
	con redigo.Conn
}

func (r redis) GetVersion(ctx context.Context) (*dbinfo.DBVersion, error) {
	info, err := redigo.String(r.con.Do("INFO", "SERVER"))
	if err != nil {
		return nil, err
	}

	strVersion, err := getString(info, "redis_version")
	if err != nil {
		return nil, err
	}

	strOs, err := getString(info, "os")
	if err != nil {
		return nil, err
	}

	strGcc, err := getString(info, "gcc_version")
	if err != nil {
		return nil, err
	}

	v := getValue(strVersion)
	os := getValue(strOs)
	gcc := getValue(strGcc)

	version := fmt.Sprintf("Redis version: %s OS %s gcc_version %s \n", v, os, gcc)

	res := &dbinfo.DBVersion{
		Version: version,
	}

	return res, nil
}

func (r redis) GetActiveClient(ctx context.Context) (*dbinfo.DBActiveClient, error) {
	info, err := redigo.String(r.con.Do("INFO", "CLIENTS"))
	if err != nil {
		return nil, err
	}

	str, err := getString(info, "connected_clients")
	if err != nil {
		return nil, err
	}

	client, err := strconv.ParseInt(getValue(str), 10, 64)
	if err != nil {
		return nil, err
	}

	res := &dbinfo.DBActiveClient{
		ActiveClient: int(client),
	}
	return res, nil
}

func (r redis) GetHealth(ctx context.Context) (*dbinfo.DBHealth, error) {
	defer r.con.Close()
	size, err := redigo.Int(r.con.Do("DBSIZE"))
	if err != nil {
		return nil, err
	}

	info, err := redigo.String(r.con.Do("INFO"))
	if err != nil {
		return nil, err
	}
	usage, err := getUsage(info)
	if err != nil {
		return nil, err
	}

	keys, err := getKeys(info)
	if err != nil {
		return nil, err
	}

	countLog, err := r.getSlowlogCount(r.con)
	if err != nil {
		return nil, err
	}

	stats, err := r.getMemoryStats(r.con)
	if err != nil {
		return nil, err
	}

	res := &dbinfo.DBHealth{
		RedisHealth: dbinfo.RedisHealth{
			AvailableKey: size,
			ExpiredKeys:  keys[0],
			EvictedKeys:  keys[1],
			MemoryUsage:  usage,
			SlowlogCount: countLog,
			MemoryStats: dbinfo.MemoryStats{
				PeakAllocated:    stats[1].(int64),
				TotalAllowed:     stats[3].(int64),
				StartupAllocated: stats[5].(int64),
				PeakPercentage:   0,
				Fragmentation:    0,
			},
		},
	}
	return res, nil
}

func New(con *dbinfo.Conn) dbinfo.Store {
	return redis{con: con.Con}
}
