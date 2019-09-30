package service

import (
	"encoding/json"
	"net/http"
	"pointone/database"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type AddWeiboReq struct {
	UserID int64   `form:"user_id" json:"user_id" binding:"required"`
	Text   string  `form:"text" json:"text" binding:"required"`
	Ats    []int64 `form:"ats" json:"ats" binding:"required"`
}

type AddWeiboResp struct {
	Status  string
	ErrCode int
	Msg     string
}

func AddWeibo(c *gin.Context) {
	var req AddWeiboReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Status":  "error",
			"Msg":     err.Error(),
			"ErrCode": 1,
		})
		return
	}

	atsJsonByte, err := json.Marshal(req.Ats)
	if err != nil {
		log.WithError(err).Warning("marsha error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	weibo := &database.Weibo{
		Text:     req.Text,
		UserID:   req.UserID,
		Ats:      string(atsJsonByte),
		CreateAt: time.Now(),
	}

	err = database.CreateWeibo(nil, weibo)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	// 异步处理被 at user 的关系
	go addAtUserRef(weibo.ID, weibo.UserID, req.Ats)

	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
	return
}

func addAtUserRef(weiboID, userID int64, atUserIDs []int64) {
	for _, atUserID := range atUserIDs {
		err := database.CreateAtUserWeiboRef(nil, &database.AtUserWeiboRef{
			AtUserID: atUserID,
			WeiboID:  weiboID,
			UserID:   userID,
		})

		if err != nil {
			log.WithError(err).Infof("create at_user_ref error,weiboID:%v", weiboID)
		}
	}
}

func GetSuggest(c *gin.Context) {
	userIDStr := c.Param("userID")
	userID, err := strconv.Atoi(userIDStr)

	if err != nil {
		log.WithError(err).Infof("get userID error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = database.GetUserByID(int64(userID))

	if err != nil {
		log.WithError(err).Infof("get user error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Infof("process with userID:%v", userID)

	ownerUserIDs, weiboIDMap, err := database.GetOwnerUserIDsAndWeiboIDsMapByAtUserID(int64(userID))
	if err != nil {
		log.WithError(err).Infof("GetOwnerUserIDsAndWeiboIDsMapByAtUserID meet error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	weibos, err := database.GetWeibosByUserIDs(ownerUserIDs)

	if err != nil {
		log.WithError(err).Infof("GetWeibosByUserIDs meet error")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 区分是否在在一条微博中被 at
	var relateUserIDs []int64
	var littleRelateUserIDs []int64

	for _, weibo := range weibos {
		var atUserIDs []int64
		err := json.Unmarshal([]byte(weibo.Ats), &atUserIDs)
		if err != nil {
			log.WithError(err).Infof("Unmarshal weibo ats meet error")
			continue
		}

		// 在同一条 weibo 中被 at
		if _, ok := weiboIDMap[weibo.ID]; ok {
			relateUserIDs = append(relateUserIDs, atUserIDs...)
		} else {
			littleRelateUserIDs = append(littleRelateUserIDs, atUserIDs...)
		}
	}

	resultUserIDs := mergeSuggestIDs(relateUserIDs, littleRelateUserIDs, int64(userID))

	c.JSON(http.StatusOK, gin.H{"msg": "ok", "userIDs": resultUserIDs})
	return
}

// 合并推荐好友
func mergeSuggestIDs(relateUserIDs, littleRelateUserIDs []int64, atUserID int64) []int64 {
	userIDMap := make(map[int64]struct{})
	userIDMap[atUserID] = struct{}{}

	suggestIDs := make([]int64, 0)

	for _, userID := range relateUserIDs {
		if _, ok := userIDMap[userID]; !ok {
			suggestIDs = append(suggestIDs, userID)
			userIDMap[userID] = struct{}{}
		}
	}

	for _, userID := range littleRelateUserIDs {
		if _, ok := userIDMap[userID]; !ok {
			suggestIDs = append(suggestIDs, userID)
			userIDMap[userID] = struct{}{}
		}
	}

	return suggestIDs
}
