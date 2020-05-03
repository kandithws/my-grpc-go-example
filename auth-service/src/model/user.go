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

type User struct {
	BaseModel
	XID      string `json:"gid" gorm:"unique; not null"`
	Username string `json:"username" gorm:"unique; not null"`
	Password string
	Email    string `json:"email" gorm:"unique; not null"`
}

func (u *User) BeforeCreate() {
	u.XID = xid.New().String()
}
