package jwt

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

type TokenValidator interface {
	ValidateToken(token string) (bool, error)
}

func GinJWTMiddleware(tokenValidator TokenValidator, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Authorization header missing")
			c.AbortWithStatusJSON(401, gin.H{"error": "Authorization header is required"})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			logger.Warn("Invalid authorization header format")
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid authorization header format"})
			return
		}

		tokenStr := fields[1]
		valid, err := tokenValidator.ValidateToken(tokenStr)
		if err != nil || !valid {
			logger.Warn("Invalid or expired token:", err)
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
			return
		}

		parsedToken, _ := jwt.ParseWithClaims(tokenStr, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if claims, ok := parsedToken.Claims.(*jwt.StandardClaims); ok && parsedToken.Valid {
			c.Set("userEmail", claims.Subject)
		} else {
			c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims"})
			return
		}

		c.Next()
	}
}
