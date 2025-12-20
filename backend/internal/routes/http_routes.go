package routes

import (
	"microblog/backend/internal/database"
	"microblog/backend/internal/handler"
	"microblog/backend/internal/helper"
	"microblog/backend/internal/model"
	"microblog/backend/pkg/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var R *gin.Engine

func Routes() {

	// model.User endpoints
	backendAPI := R.Group(util.GetPathOnly(util.Getenv("VITE_BACKEND", "/api")))
	backendAPI.GET("/options", handler.GetOptions())
	backendAPI.OPTIONS("/auth/login", handler.GetAuthLogin())
	backendAPI.GET("/auth/login", handler.GetAuthLogin())
	backendAPI.GET("/auth/logout", handler.GetAuthLogout())
	backendAPI.GET("/auth/verify", handler.VerifyAuth()) // Test auth endpoint
	// backendAPI.GET("/roles", handler.GetRoles())
	// backendAPI.GET("/users", handler.GET_DEFAULT_TableDataHandler(database.DB, &model.User{}, []string{"UserRole"}))
	// r.POST("/users", handler.POST_DEFAULT_TableDataHandler(database.DB, &model.User{}, []string{"UserRole"}))
	// r.PATCH("/users", handler.PATCH_DEFAULT_TableDataHandler(database.DB, &model.User{}, []string{"UserRole"}))
	// r.PUT("/users", handler.PUT_DEFAULT_TableDataHandler(database.DB, &model.User{}, []string{"UserRole"}))
	// r.DELETE("/users", handler.DELETE_DEFAULT_TableDataHandler(database.DB, &model.User{}))
	// backendAPI.POST("/register", RegisterHandler)
	// backendAPI.POST("/login", LoginHandler)
	backendAPI.GET("/users", handler.GET_DEFAULT_TABLE(database.DB, &model.User{}, []string{"UserRole"}))
	backendAPI.Any("/users/me", GetOwnProfileHandler)
	// Thread endpoints
	backendAPI.GET("/threads", handler.GET_THREADS_HANDLER(database.DB, []string{"User"}))
	backendAPI.POST("/threads", CreateThreadHandler)
	backendAPI.GET("/threads/:threadId", GetThreadDetailHandler)
	backendAPI.POST("/threads/:threadId/up-vote", UpVoteThreadHandler)
	backendAPI.POST("/threads/:threadId/down-vote", DownVoteThreadHandler)
	backendAPI.POST("/threads/:threadId/neutral-vote", NeutralVoteThreadHandler)

	// Comment vote endpoints
	backendAPI.GET("/threads/:threadId/comments", handler.GET_THREADS_ID_COMMENTS_HANDLER(database.DB, []string{"User"}))
	backendAPI.POST("/threads/:threadId/comments", CreateThreadCommentHandler)
	backendAPI.POST("/threads/:threadId/comments/:commentId/up-vote", UpVoteCommentHandler)
	backendAPI.POST("/threads/:threadId/comments/:commentId/down-vote", DownVoteCommentHandler)
	backendAPI.POST("/threads/:threadId/comments/:commentId/neutral-vote", NeutralVoteCommentHandler)

	// CRUD endpoints for threads
	backendAPI.PUT("/threads/:threadId", UpdateThreadHandler)
	backendAPI.DELETE("/threads/:threadId", DeleteThreadHandler)

	// CRUD endpoints for comments
	backendAPI.PUT("/threads/:threadId/comments/:commentId", UpdateCommentHandler)
	backendAPI.DELETE("/threads/:threadId/comments/:commentId", DeleteCommentHandler)

	// Leaderboard
	backendAPI.GET("/leaderboards", GetLeaderboardsHandler)
}

// Example handler: Get all users
func GetAllUsersHandler(c *gin.Context) {
	var users []model.User
	if err := database.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
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
	user, err := helper.GetFirebaseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    user,
	})
}

// Get all threads
// func GetAllThreadsHandler(c *gin.Context) {
// 	var threads []model.Thread
// 	if err := database.DB.Preload("Comments").Preload("Votes").Find(&threads).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"success": false,
// 			"error":   err.Error(),
// 			"message": "failed to get threads",
// 			"data":    gin.H{},
// 		})
// 		return
// 	}

