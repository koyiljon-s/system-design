// model/user.go 
package model

import (
	"time"
	"gorm.io/gorm"
	"github.com/google/uuid"
)

type BaseModel struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type User struct {
	BaseModel      
	Name           string        `gorm:"type:varchar(100);not null" json:"name"`
	Email          string        `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash   *string       `gorm:"column:password_hash;type:varchar(255)" json:"-"`
	Password       string        `gorm:"-" json:"password,omitempty"`

	// OAuth
	Provider       string        `gorm:"type:varchar(30);index" json:"-"`
	ProviderID     string        `gorm:"type:varchar(255);index" json:"-"`
	PictureURL     *string       `gorm:"type:text" json:"picture_url,omitempty"`
}