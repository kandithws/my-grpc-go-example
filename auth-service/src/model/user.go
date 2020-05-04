package model

import (
	"time"

	"github.com/golang/protobuf/ptypes"
	timestamppb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/rs/xid"
)

type BaseModel struct {
	ID          uint                   `json:"-" gorm:"primary_key"`
	CreatedAt   *time.Time             `json:"-"`
	UpdatedAt   *time.Time             `json:"-"`
	CreatedAtPb *timestamppb.Timestamp `json:"created_at"`
	UpdatedAtPb *timestamppb.Timestamp `json:"updated_at"`
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

func (u *User) AfterFind() {
	cstamp, _ := ptypes.TimestampProto(*u.CreatedAt)
	ustamp, _ := ptypes.TimestampProto(*u.UpdatedAt)
	u.CreatedAtPb = cstamp
	u.UpdatedAtPb = ustamp
}
