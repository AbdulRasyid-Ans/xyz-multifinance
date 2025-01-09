package utils

import (
	"math/rand"
	"os"
	"time"
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

func GenerateUniqueString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	uniqueString := string(b)

	return uniqueString
}
