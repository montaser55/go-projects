package models

import (
	"log"
	"time"

	"github.com/montaser55/two-factor-authentication-service/pkg/config"
	"gorm.io/gorm"
)

var db *gorm.DB

type UserTfaInfo struct {
	Id         int64     `gorm:"primaryKey;not null" json:"id"`
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

func CreateUserTfaInfo(userTfaInfo *UserTfaInfo) *UserTfaInfo {
	err := db.Create(userTfaInfo).Error
	if err != nil {
		log.Panic("UserTfaInfo not added")
	}
	return userTfaInfo
}

func GetUserTfaInfoByUserId(userId int64) *UserTfaInfo {
	var userTfaInfo UserTfaInfo
	db := db.Where("user_id=?", userId).First(&userTfaInfo)

	if db.Error != nil {
		log.Printf("%v", db.Error)
		return nil
	}
	return &userTfaInfo
}

func UpdateUserTfaInfo(userTfaInfo *UserTfaInfo) *UserTfaInfo {
	db.Save(userTfaInfo)
	return userTfaInfo
}
