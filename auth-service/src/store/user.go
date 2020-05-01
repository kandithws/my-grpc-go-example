package store

import (
	"github.com/jinzhu/gorm"
	"github.com/kandithws/sharespace-api/auth-service/src/common/db"
	"github.com/kandithws/sharespace-api/auth-service/src/model"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore() *UserStore {
	client := db.Client()
	// Migrate Here
	client.AutoMigrate(&model.User{})
	return &UserStore{db: client}
}
