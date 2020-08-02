package util

import (
	"github.com/anthony-dong/easy-swagger/logger"
	"strings"
)

func GetHost(addr string) (host, port string) {
	split := strings.Split(addr, ":")
	if split == nil || len(split) < 2 {
		logger.FatalF("%s not support this pattern", addr)
	}
	if split[len(split)-2] == "" {
		return "localhost", split[len(split)-1]
	}
	return split[len(split)-2], split[len(split)-1]
}
