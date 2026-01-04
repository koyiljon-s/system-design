// model/user.go 
package model

import (
	"gorm.io/gorm"
	"primejobs/user-service/internal/service/utils"
)

type User struct {
	gorm.Model     // Includes: ID (uint), CreatedAt (time.Time), UpdatedAt (time.Time), DeletedAt (gorm.DeletedAt)
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