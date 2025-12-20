package database

import (
	"fmt"
	"time"

	"microblog/backend/internal/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func AutoMigrateDB(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&model.UserRole{},
		&model.UserAbilityRule{},
		&model.User{},
	); err != nil {
		return err
	}

	// Pastikan role default tersedia
	db.FirstOrCreate(&model.UserRole{
		ID:    1,
		Title: "Super Admin",
		Name:  "superadmin",
		Icon:  "bx bx-sparkle",
	})

	db.FirstOrCreate(&model.UserRole{
		ID:    2,
		Title: "Default", // normal role
		Name:  "default",
		Icon:  "bx bx-radio-circle",
	})

	db.FirstOrCreate(&model.UserRole{
		ID:    3,
		Title: "verified", // super role
		Name:  "verified",
		Icon:  "bx bx-radio-circle",
	})

	// Isi ability rule untuk role default
	var count int64
	db.Model(&model.UserAbilityRule{}).Where("role_id IN ?", []int{1, 2}).Count(&count)
	if count == 0 {
		rules := []model.UserAbilityRule{
			{RoleID: 1, Subject: "*", Read: true},
			{RoleID: 2, Subject: "/", Read: true},
			{RoleID: 2, Subject: "/profile", Read: true, Update: true},
		}
		if err := db.Create(&rules).Error; err != nil {
			return fmt.Errorf("failed creating default abilities: %w", err)
		}
	}
	// Auto migrate all tables
	if err := db.AutoMigrate(
		&model.Thread{},
		&model.Comment{},
		&model.ThreadVote{},
		&model.CommentVote{},
	); err != nil {
		logrus.Fatalf("AutoMigrate failed: %v", err)
	}

	// Seed dummy user if table is empty
	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	if userCount == 0 {
		dummyUser := model.User{
			Name:     "Admin",
			Email:    "admin@nomo.com",
			Avatar:   "https://ui-avatars.com/api/?name=Admin&background=random",
			Password: "admin123",
		}
		if err := db.Create(&dummyUser).Error; err != nil {
			return fmt.Errorf("failed creating dummy user: %w", err)
		}

		// Seed dummy thread (welcoming message)
		dummyThread := model.Thread{
			Title:     "Welcome to Nomo Forum!",
			Body:      "This is your first thread. Feel free to post and comment.",
			Category:  "general",
			CreatedAt: time.Now(),
			UserID:    dummyUser.ID,
		}
		if err := db.Create(&dummyThread).Error; err != nil {
			return fmt.Errorf("failed creating dummy thread: %w", err)
		}

		// Seed dummy comment
		dummyComment := model.Comment{
			Content:   "Say hello to everyone!",
			CreatedAt: time.Now(),
			UserID:    dummyUser.ID,
			ThreadID:  dummyThread.ID,
		}
		if err := db.Create(&dummyComment).Error; err != nil {
			return fmt.Errorf("failed creating dummy comment: %w", err)
		}
	}
	return nil
}
