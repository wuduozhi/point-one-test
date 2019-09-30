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

// 存被 at 的人与 weibo，user 的信息
type AtUserWeiboRef struct {
	ID       int64
	WeiboID  int64
	AtUserID int64
	UserID   int64
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

func GetUserByID(userID int64) (*User, error) {
	user := &User{}
	err := DB.Model(&User{}).Where("id = ?", userID).Find(&user).Error

	return user, err
}

func GetOwnerUserIDsAndWeiboIDsMapByAtUserID(atUserID int64) ([]int64, map[int64]struct{}, error) {
	var atUserWeiboRefs []AtUserWeiboRef
	err := DB.Model(&AtUserWeiboRef{}).Where("at_user_id = ?", atUserID).Find(&atUserWeiboRefs).Error

	if err != nil {
		return nil, nil, err
	}

	var ownerUserIDs []int64
	weiboIDAtUserIDMap := make(map[int64]struct{})
	for _, ref := range atUserWeiboRefs {
		ownerUserIDs = append(ownerUserIDs, ref.UserID)
		weiboIDAtUserIDMap[ref.ID] = struct{}{}
	}

	return ownerUserIDs, weiboIDAtUserIDMap, nil
}
