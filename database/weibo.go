package database

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Weibo struct {
	ID       int64
	UserID   int64
	Text     string
	Ats      string
	CreateAt time.Time
}

func (this *Weibo) Create() {
	DB.Create(this)
}

func CreateWeibo(tx *gorm.DB, weibo *Weibo) error {
	db := DB
	if tx != nil {
		db = tx
	}

	err := db.Model(&Weibo{}).Create(weibo).Error

	return err
}
