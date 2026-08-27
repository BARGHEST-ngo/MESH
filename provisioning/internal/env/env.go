package env

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
)

func GetEnv(variable string) string {
	str := os.Getenv(variable)
	if str == "" {
		Fatal("'%s' must be set", variable)
	}
	return str
}

func GetEnvInt(variable string) int {
	str := GetEnv(variable)
	varInt, err := strconv.Atoi(str)
	if err != nil {
		Fatal("failed to parse '%s'", variable)
	}
	return varInt
}

func Fatal(msg string, args ...any) {
	slog.Error(fmt.Sprintf(msg, args...))
	os.Exit(1)
}
