package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"web_server/config"
	"web_server/db/models"
	"web_server/internal/store"
	"web_server/pkg/jwt"
	"web_server/pkg/response"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, response.Error(401, "unauthorized"))
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		claims, err := jwt.Parse(config.Default().JWT.Secret, tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.Error(401, "token invalid"))
			c.Abort()
			return
		}
		var u models.User
		if err := store.DB().Where("id = ?", claims.UserID).Preload("Role").First(&u).Error; err != nil {
			c.JSON(http.StatusUnauthorized, response.Error(401, "user not found"))
			c.Abort()
			return
		}
		c.Set("currentUser", &u)
		c.Next()
	}
}
