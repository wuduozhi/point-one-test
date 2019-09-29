package service

import (
	"encoding/json"
	"net/http"
	"pointone/database"
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
	tx := database.DB.Begin()

	err = database.CreateWeibo(tx, weibo)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusOK, gin.H{"error": err.Error()})
		return
	}

	for _, atUserID := range req.Ats {
		err = database.CreateAtUserWeiboRef(tx, &database.AtUserWeiboRef{
			AtUserID: atUserID,
			WeiboID:  weibo.ID,
		})

		if err != nil {
			tx.Rollback()
			c.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"msg": "ok"})
	return
}
