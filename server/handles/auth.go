package handles

import (
	"errors"
	"fmt"
	"strings"
	"wx_video_help/db"
	"wx_video_help/server/common"
	"wx_video_help/server/middlewares"
	"wx_video_help/utils"
	"wx_video_help/wx"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CheckCookie(c *gin.Context) {
	u := c.MustGet("account").(*db.SphAccount)
	w, err := wx.New(u, nil)
	if err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	if w.Err != nil {
		utils.ResErrWithMsg(c, w.Err.Error())
		return
	}
	utils.ResOk(c)
}

func GetLoginUrl(c *gin.Context) {
	u := wx.User{
		SphAccount: &db.SphAccount{
			Cookie: "",
		},
	}
	loginRes, err := u.GetLoginCode()
	if err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	utils.ResOkWithData(c, gin.H{
		"url":   "https://channels.weixin.qq.com/mobile/confirm_login.html?token=" + loginRes.Data.Token,
		"token": loginRes.Data.Token,
	})
}

func GetLoginStatus(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		utils.ResErrWithMsg(c, "token为空")
		return
	}
	u := wx.User{}
	statusRes, err := u.CheckLoginStatus(token)
	if err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	if statusRes.Data.Cookie == "" {
		utils.ResOkWithData(c, statusRes.Data)
		return
	}
	middlewares.GetUser(c)
	sphUser, exist := c.Get("user")
	fmt.Println("exist", exist)
	if exist && sphUser != nil {
		fmt.Println("add")
		user := sphUser.(*db.SphUser)
		aid, err := addAccount(statusRes.Data.Cookie, user.ID)
		if err != nil {
			utils.ResErrWithMsg(c, err.Error())
			return
		}
		utils.ResOkWithData(c, gin.H{"type": "add", "aid": aid})
	} else {
		fmt.Printf("login")
		tokenString, aid, err := login(statusRes.Data.Cookie)
		if err != nil {
			utils.ResErrWithMsg(c, err.Error())
			return
		}
		utils.ResOkWithData(c, gin.H{"token": tokenString, "aid": aid, "type": "login"})
	}
}

func Login(c *gin.Context) {
	var params struct {
		Cookie string `json:"cookie"`
	}
	if err := c.BindJSON(&params); err != nil {
		utils.ResErrWithMsg(c, "参数错误，"+err.Error())
		return
	}
	if params.Cookie == "" {
		utils.ResErrWithMsg(c, "请填写cookie")
		return
	}
	sphUser, exist := c.Get("user")
	if exist && sphUser != nil {
		user := sphUser.(*db.SphUser)
		aid, err := addAccount(params.Cookie, user.ID)
		if err != nil {
			utils.ResErrWithMsg(c, err.Error())
			return
		}
		utils.ResOkWithData(c, gin.H{"type": "add", "aid": aid})
	} else {
		tokenString, aid, err := login(params.Cookie)
		if err != nil {
			utils.ResErrWithMsg(c, err.Error())
			return
		}
		utils.ResOkWithData(c, gin.H{"token": tokenString, "aid": aid, "type": "login"})
	}
}

func addAccount(cookie string, uid int64) (int64, error) {
	sph := wx.NewWithCookie(cookie)
	auth, err := sph.GetAuthData()
	if err != nil {
		return -1, err
	}
	var account = db.SphAccount{}
	if err := db.Conn.Take(&account, "uniqid like ?", auth.Data.FinderUser.UniqID).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return -1, err
		} else {
			account.UID = uid
			account.NickName = auth.Data.FinderUser.Nickname
			account.Uniqid = auth.Data.FinderUser.UniqID
			account.Cookie = cookie
			if err := db.Conn.Create(&account).Error; err != nil {
				return -1, err
			}
		}
	} else {
		if err := db.Conn.Model(&db.SphAccount{}).Where("id = ?", account.ID).Update("cookie", cookie).Error; err != nil {
			return -1, err
		}
	}
	return account.ID, nil
}
func login(cookie string) (string, int64, error) {
	cookie = strings.TrimSpace(cookie)
	sph := wx.NewWithCookie(cookie)
	auth, err := sph.GetAuthData()
	if err != nil {
		return "", -1, err
	}

	var user = db.SphUser{}
	var account = db.SphAccount{}
	if err := db.Conn.Take(&account, "uniqid like ?", auth.Data.FinderUser.UniqID).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", -1, err
		} else {
			user.NickName = auth.Data.FinderUser.Nickname
			if err := db.Conn.Create(&user).Error; err != nil {
				return "", -1, err
			}
			account.UID = user.ID
			account.NickName = auth.Data.FinderUser.Nickname
			account.Uniqid = auth.Data.FinderUser.UniqID
			account.Cookie = cookie
			if err := db.Conn.Create(&account).Error; err != nil {
				return "", -1, err
			}
		}
	} else {
		if err := db.Conn.Model(&db.SphAccount{}).Where("id = ?", account.ID).Update("cookie", cookie).Error; err != nil {
			return "", -1, err
		}
		if err := db.Conn.Take(&user, "id = ?", account.ID).Error; err != nil {
			return "", -1, err
		}
	}
	tokenString, err := common.GenerateToken(&user)
	if err != nil {
		return "", -1, err
	}
	return tokenString, account.ID, nil
}
