package redis

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"

	redigo "github.com/gomodule/redigo/redis"
)

func getString(info string, prefix string) (string, error) {
	var text string
	reader := bufio.NewReader(bytes.NewBufferString(info))
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, prefix) {
			text = line
			break
		}
	}
	return text, nil
}

func getValue(str string) string {
	split := strings.Split(str, ":")
	return split[1]
}

//getting health for redis
func getUsage(info string) (string, error) {
	str, err := getString(info, "used_memory")
	if err != nil {
		return "", err
	}

	usage := getValue(str)
	return usage, nil
}

func getKeys(info string) ([]string, error) {
	strExpired, err := getString(info, "expired_keys")
	if err != nil {
		return nil, err
	}

	strEvicted, err := getString(info, "evicted_keys")
	if err != nil {
		return nil, err
	}

	exp := getValue(strExpired)
	evi := getValue(strEvicted)
	return []string{exp, evi}, nil
}

func (r *redis) getSlowlogCount(con redigo.Conn) (int, error) {
	count, err := redigo.Int(con.Do("SLOWLOG", "LEN"))
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r redis) getMemoryStats(con redigo.Conn) ([]interface{}, error) {
	stats, err := redigo.Values(con.Do("MEMORY", "STATS"))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return stats, nil
}
