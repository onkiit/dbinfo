package cassandra

import (
	"bufio"
	"bytes"
	"strings"
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
