package routes

import (
	"microblog/backend/internal/database"
	"microblog/backend/internal/model"
	"microblog/backend/pkg/types"
	"microblog/backend/pkg/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var R *gin.Engine

func Routes() {

	// model.User endpoints
	nomoAPI := R.Group(util.GetPathOnly(util.Getenv("VITE_BACKEND", "/api")))
	nomoAPI.POST("/register", RegisterHandler)
	nomoAPI.POST("/login", LoginHandler)
	nomoAPI.GET("/users", GetAllUsersHandler)
	nomoAPI.Any("/users/me", GetOwnProfileHandler)
	// Thread endpoints
	nomoAPI.GET("/threads", GetAllThreadsHandler)
	nomoAPI.POST("/threads", CreateThreadHandler)
	nomoAPI.GET("/threads/:threadId", GetThreadDetailHandler)
	nomoAPI.POST("/threads/:threadId/comments", CreateThreadCommentHandler)
	nomoAPI.POST("/threads/:threadId/up-vote", UpVoteThreadHandler)
	nomoAPI.POST("/threads/:threadId/down-vote", DownVoteThreadHandler)
	nomoAPI.POST("/threads/:threadId/neutral-vote", NeutralVoteThreadHandler)

	// Comment vote endpoints
	nomoAPI.POST("/threads/:threadId/comments/:commentId/up-vote", UpVoteCommentHandler)
	nomoAPI.POST("/threads/:threadId/comments/:commentId/down-vote", DownVoteCommentHandler)
	nomoAPI.POST("/threads/:threadId/comments/:commentId/neutral-vote", NeutralVoteCommentHandler)

	// CRUD endpoints for threads
	nomoAPI.PUT("/threads/:threadId", UpdateThreadHandler)
	nomoAPI.DELETE("/threads/:threadId", DeleteThreadHandler)

	// CRUD endpoints for comments
	nomoAPI.PUT("/threads/:threadId/comments/:commentId", UpdateCommentHandler)
	nomoAPI.DELETE("/threads/:threadId/comments/:commentId", DeleteCommentHandler)

	// Leaderboard
	nomoAPI.GET("/leaderboards", GetLeaderboardsHandler)
}

// Example handler: Register
func RegisterHandler(c *gin.Context) {
	type RegisterRequest struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	// Check for duplicate email
	var existingUser model.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "fail",
			"message": "email already registered",
		})
		return
	}

	// Generate avatar URL
	avatarURL := "https://ui-avatars.com/api/?name=" + req.Name + "&background=random"

	// Create user object
	user := model.User{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Email:    types.Email(req.Email),
		Avatar:   types.Avatar(avatarURL),
		Password: types.Password(req.Password), // In production, hash the password
	}

	// Save to DB
	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to create user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user created",
		"data": gin.H{
			"user": gin.H{
				"id":     user.ID,
				"name":   user.Name,
				"email":  user.Email,
				"avatar": user.Avatar,
			},
		},
	})
}

// RandString generates a random string of n length
func RandString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(time.Nanosecond) // ensure different seed
	}
	return string(b)
}

// Example handler: Login
func LoginHandler(c *gin.Context) {
	type LoginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	var user model.User
	if err := database.DB.Where("email = ? AND password = ?", req.Email, req.Password).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "email or password is wrong",
			"data":    gin.H{},
		})
		return
	}

	tokenStr := RandString(64)
	token := model.Token{
		ID:       uuid.New().String(),
		UserID:   user.ID,
		CreateAt: time.Now(),
		Token:    tokenStr,
	}
	if err := database.DB.Create(&token).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to save token",
			"data":    gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user logged in",
		"data": gin.H{
			"token": tokenStr,
		},
	})
}

// Example handler: Get all users
func GetAllUsersHandler(c *gin.Context) {
	var users []model.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to get users",
			"data":    gin.H{},
		})
		return
	}

	var result []gin.H
	for _, u := range users {
		result = append(result, gin.H{
			"id":     u.ID,
			"name":   u.Name,
			"email":  u.Email,
			"avatar": u.Avatar,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"users": result,
		},
	})
}

