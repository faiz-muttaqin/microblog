package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Thread struct {
	ID             string       `json:"id" gorm:"primaryKey;column:id;size:36" ui:"sortable"`
	Title          string       `json:"title" gorm:"column:title;size:255" ui:"creatable;visible;visibility;editable;filterable;;sortable"`
	Body           string       `json:"body" gorm:"column:body;type:text" ui:"creatable;visible;visibility;editable;filterable;;sortable"`
	Category       string       `json:"category" gorm:"column:category;size:100" ui:"creatable;visible;visibility;editable;filterable;;sortable"`
	CreatedAt      time.Time    `json:"created_at" gorm:"column:created_at" ui:"creatable;visible;visibility;editable;filterable;;sortable"`
	UpdatedAt      time.Time    `json:"updated_at" gorm:"column:updated_at" ui:"creatable;visible;visibility;editable;filterable;;sortable"`
	UserID         string       `json:"user_id" gorm:"column:user_id;size:36" ui:"creatable;visible;visibility;editable;filterable;;sortable"`
	User           User         `json:"user" gorm:"foreignKey:UserID" ui:"creatable;visible;visibility;editable;filterable;;sortable"`
	TotalUpVotes   int          `json:"total_up_votes" gorm:"column:total_up_votes" ui:"creatable;visible;visibility;editable;filterable;;sortable"`
	TotalDownVotes int          `json:"total_down_votes" gorm:"column:total_down_votes" ui:"creatable;visible;visibility;editable;filterable;;sortable"`
	UpVotedByMe    bool         `json:"up_voted_by_me" gorm:"-" ui:"visible;sortable"`
	DownVotedByMe  bool         `json:"down_voted_by_me" gorm:"-" ui:"visible;sortable"`
	TotalComments  int          `json:"total_comments" gorm:"column:total_comments" ui:"creatable;visible;visibility;editable;filterable;;sortable"`
	Votes          []ThreadVote `json:"votes" gorm:"foreignKey:ThreadID" ui:"visible;sortable"`
	Comments       []Comment    `json:"comments" gorm:"foreignKey:ThreadID" ui:"visible;sortable"`
}

func (t *Thread) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// TableName overrides the default table name for Thread model
func (Thread) TableName() string {
	return "threads"
}

type Comment struct {
	ID             string        `json:"id" gorm:"primaryKey;column:id;size:36" ui:"visible;sortable"`
	ThreadID       string        `json:"thread_id" gorm:"column:thread_id;size:36" ui:"visible;sortable"`
	UserID         string        `json:"user_id" gorm:"column:user_id;size:36" ui:"visible;sortable"`
	User           User          `json:"user" gorm:"foreignKey:UserID" ui:"visible;sortable"`
	Content        string        `json:"content" gorm:"column:content;type:text" ui:"creatable;visible;sortable"`
	CreatedAt      time.Time     `json:"createdAt" gorm:"column:created_at" ui:"visible;filterable;sortable"`
	UpdatedAt      time.Time     `json:"updatedAt" gorm:"column:updated_at" ui:"visible;filterable;sortable"`
	TotalUpVotes   int           `json:"total_up_votes" gorm:"column:total_up_votes" ui:"visible;filterable;sortable"`
	TotalDownVotes int           `json:"total_down_votes" gorm:"column:total_down_votes" ui:"visible;filterable;sortable"`
	UpVotedByMe    bool          `json:"up_voted_by_me" gorm:"-" ui:"visible"`
	DownVotedByMe  bool          `json:"down_voted_by_me" gorm:"-" ui:"visible"`
	Votes          []CommentVote `json:"votes" gorm:"foreignKey:CommentID" ui:"visible;visibility;sortable"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}
func (c *Comment) AfterCreate(tx *gorm.DB) error {
	if err := tx.Exec(`
			UPDATE threads
			SET 
			total_comments = (
				SELECT COUNT(*)
				FROM comments
				WHERE thread_id = ?
			)
			WHERE id = ?
		`, c.ThreadID, c.ThreadID).Error; err != nil {
		logrus.Println(err)
	}
	return nil
}

// TableName overrides the default table name for Comment model
func (Comment) TableName() string {
	return "comments"
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

// TableName overrides the default table name for ThreadVote model
func (ThreadVote) TableName() string {
	return "thread_votes"
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

// TableName overrides the default table name for CommentVote model
func (CommentVote) TableName() string {
	return "comment_votes"
}
