package middlewares

import (
	"net/http"
	"strings"
	"wx_video_help/db"
	"wx_video_help/server/common"
	"wx_video_help/utils"

	"github.com/gin-gonic/gin"
)

func AuthUser(c *gin.Context) {
	authUser(c)
	c.Next()
}

func GetUser(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	if tokenString == "" {
		c.Next()
		return
	}
	claims, err := common.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	var u db.SphUser
	if err := db.Conn.Take(&u, "id = ?", claims.UserID).Error; err != nil {
		utils.ResErrWithMsg(c, err.Error())
		c.Abort()
		return
	}
	c.Set("userid", claims.UserID)
	c.Set("user", &u)
	c.Next()
}
func AuthAccount(c *gin.Context) {
	authUser(c)
	authAccount(c)
	c.Next()
}
func AuthMsg(c *gin.Context) {
	authUser(c)
	authAccount(c)
	authMsg(c)
	c.Next()
}

func authUser(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	claims, err := common.ParseToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	var u db.SphUser
	if err := db.Conn.Take(&u, "id = ?", claims.UserID).Error; err != nil {
		utils.ResErrWithMsg(c, err.Error())
		c.Abort()
		return
	}
	c.Set("userid", claims.UserID)
	c.Set("user", &u)
}
func authAccount(c *gin.Context) {
	accountId := c.Query("aid")
	var account db.SphAccount
	if accountId == "" {
		utils.ResErrWithMsg(c, "请选择一个视频号")
		c.Abort()
		return
	} else {
		if err := db.Conn.Take(&account, "id = ?", accountId).Error; err != nil {
			utils.ResErrWithMsg(c, err.Error())
			c.Abort()
			return
		}
	}
	userid := c.MustGet("userid").(int64)
	if userid != account.UID {
		utils.ResErrWithMsg(c, "无权限")
		c.Abort()
		return
	}
	c.Set("accountid", account.ID)
	c.Set("account", &account)
}

func authMsg(c *gin.Context) {
	msgId := c.Query("mid")
	var msg db.SphMsg
	if msgId == "" {
		utils.ResErrWithMsg(c, "请选择一个消息")
		c.Abort()
		return
	} else {
		if err := db.Conn.Take(&msg, "id = ?", msgId).Error; err != nil {
			utils.ResErrWithMsg(c, err.Error())
			c.Abort()
			return
		}
	}
	accountid := c.MustGet("accountid").(int64)
	if accountid != msg.AID {
		utils.ResErrWithMsg(c, "无权限")
		c.Abort()
		return
	}
	c.Set("msgid", msg.ID)
	c.Set("msg", &msg)
}