// Get own profile
func GetOwnProfileHandler(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "missing or invalid authorization header",
			"data":    gin.H{},
		})
		return
	}
	tokenStr := authHeader[7:] // Remove "Bearer " prefix

	// Find token in DB
	var token model.Token
	if err := database.DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
		})
		return
	}

	// Find user by token.UserID
	var user model.User
	if err := database.DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "user not found",
			"data":    gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"user": gin.H{
				"id":     user.ID,
				"name":   user.Name,
				"email":  user.Email,
				"avatar": user.Avatar,
			},
		},
	})
}

// Get all threads
func GetAllThreadsHandler(c *gin.Context) {
	var threads []model.Thread
	if err := database.DB.Preload("Comments").Preload("Votes").Find(&threads).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to get threads",
			"data":    gin.H{},
		})
		return
	}

	var result []gin.H
	for _, t := range threads {
		upVotesBy := []string{}
		downVotesBy := []string{}
		for _, v := range t.Votes {
			switch v.VoteType {
			case "up":
				upVotesBy = append(upVotesBy, v.UserID)
			case "down":
				downVotesBy = append(downVotesBy, v.UserID)
			}
		}
		result = append(result, gin.H{
			"id":            t.ID,
			"title":         t.Title,
			"body":          t.Body,
			"category":      t.Category,
			"createdAt":     t.CreatedAt,
			"ownerId":       t.OwnerID,
			"totalComments": len(t.Comments),
			"upVotesBy":     upVotesBy,
			"downVotesBy":   downVotesBy,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "threads retrieved",
		"data": gin.H{
			"threads": result,
		},
	})
}

// Create thread
func CreateThreadHandler(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "Invalid token signature",
			"data":    gin.H{},
		})
		return
	}
	tokenStr := authHeader[7:]

	// Find token in DB
	var token model.Token
	if err := database.DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "Invalid token signature",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user model.User
	if err := database.DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "Invalid token signature",
			"data":    gin.H{},
		})
		return
	}

	// Parse request
	type CreateThreadRequest struct {
		Title    string `json:"title" binding:"required"`
		Body     string `json:"body" binding:"required"`
		Category string `json:"category" binding:"required"`
	}
	var req CreateThreadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	// Create thread object
	thread := model.Thread{
		ID:        uuid.New().String(),
		Title:     req.Title,
		Body:      req.Body,
		Category:  req.Category,
		CreatedAt: time.Now(),
		OwnerID:   user.ID,
	}

	// Save to DB
	if err := database.DB.Create(&thread).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to create thread",
			"data":    gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "thread created",
		"data": gin.H{
			"thread": gin.H{
				"id":        thread.ID,
				"title":     thread.Title,
				"body":      thread.Body,
				"category":  thread.Category,
				"createdAt": thread.CreatedAt.Format(time.RFC3339),
				"owner": gin.H{
					"id":     user.ID,
					"name":   user.Name,
					"email":  user.Email,
					"avatar": user.Avatar,
				},
			},
		},
	})
}
func UpdateThreadHandler(c *gin.Context) {
	threadID := c.Param("threadId")

	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid authorization header",
			"data":    gin.H{},
		})
		return
	}
	tokenStr := authHeader[7:]

	var token model.Token
	if err := database.DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	var thread model.Thread
	if err := database.DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "thread not found",
			"data":    gin.H{},
		})
		return
	}

	if thread.OwnerID != token.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "not authorized to update this thread",
			"data":    gin.H{},
		})
		return
	}

	type UpdateThreadRequest struct {
		Title    string `json:"title"`
		Body     string `json:"body"`
		Category string `json:"category"`
	}
	var req UpdateThreadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	updated := false
	if req.Title != "" {
		thread.Title = req.Title
		updated = true
	}
	if req.Body != "" {
		thread.Body = req.Body
		updated = true
	}
	if req.Category != "" {
		thread.Category = req.Category
		updated = true
	}
	if updated {
		thread.UpdatedAt = time.Now()
		if err := database.DB.Save(&thread).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to update thread",
				"data":    gin.H{},
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "thread updated",
		"data": gin.H{
			"thread": gin.H{
				"id":        thread.ID,
				"title":     thread.Title,
				"body":      thread.Body,
				"category":  thread.Category,
				"createdAt": thread.CreatedAt.Format(time.RFC3339),
				"updatedAt": thread.UpdatedAt.Format(time.RFC3339),
				"ownerId":   thread.OwnerID,
			},
		},
	})
}

