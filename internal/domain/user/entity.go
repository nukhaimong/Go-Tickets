package user

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name" validate:"required" gorm:"type:varchar(100);not null"`
	Email    string `json:"email" validate:"required,email" gorm:"type:varchar(225); uniqueIndex;not null"`
	Password string `json:"password" validate:"required,min=6" gorm:"type:varchar(100);not null"`
}

func (u *User) hashPassword(password string) error {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashPassword)
	return nil
}

func (u *User) checkPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}
