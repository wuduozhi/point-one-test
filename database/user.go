package database

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	ID       int64
	Name     string
	CreateAt time.Time
}

type AtUserWeiboRef struct {
	ID       int64
	WeiboID  int64
	AtUserID int64
}

func (this *User) Create() {
	DB.Create(this)
}

func CreateAtUserWeiboRef(tx *gorm.DB, atUserRef *AtUserWeiboRef) error {
	db := DB
	if tx != nil {
		db = tx
	}

	err := db.Model(&AtUserWeiboRef{}).Create(atUserRef).Error

	return err
}