func DeleteThreadHandler(c *gin.Context) {
	threadID := c.Param("threadId")

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid authorization header",
			"data":    gin.H{},
		})
		return
	}
	tokenStr := authHeader[7:]

	var token model.Token
	if err := database.DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	var thread model.Thread
	if err := database.DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "thread not found",
			"data":    gin.H{},
		})
		return
	}

	if thread.OwnerID != token.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "not authorized to delete this thread",
			"data":    gin.H{},
		})
		return
	}

	if err := database.DB.Delete(&thread).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to delete thread",
			"data":    gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "thread deleted",
		"data":    gin.H{},
	})
}

func UpdateCommentHandler(c *gin.Context) {
	threadID := c.Param("threadId")
	commentID := c.Param("commentId")

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid authorization header",
			"data":    gin.H{},
		})
		return
	}
	tokenStr := authHeader[7:]

	var token model.Token
	if err := database.DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	var comment model.Comment
	if err := database.DB.Where("id = ? AND thread_id = ?", commentID, threadID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "comment not found",
			"data":    gin.H{},
		})
		return
	}

	if comment.OwnerID != token.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "not authorized to update this comment",
			"data":    gin.H{},
		})
		return
	}

	type UpdateCommentRequest struct {
		Content string `json:"content"`
	}
	var req UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	if req.Content != "" {
		comment.Content = req.Content
		comment.UpdatedAt = time.Now()
		if err := database.DB.Save(&comment).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to update comment",
				"data":    gin.H{},
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "comment updated",
		"data": gin.H{
			"comment": gin.H{
				"id":        comment.ID,
				"content":   comment.Content,
				"createdAt": comment.CreatedAt.Format(time.RFC3339),
				"updatedAt": comment.UpdatedAt.Format(time.RFC3339),
				"ownerId":   comment.OwnerID,
				"threadId":  comment.ThreadID,
			},
		},
	})
}

func DeleteCommentHandler(c *gin.Context) {
	threadID := c.Param("threadId")
	commentID := c.Param("commentId")

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid authorization header",
			"data":    gin.H{},
		})
		return
	}
	tokenStr := authHeader[7:]

	var token model.Token
	if err := database.DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	var comment model.Comment
	if err := database.DB.Where("id = ? AND thread_id = ?", commentID, threadID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "comment not found",
			"data":    gin.H{},
		})
		return
	}

	if comment.OwnerID != token.UserID {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "not authorized to delete this comment",
			"data":    gin.H{},
		})
		return
	}

	if err := database.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to delete comment",
			"data":    gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "comment deleted",
		"data":    gin.H{},
	})
}

// Get thread detail
func GetThreadDetailHandler(c *gin.Context) {
	threadID := c.Param("threadId")
	var thread model.Thread
	// Preload Owner and Comments (with their Owners)
	if err := database.DB.Preload("Owner").
		Preload("Comments.Owner").
		Preload("Comments.Votes").
		Preload("Votes").
		Where("id = ?", threadID).
		First(&thread).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "thread not found",
		})
		return
	}

	upVotesBy := []string{}
	downVotesBy := []string{}
	for _, v := range thread.Votes {
		switch v.VoteType {
		case "up":
			upVotesBy = append(upVotesBy, v.UserID)
		case "down":
			downVotesBy = append(downVotesBy, v.UserID)
		}
	}

	comments := []gin.H{}
	for _, cm := range thread.Comments {
		commentUpVotesBy := []string{}
		commentDownVotesBy := []string{}
		for _, v := range cm.Votes {
			switch v.VoteType {
			case "up":
				commentUpVotesBy = append(commentUpVotesBy, v.UserID)
			case "down":
				commentDownVotesBy = append(commentDownVotesBy, v.UserID)
			}
		}
		comments = append(comments, gin.H{
			"id":        cm.ID,
			"content":   cm.Content,
			"createdAt": cm.CreatedAt.Format(time.RFC3339),
			"owner": gin.H{
				"id":     cm.Owner.ID,
				"name":   cm.Owner.Name,
				"email":  cm.Owner.Email,
				"avatar": cm.Owner.Avatar,
			},
			"upVotesBy":   commentUpVotesBy,
			"downVotesBy": commentDownVotesBy,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"detailThread": gin.H{
				"id":        thread.ID,
				"title":     thread.Title,
				"body":      thread.Body,
				"category":  thread.Category,
				"createdAt": thread.CreatedAt.Format(time.RFC3339),
				"owner": gin.H{
					"id":     thread.Owner.ID,
					"name":   thread.Owner.Name,
					"email":  thread.Owner.Email,
					"avatar": thread.Owner.Avatar,
				},
				"upVotesBy":   upVotesBy,
				"downVotesBy": downVotesBy,
				"comments":    comments,
			},
		},
	})
}

