package middleware

import (
    "strings"
    "backend/utils"
	"github.com/gin-gonic/gin"
    "net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing or invalid"})
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// ✅ ส่งค่าเข้า context
		c.Set("staff_id", uint(claims["staff_id"].(float64)))
		c.Set("hospital_id", uint(claims["hospital_id"].(float64)))
		c.Set("role", claims["role"].(string))

		c.Next()
	}
}
