// model/user.go 
package model

import (
	"time"
	"gorm.io/gorm"
	"primejobs/user-service/internal/service/utils"
	"github.com/google/uuid"
)

type BaseModel struct {
	ID             uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAT      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

type User struct {
	BaseModel      
	Name           string        `gorm:"type:varchar(100);not null" json:"name"`
	Email          string        `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	PasswordHash   string        `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	Password       string        `gorm:"-" json:"password,omitempty"`
}


// BeforeCreate hook 
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashed, err := utils.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.PasswordHash = hashed
		u.Password = ""
	}
	return nil
}

// BeforeUpdate hook
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	if u.Password != "" {
		hashed, err := utils.HashPassword(u.Password)
		if err != nil {
			return err
		}
		tx.Statement.SetColumn("password_hash", hashed)
		u.Password = ""
	}
	return nil
}