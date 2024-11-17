package entity

import (
	"fmt"
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	Content string `gorm:"type:text" json:"content"`
	UserID  uint
}

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Entries  []Entry
}

func (u *User) Save(db *gorm.DB) (*User, error) {
	err := db.Create(&u).Error
	if err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	return u, nil
}

func (u *User) BeforeSave(db *gorm.DB) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	u.Password = string(hash)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	return nil
}
