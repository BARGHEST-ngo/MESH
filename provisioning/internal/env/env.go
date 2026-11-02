package env

import (
	"log/slog"
	"os"
	"strconv"
)

func GetEnv(variable string) string {
	str := os.Getenv(variable)
	if str == "" {
		Fatal("must be set", "variable", variable)
	}
	return str
}

func GetEnvInt(variable string) int {
	str := GetEnv(variable)
	varInt, err := strconv.Atoi(str)
	if err != nil {
		Fatal("failed to parse", "variable", variable)
	}
	return varInt
}

func Fatal(msg string, args ...any) {
	slog.Error(msg, args...)
	os.Exit(1)
}
