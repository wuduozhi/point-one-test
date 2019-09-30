package main

import (
	"pointone/config"
	"pointone/database"
	service "pointone/services"

	"github.com/gin-gonic/gin"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		return
	}

	err = database.InitMysql(config.Config.Mysql.DBUrl)
	if err != nil {
		return
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/weibo", service.AddWeibo)
	r.GET("/suggest/:userID", service.GetSuggest)
	r.Run() // listen and serve on 0.0.0.0:8080
}

func testDb() {
	user := &database.User{
		Name: "小智e",
	}
	user.Create()

	weibo := &database.Weibo{
		Text:   "This is the first weibo",
		Ats:    "[1,2]",
		UserID: 1,
	}

	weibo.Create()
}
