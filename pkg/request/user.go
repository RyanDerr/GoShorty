package request

import (
	"fmt"
	"regexp"
)

const (
	regexAlphanumeric = `^[a-zA-Z0-9]*$`
	// This regex allows alphanumeric characters, as well as the following special characters: .&/@*%-!
	regexPassword     = `^[a-zA-Z0-9.&/@*%!]*$`
)

type UserAuthInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type EntryInput struct {
	Content string `json:"content" binding:"required"`
}

func (u *UserAuthInput) Validate() error {
	switch {
	case u.Username == "":
		return fmt.Errorf("username is required")
	case len(u.Username) < 4:
		return fmt.Errorf("username must be at least 4 characters")
	case len(u.Username) > 32:
		return fmt.Errorf("username must be at most 32 characters")
	case !regexp.MustCompile(regexAlphanumeric).MatchString(u.Username):
		return fmt.Errorf("username must be alphanumeric")
	}

	switch {
	case u.Password == "":
		return fmt.Errorf("password is required")
	case len(u.Password) < 8:
		return fmt.Errorf("password must be at least 8 characters")
	case len(u.Password) > 32:
		return fmt.Errorf("password must be at most 32 characters")
	case !regexp.MustCompile(regexPassword).MatchString(u.Password):
		return fmt.Errorf("password contains invalid characters")
	}
	return nil
}

func (e *EntryInput) Validate() error {
	if e.Content == "" {
		return fmt.Errorf("content is required")
	}
	return nil
}
