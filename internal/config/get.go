package config

import (
	"os"
	"strconv"
)

func getEnvStr(key string) string {
	v := os.Getenv(key)
	return v
}

// "true" and "false"
func getEnvBool(key string) bool {
	v := os.Getenv(key)
	return v == "true"
}

func getEnvInt(key string) int {
	s := getEnvStr(key)
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}
