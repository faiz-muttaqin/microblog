package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// type Token struct {
// 	ID       string    `json:"id" gorm:"primaryKey;column:id;size:36"`
// 	UserID   string    `json:"user_id" gorm:"column:user_id;size:36"`
// 	CreateAt time.Time `json:"created_at" gorm:"autoCreateTime;column:created_at"`
// 	Token    string    `json:"token" gorm:"column:token;size:255;uniqueIndex"`
// }

// func (t *Token) BeforeCreate(tx *gorm.DB) error {
// 	if t.ID == "" {
// 		t.ID = uuid.New().String()
// 	}
// 	return nil
// }

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

func (t *Thread) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
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

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

type ThreadVote struct {
	ID       string `json:"id" gorm:"primaryKey;column:id;size:36"`
	ThreadID string `json:"thread_id" gorm:"column:thread_id;size:36"`
	UserID   string `json:"user_id" gorm:"column:user_id;size:36"`
	VoteType string `json:"vote_type" gorm:"column:vote_type;size:10"`
}

func (tv *ThreadVote) BeforeCreate(tx *gorm.DB) error {
	if tv.ID == "" {
		tv.ID = uuid.New().String()
	}
	return nil
}

type CommentVote struct {
	ID        string `json:"id" gorm:"primaryKey;column:id;size:36"`
	CommentID string `json:"comment_id" gorm:"column:comment_id;size:36"`
	UserID    string `json:"user_id" gorm:"column:user_id;size:36"`
	VoteType  string `json:"vote_type" gorm:"column:vote_type;size:10"`
}

func (cv *CommentVote) BeforeCreate(tx *gorm.DB) error {
	if cv.ID == "" {
		cv.ID = uuid.New().String()
	}
	return nil
}
