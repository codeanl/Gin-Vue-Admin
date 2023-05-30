package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string  `gorm:"type:varchar(20);not null;unique" json:"username"`
	Password     string  `gorm:"size:255;not null" json:"password"`
	Mobile       string  `gorm:"type:varchar(11);not null;unique" json:"mobile"`
	Avatar       string  `gorm:"type:varchar(255)" json:"avatar"`
	Nickname     string  `gorm:"type:varchar(20)" json:"nickname"`
	Introduction string  `gorm:"type:varchar(255)" json:"introduction"`
	Status       uint    `gorm:"type:tinyint(1);default:1;comment:'1正常, 2禁用'" json:"status"`
	Creator      string  `gorm:"type:varchar(20);" json:"creator"`
	Roles        []*Role `gorm:"many2many:user_roles" json:"roles"`
}

const (
	PassWordCost = 12 //密码加密难度
)

// SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
