package middlewares

import (
	"marketplace/internal/auth"
	"net/http"

	"github.com/labstack/echo/v4"
)

func JwtMiddleware(jwtProvider auth.JWTProvider) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Извлекаем токен из заголовка
			tokenString, err := jwtProvider.ExtractToken(c.Request())
			if err != nil {
				if err.Error() == "authorization header is needed" {
					return next(c)
				}
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": err.Error()})
			}

			// Получаем userID из токена
			userID, err := jwtProvider.ParseIdFromToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"message": "invalid token"})
			}

			// Кладём userID в контекст
			c.Set("userID", userID)

			return next(c)
		}
	}
}