// 	var result []gin.H
// 	for _, t := range threads {
// 		upVotesBy := []string{}
// 		downVotesBy := []string{}
// 		for _, v := range t.Votes {
// 			switch v.VoteType {
// 			case "up":
// 				upVotesBy = append(upVotesBy, v.UserID)
// 			case "down":
// 				downVotesBy = append(downVotesBy, v.UserID)
// 			}
// 		}
// 		result = append(result, gin.H{
// 			"id":            t.ID,
// 			"title":         t.Title,
// 			"body":          t.Body,
// 			"category":      t.Category,
// 			"createdAt":     t.CreatedAt,
// 			"userId":        t.UserID,
// 			"totalComments": len(t.Comments),
// 			"upVotesBy":     upVotesBy,
// 			"downVotesBy":   downVotesBy,
// 		})
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"success": true,
// 		"message": "threads retrieved",
// 		"data": gin.H{
// 			"threads": result,
// 		},
// 	})
// }

// Create thread
func CreateThreadHandler(c *gin.Context) {
	user, err := helper.GetFirebaseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
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
			"success": false,
			"error":   err.Error(),
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
		UserID:    user.ID,
	}

	// Save to DB
	if err := database.DB.Create(&thread).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "failed to create thread",
			"data":    gin.H{},
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "thread created",
		"data":    thread,
	})
}
func UpdateThreadHandler(c *gin.Context) {
	threadID := c.Param("threadId")

	userData, err := helper.GetFirebaseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var thread model.Thread
	if err := database.DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "thread not found",
			"data":    gin.H{},
		})
		return
	}

	if thread.UserID != userData.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Forbidden",
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
			"success": false,
			"error":   err.Error(),
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
				"success": false,
				"error":   err.Error(),
				"message": "failed to update thread",
				"data":    gin.H{},
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "thread updated",
		"data": gin.H{
			"thread": gin.H{
				"id":        thread.ID,
				"title":     thread.Title,
				"body":      thread.Body,
				"category":  thread.Category,
				"createdAt": thread.CreatedAt.Format(time.RFC3339),
				"updatedAt": thread.UpdatedAt.Format(time.RFC3339),
				"userId":    thread.UserID,
			},
		},
	})
}

func DeleteThreadHandler(c *gin.Context) {
	threadID := c.Param("threadId")

	user, err := helper.GetFirebaseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var thread model.Thread
	if err := database.DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "thread not found",
			"data":    gin.H{},
		})
		return
	}

	if thread.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Forbidden",
			"message": "not authorized to delete this thread",
			"data":    gin.H{},
		})
		return
	}

	if err := database.DB.Delete(&thread).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "failed to delete thread",
			"data":    gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "thread deleted",
		"data":    gin.H{},
	})
}

func UpdateCommentHandler(c *gin.Context) {
	threadID := c.Param("threadId")
	commentID := c.Param("commentId")
	user, err := helper.GetFirebaseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var comment model.Comment
	if err := database.DB.Where("id = ? AND thread_id = ?", commentID, threadID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "comment not found",
			"data":    gin.H{},
		})
		return
	}

	if comment.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Forbidden",
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
			"success": false,
			"error":   err.Error(),
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
				"success": false,
				"error":   err.Error(),
				"message": "failed to update comment",
				"data":    gin.H{},
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "comment updated",
		"data": gin.H{
			"comment": gin.H{
				"id":        comment.ID,
				"content":   comment.Content,
				"createdAt": comment.CreatedAt.Format(time.RFC3339),
				"updatedAt": comment.UpdatedAt.Format(time.RFC3339),
				"userId":    comment.UserID,
				"threadId":  comment.ThreadID,
			},
		},
	})
}

func DeleteCommentHandler(c *gin.Context) {
	threadID := c.Param("threadId")
	commentID := c.Param("commentId")
	user, err := helper.GetFirebaseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	var comment model.Comment
	if err := database.DB.Where("id = ? AND thread_id = ?", commentID, threadID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "comment not found",
			"data":    gin.H{},
		})
		return
	}

	if comment.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"error":   "Forbidden",
			"message": "not authorized to delete this comment",
			"data":    gin.H{},
		})
		return
	}

	if err := database.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "failed to delete comment",
			"data":    gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "comment deleted",
		"data":    gin.H{},
	})
}

