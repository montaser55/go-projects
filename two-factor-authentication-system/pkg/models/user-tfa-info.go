package models

import (
	"log"
	"time"

	"github.com/montaser55/two-factor-authentication-service/pkg/config"
	"gorm.io/gorm"
)

var db *gorm.DB

type UserTfaInfo struct {
	Id         int64     `gorm:"primaryKey;not null" sql:"nextval(user_tfa_info_seq)" json:"id"`
	UserId     int64     `gorm:"not null;unique" json:"userId"`
	Sms        bool      `json:"sms"`
	App        bool      `json:"app"`
	SecretKey  string    `json:"secretKey"`
	TryCounter int       `json:"tryCounter"`
	Version    int64     `json:"version"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

func (UserTfaInfo) TableName() string {
	return "user_tfa_info"
}

func init() {
	config.Connect()
	db = config.GetDB()
	db.AutoMigrate(&UserTfaInfo{})
}

func GetAllUserTfaInfos() []UserTfaInfo {
	var UserTfaInfos []UserTfaInfo
	db.Find(&UserTfaInfos)
	return UserTfaInfos
}

func (userTfaInfo *UserTfaInfo) CreateUserTfaInfo() *UserTfaInfo {
	err := db.Create(&userTfaInfo).Error
	if err != nil {
		log.Panic("UserTfaInfo not added")
	}
	return userTfaInfo
}