// Create thread comment
func CreateThreadCommentHandler(c *gin.Context) {
	// Get thread ID from URL
	threadID := c.Param("threadId")

	// Find thread
	var thread model.Thread
	if err := database.DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "thread is not found",
			"data":    gin.H{},
		})
		return
	}

	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "Invalid token signature",
			"data":    gin.H{},
		})
		return
	}
	tokenStr := authHeader[7:]

	// Find token in DB
	var token model.Token
	if err := database.DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "Invalid token signature",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user model.User
	if err := database.DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "Invalid token signature",
			"data":    gin.H{},
		})
		return
	}

	// Parse request
	type CreateCommentRequest struct {
		Content string `json:"content" binding:"required"`
	}
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
			"data":    gin.H{},
		})
		return
	}

	// Create comment object
	comment := model.Comment{
		ID:        uuid.New().String(),
		Content:   req.Content,
		CreatedAt: time.Now(),
		OwnerID:   user.ID,
		ThreadID:  thread.ID,
	}

	// Save to DB
	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to create comment",
			"data":    gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "comment created",
		"data": gin.H{
			"comment": gin.H{
				"id":        comment.ID,
				"content":   comment.Content,
				"createdAt": comment.CreatedAt.Format(time.RFC3339),
				"owner": gin.H{
					"id":     user.ID,
					"name":   user.Name,
					"email":  user.Email,
					"avatar": user.Avatar,
				},
			},
		},
	})
}

// Upvote thread
func UpVoteThreadHandler(c *gin.Context) {
	// Get thread ID from URL
	threadID := c.Param("threadId")

	// Validate thread exists
	var thread model.Thread
	if err := database.DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "thread not found",
			"data":    gin.H{},
		})
		return
	}

	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid authorization header",
			"data":    gin.H{},
		})
		return
	}
	tokenStr := authHeader[7:]

	// Find token in DB
	var token model.Token
	if err := database.DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user model.User
	if err := database.DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "user not found",
			"data":    gin.H{},
		})
		return
	}

	// Check if user already voted
	var vote model.ThreadVote
	if err := database.DB.Where("thread_id = ? AND user_id = ?", threadID, user.ID).First(&vote).Error; err == nil {
		// Update voteType to "up"
		vote.VoteType = "up"
		if err := database.DB.Save(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to update vote",
				"data":    gin.H{},
			})
			return
		}
	} else {
		// Create new vote
		vote = model.ThreadVote{
			ID:       uuid.New().String(),
			ThreadID: threadID,
			UserID:   user.ID,
			VoteType: "up",
		}
		if err := database.DB.Create(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to create vote",
				"data":    gin.H{},
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "thread up voted",
		"data": gin.H{
			"vote": gin.H{
				"id":       vote.ID,
				"threadId": vote.ThreadID,
				"userId":   vote.UserID,
				"voteType": 1,
			},
		},
	})
}

