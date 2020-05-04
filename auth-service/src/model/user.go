package model

import (
	"time"

	"github.com/rs/xid"
)

type BaseModel struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Use validate only for GRPC forms validation
type User struct {
	BaseModel
	XID      string `json:"id" gorm:"unique; not null"`
	Username string `json:"username" gorm:"unique; not null" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" gorm:"unique; not null" validate:"required"`
}

func (u *User) BeforeCreate() {
	u.XID = xid.New().String()
}
