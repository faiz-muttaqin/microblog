package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var R *gin.Engine

type User struct {
	ID       string  `json:"id" gorm:"primaryKey;column:id;size:36"`
	Name     string  `json:"name" gorm:"column:name;size:100"`
	Email    string  `json:"email" gorm:"column:email;size:255;uniqueIndex"`
	Avatar   string  `json:"avatar" gorm:"column:avatar;size:500"`
	Password string  `json:"password" gorm:"column:password;size:255"`
	Tokens   []Token `json:"tokens" gorm:"foreignKey:UserID"`
}

type Token struct {
	ID       string    `json:"id" gorm:"primaryKey;column:id;size:36"`
	UserID   string    `json:"user_id" gorm:"column:user_id;size:36"`
	CreateAt time.Time `json:"created_at" gorm:"autoCreateTime;column:created_at"`
	Token    string    `json:"token" gorm:"column:token;size:255;uniqueIndex"`
}

type Thread struct {
	ID        string       `json:"id" gorm:"primaryKey;column:id;size:36"`
	Title     string       `json:"title" gorm:"column:title;size:255"`
	Body      string       `json:"body" gorm:"column:body;type:text"`
	Category  string       `json:"category" gorm:"column:category;size:100"`
	CreatedAt time.Time    `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time    `json:"updatedAt" gorm:"column:updated_at"`
	OwnerID   string       `json:"owner_id" gorm:"column:owner_id;size:36"`
	Owner     User         `json:"owner" gorm:"foreignKey:OwnerID"`
	Comments  []Comment    `json:"comments" gorm:"foreignKey:ThreadID"`
	Votes     []ThreadVote `json:"votes" gorm:"foreignKey:ThreadID"`
}

type Comment struct {
	ID        string        `json:"id" gorm:"primaryKey;column:id;size:36"`
	Content   string        `json:"content" gorm:"column:content;type:text"`
	CreatedAt time.Time     `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt time.Time     `json:"updatedAt" gorm:"column:updated_at"`
	OwnerID   string        `json:"owner_id" gorm:"column:owner_id;size:36"`
	Owner     User          `json:"owner" gorm:"foreignKey:OwnerID"`
	ThreadID  string        `json:"thread_id" gorm:"column:thread_id;size:36"`
	Votes     []CommentVote `json:"votes" gorm:"foreignKey:CommentID"`
}

type ThreadVote struct {
	ID       string `json:"id" gorm:"primaryKey;column:id;size:36"`
	ThreadID string `json:"thread_id" gorm:"column:thread_id;size:36"`
	UserID   string `json:"user_id" gorm:"column:user_id;size:36"`
	VoteType string `json:"vote_type" gorm:"column:vote_type;size:10"`
}

type CommentVote struct {
	ID        string `json:"id" gorm:"primaryKey;column:id;size:36"`
	CommentID string `json:"comment_id" gorm:"column:comment_id;size:36"`
	UserID    string `json:"user_id" gorm:"column:user_id;size:36"`
	VoteType  string `json:"vote_type" gorm:"column:vote_type;size:10"`
}