// Downvote thread
func DownVoteThreadHandler(c *gin.Context) {
	// Get thread ID from URL
	threadID := c.Param("threadId")

	// Validate thread exists
	var thread model.Thread
	if err := database.DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "thread not found",
			"data":    gin.H{},
		})
		return
	}

	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid authorization header",
			"data":    gin.H{},
		})
		return
	}
	tokenStr := authHeader[7:]

	// Find token in DB
	var token model.Token
	if err := database.DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user model.User
	if err := database.DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "user not found",
			"data":    gin.H{},
		})
		return
	}

	// Check if user already voted
	var vote model.ThreadVote
	if err := database.DB.Where("thread_id = ? AND user_id = ?", threadID, user.ID).First(&vote).Error; err == nil {
		// Update voteType to "down"
		vote.VoteType = "down"
		if err := database.DB.Save(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to update vote",
				"data":    gin.H{},
			})
			return
		}
	} else {
		// Create new vote
		vote = model.ThreadVote{
			ID:       uuid.New().String(),
			ThreadID: threadID,
			UserID:   user.ID,
			VoteType: "down",
		}
		if err := database.DB.Create(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to create vote",
				"data":    gin.H{},
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "thread down voted",
		"data": gin.H{
			"vote": gin.H{
				"id":       vote.ID,
				"threadId": vote.ThreadID,
				"userId":   vote.UserID,
				"voteType": -1,
			},
		},
	})
}

// Neutral vote thread
func NeutralVoteThreadHandler(c *gin.Context) {
	// Get thread ID from URL
	threadID := c.Param("threadId")

	// Validate thread exists
	var thread model.Thread
	if err := database.DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "thread not found",
			"data":    gin.H{},
		})
		return
	}

	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid authorization header",
			"data":    gin.H{},
		})
		return
	}
	tokenStr := authHeader[7:]

	// Find token in DB
	var token model.Token
	if err := database.DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user model.User
	if err := database.DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "user not found",
			"data":    gin.H{},
		})
		return
	}

	// Check if user already voted
	var vote model.ThreadVote
	if err := database.DB.Where("thread_id = ? AND user_id = ?", threadID, user.ID).First(&vote).Error; err == nil {
		// Update voteType to "neutral"
		vote.VoteType = "neutral"
		if err := database.DB.Save(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to update vote",
				"data":    gin.H{},
			})
			return
		}
	} else {
		// Create new vote
		vote = model.ThreadVote{
			ID:       uuid.New().String(),
			ThreadID: threadID,
			UserID:   user.ID,
			VoteType: "neutral",
		}
		if err := database.DB.Create(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to create vote",
				"data":    gin.H{},
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "thread vote neutralized",
		"data": gin.H{
			"vote": gin.H{
				"id":       vote.ID,
				"threadId": vote.ThreadID,
				"userId":   vote.UserID,
				"voteType": 0,
			},
		},
	})
}

// Upvote comment
func UpVoteCommentHandler(c *gin.Context) {
	threadID := c.Param("threadId")
	commentID := c.Param("commentId")

	// Validate comment exists
	var comment model.Comment
	if err := database.DB.Where("id = ? AND thread_id = ?", commentID, threadID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "comment not found",
			"data":    gin.H{},
		})
		return
	}

	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid authorization header",
			"data":    gin.H{},
		})
		return
	}
	tokenStr := authHeader[7:]

	// Find token in DB
	var token model.Token
	if err := database.DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user model.User
	if err := database.DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "user not found",
			"data":    gin.H{},
		})
		return
	}

	// Check if user already voted
	var vote model.CommentVote
	if err := database.DB.Where("comment_id = ? AND user_id = ?", commentID, user.ID).First(&vote).Error; err == nil {
		// Update voteType to "up"
		vote.VoteType = "up"
		if err := database.DB.Save(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to update vote",
				"data":    gin.H{},
			})
			return
		}
	} else {
		// Create new vote
		vote = model.CommentVote{
			ID:        uuid.New().String(),
			CommentID: commentID,
			UserID:    user.ID,
			VoteType:  "up",
		}
		if err := database.DB.Create(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to create vote",
				"data":    gin.H{},
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "thread comment up voted",
		"data": gin.H{
			"vote": gin.H{
				"id":        vote.ID,
				"threadId":  threadID,
				"commentId": commentID,
				"userId":    user.ID,
				"voteType":  1,
			},
		},
	})
}

