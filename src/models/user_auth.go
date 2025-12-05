package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserAuth struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Username  string `json:"username" gorm:"size:50;not null;unique"`
	Email     string `json:"email" gorm:"size:100;not null;unique"`
	Password  string `json:"-" gorm:"not null"` // json:"-" hides this from API responses
	FullName  string `json:"fullName" gorm:"size:150;not null"`
	Role      string `json:"role" gorm:"type:varchar(20);default:'viewer'"` // admin/manager/viewer
	IsActive  bool   `json:"isActive" gorm:"default:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (UserAuth) TableName() string { return "auth_users" }

// mirrors the instance method UserAuth.prototype.validatePassword
func (u *UserAuth) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// mirrors the beforeCreate/beforeUpdate Sequelize hooks that hash the password
func (u *UserAuth) HashPassword() error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashed)
	return nil
}