func AutoMigrateDB(db *gorm.DB) {
	// Auto migrate all tables
	if err := db.AutoMigrate(&User{}, &Token{}, &Thread{}, &Comment{}, &ThreadVote{}, &CommentVote{}); err != nil {
		logrus.Fatalf("AutoMigrate failed: %v", err)
	}

	// Seed dummy user if table is empty
	var userCount int64
	db.Model(&User{}).Count(&userCount)
	if userCount == 0 {
		dummyUser := User{
			ID:       uuid.New().String(),
			Name:     "Admin",
			Email:    "admin@nomo.com",
			Avatar:   "https://ui-avatars.com/api/?name=Admin&background=random",
			Password: "admin123",
		}
		db.Create(&dummyUser)

		// Seed dummy thread (welcoming message)
		dummyThread := Thread{
			ID:        uuid.New().String(),
			Title:     "Welcome to Nomo Forum!",
			Body:      "This is your first thread. Feel free to post and comment.",
			Category:  "general",
			CreatedAt: time.Now(),
			OwnerID:   dummyUser.ID,
		}
		db.Create(&dummyThread)

		// Seed dummy comment
		dummyComment := Comment{
			ID:        uuid.New().String(),
			Content:   "Say hello to everyone!",
			CreatedAt: time.Now(),
			OwnerID:   dummyUser.ID,
			ThreadID:  dummyThread.ID,
		}
		db.Create(&dummyComment)
	}
}
func main() {
	// Load environment variables using godotenv
	_ = godotenv.Load()
	dbDsn := os.Getenv("DB_PATH")

	var err error
	DB, err = InitAndCheckDB(dbDsn)
	if err != nil {
		logrus.Fatalf("Database setup failed: %v", err)
	}
	go func() {
		AutoMigrateDB(DB)
	}()

	R = gin.Default()
	R.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// User endpoints
	nomoAPI := R.Group("/nomo")
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

	R.Run() // listen and serve on specified port
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
	var existingUser User
	if err := DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"status":  "fail",
			"message": "email already registered",
		})
		return
	}

	// Generate avatar URL
	avatarURL := "https://ui-avatars.com/api/?name=" + req.Name + "&background=random"

	// Create user object
	user := User{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Email:    req.Email,
		Avatar:   avatarURL,
		Password: req.Password, // In production, hash the password
	}

	// Save to DB
	if err := DB.Create(&user).Error; err != nil {
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

	var user User
	if err := DB.Where("email = ? AND password = ?", req.Email, req.Password).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "email or password is wrong",
			"data":    gin.H{},
		})
		return
	}

	tokenStr := RandString(64)
	token := Token{
		ID:       uuid.New().String(),
		UserID:   user.ID,
		CreateAt: time.Now(),
		Token:    tokenStr,
	}
	if err := DB.Create(&token).Error; err != nil {
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
	var users []User
	if err := DB.Find(&users).Error; err != nil {
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
	var token Token
	if err := DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
		})
		return
	}

	// Find user by token.UserID
	var user User
	if err := DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
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
	var threads []Thread
	if err := DB.Preload("Comments").Preload("Votes").Find(&threads).Error; err != nil {
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
	var token Token
	if err := DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "Invalid token signature",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user User
	if err := DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
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
	thread := Thread{
		ID:        uuid.New().String(),
		Title:     req.Title,
		Body:      req.Body,
		Category:  req.Category,
		CreatedAt: time.Now(),
		OwnerID:   user.ID,
	}

	// Save to DB
	if err := DB.Create(&thread).Error; err != nil {
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

	var token Token
	if err := DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	var thread Thread
	if err := DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
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
		if err := DB.Save(&thread).Error; err != nil {
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

	var token Token
	if err := DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	var thread Thread
	if err := DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
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

	if err := DB.Delete(&thread).Error; err != nil {
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

	var token Token
	if err := DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	var comment Comment
	if err := DB.Where("id = ? AND thread_id = ?", commentID, threadID).First(&comment).Error; err != nil {
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
		if err := DB.Save(&comment).Error; err != nil {
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

	var token Token
	if err := DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	var comment Comment
	if err := DB.Where("id = ? AND thread_id = ?", commentID, threadID).First(&comment).Error; err != nil {
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

	if err := DB.Delete(&comment).Error; err != nil {
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
	var thread Thread
	// Preload Owner and Comments (with their Owners)
	if err := DB.Preload("Owner").
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
	var thread Thread
	if err := DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
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
	var token Token
	if err := DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "Invalid token signature",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user User
	if err := DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
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
	comment := Comment{
		ID:        uuid.New().String(),
		Content:   req.Content,
		CreatedAt: time.Now(),
		OwnerID:   user.ID,
		ThreadID:  thread.ID,
	}

	// Save to DB
	if err := DB.Create(&comment).Error; err != nil {
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
	var thread Thread
	if err := DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
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
	var token Token
	if err := DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user User
	if err := DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "user not found",
			"data":    gin.H{},
		})
		return
	}

	// Check if user already voted
	var vote ThreadVote
	if err := DB.Where("thread_id = ? AND user_id = ?", threadID, user.ID).First(&vote).Error; err == nil {
		// Update voteType to "up"
		vote.VoteType = "up"
		if err := DB.Save(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to update vote",
				"data":    gin.H{},
			})
			return
		}
	} else {
		// Create new vote
		vote = ThreadVote{
			ID:       uuid.New().String(),
			ThreadID: threadID,
			UserID:   user.ID,
			VoteType: "up",
		}
		if err := DB.Create(&vote).Error; err != nil {
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
	var thread Thread
	if err := DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
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
	var token Token
	if err := DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user User
	if err := DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "user not found",
			"data":    gin.H{},
		})
		return
	}

	// Check if user already voted
	var vote ThreadVote
	if err := DB.Where("thread_id = ? AND user_id = ?", threadID, user.ID).First(&vote).Error; err == nil {
		// Update voteType to "down"
		vote.VoteType = "down"
		if err := DB.Save(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to update vote",
				"data":    gin.H{},
			})
			return
		}
	} else {
		// Create new vote
		vote = ThreadVote{
			ID:       uuid.New().String(),
			ThreadID: threadID,
			UserID:   user.ID,
			VoteType: "down",
		}
		if err := DB.Create(&vote).Error; err != nil {
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
	var thread Thread
	if err := DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
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
	var token Token
	if err := DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user User
	if err := DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "user not found",
			"data":    gin.H{},
		})
		return
	}

	// Check if user already voted
	var vote ThreadVote
	if err := DB.Where("thread_id = ? AND user_id = ?", threadID, user.ID).First(&vote).Error; err == nil {
		// Update voteType to "neutral"
		vote.VoteType = "neutral"
		if err := DB.Save(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to update vote",
				"data":    gin.H{},
			})
			return
		}
	} else {
		// Create new vote
		vote = ThreadVote{
			ID:       uuid.New().String(),
			ThreadID: threadID,
			UserID:   user.ID,
			VoteType: "neutral",
		}
		if err := DB.Create(&vote).Error; err != nil {
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
	var comment Comment
	if err := DB.Where("id = ? AND thread_id = ?", commentID, threadID).First(&comment).Error; err != nil {
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
	var token Token
	if err := DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user User
	if err := DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "user not found",
			"data":    gin.H{},
		})
		return
	}

	// Check if user already voted
	var vote CommentVote
	if err := DB.Where("comment_id = ? AND user_id = ?", commentID, user.ID).First(&vote).Error; err == nil {
		// Update voteType to "up"
		vote.VoteType = "up"
		if err := DB.Save(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to update vote",
				"data":    gin.H{},
			})
			return
		}
	} else {
		// Create new vote
		vote = CommentVote{
			ID:        uuid.New().String(),
			CommentID: commentID,
			UserID:    user.ID,
			VoteType:  "up",
		}
		if err := DB.Create(&vote).Error; err != nil {
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
	var comment Comment
	if err := DB.Where("id = ? AND thread_id = ?", commentID, threadID).First(&comment).Error; err != nil {
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
	var token Token
	if err := DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user User
	if err := DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "user not found",
			"data":    gin.H{},
		})
		return
	}

	// Check if user already voted
	var vote CommentVote
	if err := DB.Where("comment_id = ? AND user_id = ?", commentID, user.ID).First(&vote).Error; err == nil {
		// Update voteType to "down"
		vote.VoteType = "down"
		if err := DB.Save(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to update vote",
				"data":    gin.H{},
			})
			return
		}
	} else {
		// Create new vote
		vote = CommentVote{
			ID:        uuid.New().String(),
			CommentID: commentID,
			UserID:    user.ID,
			VoteType:  "down",
		}
		if err := DB.Create(&vote).Error; err != nil {
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
	var comment Comment
	if err := DB.Where("id = ? AND thread_id = ?", commentID, threadID).First(&comment).Error; err != nil {
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
	var token Token
	if err := DB.Where("token = ?", tokenStr).First(&token).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "invalid token",
			"data":    gin.H{},
		})
		return
	}

	// Find user by token.UserID
	var user User
	if err := DB.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "user not found",
			"data":    gin.H{},
		})
		return
	}

	// Check if user already voted
	var vote CommentVote
	if err := DB.Where("comment_id = ? AND user_id = ?", commentID, user.ID).First(&vote).Error; err == nil {
		// Update voteType to "neutral"
		vote.VoteType = "neutral"
		if err := DB.Save(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "fail",
				"message": "failed to update vote",
				"data":    gin.H{},
			})
			return
		}
	} else {
		// Create new vote
		vote = CommentVote{
			ID:        uuid.New().String(),
			CommentID: commentID,
			UserID:    user.ID,
			VoteType:  "neutral",
		}
		if err := DB.Create(&vote).Error; err != nil {
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
	var users []User
	if err := DB.Find(&users).Error; err != nil {
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
		DB.Model(&ThreadVote{}).
			Where("user_id = ? AND vote_type IN ?", user.ID, []string{"up", "down"}).
			Count(&threadVotesCount)

		// Count comment votes (up/down only)
		DB.Model(&CommentVote{}).
			Where("user_id = ? AND vote_type IN ?", user.ID, []string{"up", "down"}).
			Count(&commentVotesCount)

		// Count comments
		DB.Model(&Comment{}).
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

// InitPostgreSqlDB initializes and checks the PostgreSQL database connection
func InitPostgreSqlDB(dsn string) (*gorm.DB, error) {
	// Accepts a full DSN string, e.g. "postgresql://user:pass@host:port/dbname?sslmode=verify-full"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Error(err)
		fmt.Println("Failed to connect to PostgreSQL database:", err)
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return db, nil
}
func InitAndCheckDB(dsn string) (*gorm.DB, error) {

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Error(err)
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	// Get the underlying sql.DB object
	sqlDB, err := db.DB()
	if err != nil {
		logrus.Error(err)
		return nil, fmt.Errorf("failed to get db instance: %v", err)
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)           // Set the maximum number of idle connections
	sqlDB.SetMaxOpenConns(100)          // Set the maximum number of open connections
	sqlDB.SetConnMaxLifetime(time.Hour) // Set the maximum lifetime of a connection
	return db, nil
}
