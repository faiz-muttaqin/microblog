package handler

import (
	"net/http"

	"microblog/backend/internal/database"
	"microblog/backend/internal/model"

	"github.com/gin-gonic/gin"
)

// UserRole represents a role option
type UserRole struct {
	Label string `json:"label"`
	Value string `json:"value"`
	Icon  string `json:"icon"`
}

// GetRoles returns available user roles
func GetRoles() gin.HandlerFunc {
	return func(c *gin.Context) {
		var roles []model.UserRole
		database.DB.Preload("AbilityRules").Find(&roles)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    roles,
		})

	}
}