// Get thread detail
func GetThreadDetailHandler(c *gin.Context) {
	threadID := c.Param("threadId")
	// Preload User and Comments (with their Users)
	var thread model.Thread
	if err := database.DB.Preload("User").
		Preload("Comments.User").
		Preload("Comments.Votes").
		Preload("Votes").
		Where("id = ?", threadID).
		First(&thread).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "thread not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    thread,
	})
}

// Get thread comment
func GetThreadCommentsHandler(c *gin.Context) {
	// Get thread ID from URL
	threadID := c.Param("threadId")

	// Find thread
	var thread model.Thread
	if err := database.DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "thread is not found",
			"data":    gin.H{},
		})
		return
	}
	user, err := helper.GetFirebaseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
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
			"success": false,
			"error":   err.Error(),
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
		UserID:    user.ID,
		ThreadID:  thread.ID,
	}

	// Save to DB
	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "failed to create comment",
			"data":    gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
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

// Create thread comment
func CreateThreadCommentHandler(c *gin.Context) {
	// Get thread ID from URL
	threadID := c.Param("threadId")

	// Find thread
	var thread model.Thread
	if err := database.DB.Where("id = ?", threadID).First(&thread).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "thread is not found",
			"data":    gin.H{},
		})
		return
	}
	user, err := helper.GetFirebaseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
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
			"success": false,
			"error":   err.Error(),
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
		UserID:    user.ID,
		ThreadID:  thread.ID,
	}

	// Save to DB
	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "failed to create comment",
			"data":    gin.H{},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
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
	user, err := helper.GetFirebaseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	// Get thread ID from URL
	threadID := c.Param("threadId")

	// Validate thread exists
	var thread model.Thread
	if err := database.DB.Select("id").Where("id = ?", threadID).First(&thread).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "thread not found",
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
				"success": false,
				"error":   err.Error(),
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
				"success": false,
				"error":   err.Error(),
				"message": "failed to create vote",
				"data":    gin.H{},
			})
			return
		}

	}
	if err := database.DB.Exec(`
			UPDATE threads
			SET 
			total_up_votes = (
				SELECT COUNT(*)
				FROM thread_votes
				WHERE thread_id = ? AND vote_type = 'up'
			),
			total_down_votes = (
				SELECT COUNT(*)
				FROM thread_votes
				WHERE thread_id = ? AND vote_type = 'down'
			)
			WHERE id = ?
		`, threadID, threadID, threadID).Error; err != nil {
		logrus.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
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
	user, err := helper.GetFirebaseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	// Get thread ID from URL
	threadID := c.Param("threadId")

	// Validate thread exists
	var thread model.Thread
	if err := database.DB.Select("id").Where("id = ?", threadID).First(&thread).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   err.Error(),
			"message": "thread not found",
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
				"success": false,
				"error":   err.Error(),
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
				"success": false,
				"error":   err.Error(),
				"message": "failed to create vote",
				"data":    gin.H{},
			})
			return
		}
	}
	if err := database.DB.Exec(`
			UPDATE threads
			SET 
			total_up_votes = (
				SELECT COUNT(*)
				FROM thread_votes
				WHERE thread_id = ? AND vote_type = 'up'
			),
			total_down_votes = (
				SELECT COUNT(*)
				FROM thread_votes
				WHERE thread_id = ? AND vote_type = 'down'
			)
			WHERE id = ?
		`, threadID, threadID, threadID).Error; err != nil {
		logrus.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
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
			"success": false,
			"error":   err.Error(),
			"message": "thread not found",
			"data":    gin.H{},
		})
		return
	}
	user, err := helper.GetFirebaseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"error":   err.Error(),
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
				"success": false,
				"error":   err.Error(),
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
				"success": false,
				"error":   err.Error(),
				"message": "failed to create vote",
				"data":    gin.H{},
			})
			return
		}
	}
	if err := database.DB.Exec(`
			UPDATE threads
			SET 
			total_up_votes = (
				SELECT COUNT(*)
				FROM thread_votes
				WHERE thread_id = ? AND vote_type = 'up'
			),
			total_down_votes = (
				SELECT COUNT(*)
				FROM thread_votes
				WHERE thread_id = ? AND vote_type = 'down'
			)
			WHERE id = ?
		`, threadID, threadID, threadID).Error; err != nil {
		logrus.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
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
	user, err := helper.GetFirebaseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
		return
	}
	threadID := c.Param("threadId")
	commentID := c.Param("commentId")
	// Validate comment exists
	var comment model.Comment
	if err := database.DB.Select("id").Where("id = ? AND thread_id = ?", commentID, threadID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error(), "message": "comment not found", "data": gin.H{}})
		return
	}

	// Check if user already voted
	var vote model.CommentVote
	if err := database.DB.Where("comment_id = ? AND user_id = ?", commentID, user.ID).First(&vote).Error; err == nil {
		// Update voteType to "up"
		vote.VoteType = "up"
		if err := database.DB.Save(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
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
				"success": false,
				"error":   err.Error(),
				"message": "failed to create vote",
				"data":    gin.H{},
			})
			return
		}
	}
	if err := database.DB.Exec(`
			UPDATE comments SET 
			total_up_votes = (
				SELECT COUNT(*)
				FROM comment_votes
				WHERE comment_id = ? AND vote_type = 'up'
			),
			total_down_votes = (
				SELECT COUNT(*)
				FROM comment_votes
				WHERE comment_id = ? AND vote_type = 'down'
			)
			WHERE id = ?
		`, commentID, commentID, commentID).Error; err != nil {
		logrus.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
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
	user, err := helper.GetFirebaseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
		return
	}
	threadID := c.Param("threadId")
	commentID := c.Param("commentId")

	// Validate comment exists
	var comment model.Comment
	if err := database.DB.Select("id").Where("id = ? AND thread_id = ?", commentID, threadID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error(), "message": "comment not found", "data": gin.H{}})
		return
	}
	// Check if user already voted
	var vote model.CommentVote
	if err := database.DB.Where("comment_id = ? AND user_id = ?", commentID, user.ID).First(&vote).Error; err == nil {
		// Update voteType to "down"
		vote.VoteType = "down"
		if err := database.DB.Save(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
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
				"success": false,
				"error":   err.Error(),
				"message": "failed to create vote",
				"data":    gin.H{},
			})
			return
		}
	}
	if err := database.DB.Exec(`
			UPDATE comments SET 
			total_up_votes = (
				SELECT COUNT(*)
				FROM comment_votes
				WHERE comment_id = ? AND vote_type = 'up'
			),
			total_down_votes = (
				SELECT COUNT(*)
				FROM comment_votes
				WHERE comment_id = ? AND vote_type = 'down'
			)
			WHERE id = ?
		`, commentID, commentID, commentID).Error; err != nil {
		logrus.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
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
	user, err := helper.GetFirebaseUser(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": err.Error()})
		return
	}
	threadID := c.Param("threadId")
	commentID := c.Param("commentId")

	// Validate comment exists
	var comment model.Comment
	if err := database.DB.Select("id").Where("id = ? AND thread_id = ?", commentID, threadID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": err.Error(), "message": "comment not found", "data": gin.H{}})
		return
	}

	// Check if user already voted
	var vote model.CommentVote
	if err := database.DB.Where("comment_id = ? AND user_id = ?", commentID, user.ID).First(&vote).Error; err == nil {
		// Update voteType to "neutral"
		vote.VoteType = "neutral"
		if err := database.DB.Save(&vote).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
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
				"success": false,
				"error":   err.Error(),
				"message": "failed to create vote",
				"data":    gin.H{},
			})
			return
		}
	}
	if err := database.DB.Exec(`
			UPDATE comments SET 
			total_up_votes = (
				SELECT COUNT(*)
				FROM comment_votes
				WHERE comment_id = ? AND vote_type = 'up'
			),
			total_down_votes = (
				SELECT COUNT(*)
				FROM comment_votes
				WHERE comment_id = ? AND vote_type = 'down'
			)
			WHERE id = ?
		`, commentID, commentID, commentID).Error; err != nil {
		logrus.Println(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
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
			"success": false,
			"error":   err.Error(),
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
