package middleware

import (
	"strings"

	"github.com/AbdulRasyid-Ans/xyz-multifinance/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func LoggerMiddleware() echo.MiddlewareFunc {
	return middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] method=${method}, uri=${uri}, status=${status}, latency=${latency_human}` + "\n",
	})
}

func CORSMiddleware() echo.MiddlewareFunc {
	var allowedOrigins []string

	envOrigins := utils.GetEnvWithDefault("CORS_ALLOWED_ORIGINS", "http://localhost:3000")
	splitOrigins := strings.Split(envOrigins, ",")

	allowedOrigins = append(allowedOrigins, splitOrigins...)

	return middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: allowedOrigins,
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowHeaders: []string{echo.HeaderAuthorization, echo.HeaderContentType},
	})
}
