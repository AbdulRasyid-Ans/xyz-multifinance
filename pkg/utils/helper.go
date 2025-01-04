package utils

import (
	"os"
)

func GetEnv(key string) string {
	return os.Getenv(key)
}

func GetEnvWithDefault(key string, defaultValue string) string {
	out := os.Getenv(key)
	if out == "" {
		out = defaultValue
	}

	return out
}

func ParsePagination(page int, pageLimit int) (limit int, offset int) {
	if page <= 0 {
		page = 1
	}

	if pageLimit <= 0 {
		pageLimit = 10
	}

	offset = (page - 1) * pageLimit
	limit = pageLimit

	return limit, offset
}
