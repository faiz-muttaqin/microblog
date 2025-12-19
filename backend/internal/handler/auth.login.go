package handler

import (
	"net/http"

	"microblog/backend/internal/helper"
	"microblog/backend/internal/model"

	"github.com/gin-gonic/gin"
)

func GetAuthLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		userData, err := helper.GetFirebaseUser(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "Login failed: user is not active",
				"error":   err.Error(),
			})
			return
		}
		if userData.Status != "active" {

		}
		var user model.User
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Login successful",
			"data":    userData,
			"table": gin.H{
				user.TableName(): user.TableSettings("/users"),
			},
		})
	}
}