// Downvote comment
func DownVoteCommentHandler(c *gin.Context) {
	threadID := c.Param("threadId")
	commentID := c.Param("commentId")

	// Validate comment exists
	var comment model.Comment
	if err := database.DB.Where("id = ? AND thread_id = ?", commentID, threadID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "comment not found",
			"data":    gin.H{},
		})
		return
	}

	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid authorization header",
			"data":    gin.H{},
		})
		return
	}
	tokenStr := authHeader[7:]

	// Find token in DB
	var token model.Token
	if err := database.DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user model.User
	if err := database.DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "user not found",
			"data":    gin.H{},
		})
		return
	}

	// Check if user already voted
	var vote model.CommentVote
	if err := database.DB.Where("comment_id = ? AND user_id = ?", commentID, user.ID).First(&vote).Error; err == nil {
		// Update voteType to "down"
		vote.VoteType = "down"
		if err := database.DB.Save(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to update vote",
				"data":    gin.H{},
			})
			return
		}
	} else {
		// Create new vote
		vote = model.CommentVote{
			ID:        uuid.New().String(),
			CommentID: commentID,
			UserID:    user.ID,
			VoteType:  "down",
		}
		if err := database.DB.Create(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to create vote",
				"data":    gin.H{},
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "thread comment down voted",
		"data": gin.H{
			"vote": gin.H{
				"id":        vote.ID,
				"threadId":  threadID,
				"commentId": commentID,
				"userId":    user.ID,
				"voteType":  -1,
			},
		},
	})
}

// Neutral vote comment
func NeutralVoteCommentHandler(c *gin.Context) {
	threadID := c.Param("threadId")
	commentID := c.Param("commentId")

	// Validate comment exists
	var comment model.Comment
	if err := database.DB.Where("id = ? AND thread_id = ?", commentID, threadID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "comment not found",
			"data":    gin.H{},
		})
		return
	}

	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid authorization header",
			"data":    gin.H{},
		})
		return
	}
	tokenStr := authHeader[7:]

	// Find token in DB
	var token model.Token
	if err := database.DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user model.User
	if err := database.DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "user not found",
			"data":    gin.H{},
		})
		return
	}

	// Check if user already voted
	var vote model.CommentVote
	if err := database.DB.Where("comment_id = ? AND user_id = ?", commentID, user.ID).First(&vote).Error; err == nil {
		// Update voteType to "neutral"
		vote.VoteType = "neutral"
		if err := database.DB.Save(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to update vote",
				"data":    gin.H{},
			})
			return
		}
	} else {
		// Create new vote
		vote = model.CommentVote{
			ID:        uuid.New().String(),
			CommentID: commentID,
			UserID:    user.ID,
			VoteType:  "neutral",
		}
		if err := database.DB.Create(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to create vote",
				"data":    gin.H{},
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "thread comment neutral voted",
		"data": gin.H{
			"vote": gin.H{
				"id":        vote.ID,
				"threadId":  threadID,
				"commentId": commentID,
				"userId":    user.ID,
				"voteType":  0,
			},
		},
	})
}

// Get leaderboards
func GetLeaderboardsHandler(c *gin.Context) {
	// Get all users
	var users []model.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "failed to get users",
		})
		return
	}

	leaderboards := []gin.H{}

	for _, user := range users {
		var threadVotesCount int64
		var commentVotesCount int64
		var commentsCount int64

		// Count thread votes (up/down only)
		database.DB.Model(&model.ThreadVote{}).
			Where("user_id = ? AND vote_type IN ?", user.ID, []string{"up", "down"}).
			Count(&threadVotesCount)

		// Count comment votes (up/down only)
		database.DB.Model(&model.CommentVote{}).
			Where("user_id = ? AND vote_type IN ?", user.ID, []string{"up", "down"}).
			Count(&commentVotesCount)

		// Count comments
		database.DB.Model(&model.Comment{}).
			Where("owner_id = ?", user.ID).
			Count(&commentsCount)

		score := int(threadVotesCount+commentVotesCount)*5 + int(commentsCount)*20

		leaderboards = append(leaderboards, gin.H{
			"user": gin.H{
				"id":     user.ID,
				"name":   user.Name,
				"email":  user.Email,
				"avatar": user.Avatar,
			},
			"score": score,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"leaderboards": leaderboards,
		},
	})
}
