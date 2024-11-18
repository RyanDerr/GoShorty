package request

import "fmt"

type UserAuthInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type EntryInput struct {
	Content string `json:"content" binding:"required"`
}

func (u *UserAuthInput) Validate() error {
	if u.Username == "" {
		return fmt.Errorf("username is required")
	}
	if u.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

func (e *EntryInput) Validate() error {
	if e.Content == "" {
		return fmt.Errorf("content is required")
	}
	return nil
}
