package handles

import (
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"wx_video_help/db"
	"wx_video_help/utils"
	"wx_video_help/wx"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AddMsg(c *gin.Context) {
	u := c.MustGet("account").(*db.SphAccount)
	msg := db.SphMsg{
		AID:   u.ID,
		Title: "新消息",
	}
	if err := db.Conn.Create(&msg).Error; err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	utils.ResOk(c)
}

func GetMsgInfo(c *gin.Context) {
	msg := c.MustGet("msg").(*db.SphMsg)
	utils.ResOkWithData(c, gin.H{
		"msg": msg,
	})
}

func SaveMsg(c *gin.Context) {
	var params struct {
		Title string `json:"title"`
		Text  string `json:"text"`
	}
	if err := c.ShouldBind(&params); err != nil {
		utils.ResErrWithMsg(c, "参数错误，"+err.Error())
		return
	}
	msg := c.MustGet("msg").(*db.SphMsg)
	msg.Text = params.Text
	msg.Title = params.Title
	if err := db.Conn.Save(&msg).Error; err != nil {
		utils.ResErrWithMsg(c, "保存失败，"+err.Error())
		return
	}

	u := c.MustGet("account").(*db.SphAccount)
	if err := updateMsg(u); err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	utils.ResOk(c)
}

func GetImg(c *gin.Context) {
	u := c.MustGet("account").(*db.SphAccount)
	msg := c.MustGet("msg").(*db.SphMsg)
	imgMsg := ""
	if msg.URL != "" {
		w, err := wx.New(u, nil)
		if err != nil {
			utils.ResErrWithMsg(c, err.Error())
			return
		}
		img, err := w.GetMediaInfo((*wx.Image)(&msg.Image))
		if err != nil {
			utils.ResErrWithMsg(c, err.Error())
			return
		}
		imgMsg = "data:image/png;base64," + img.Data.ImgContent
	}
	utils.ResOkWithData(c, gin.H{
		"img": imgMsg,
	})
}

func UploadImg(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	// 读取文件内容
	var buf []byte
	buf, err = io.ReadAll(file)
	if err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	fileSize := fileHeader.Size
	// 计算 MD5 值
	hash := md5.Sum(buf)
	md5Str := fmt.Sprintf("%x", hash)

	// 计算 Base64 编码
	base64Str := base64.StdEncoding.EncodeToString(buf)

	img := wx.ImageFile{
		Size:   fileSize,
		Md5:    md5Str,
		Base64: "data:application/octet-stream;base64," + base64Str,
	}
	u := c.MustGet("account").(*db.SphAccount)
	w, err := wx.New(u, nil)
	if err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	imgRes, err := w.UploadImage(&img, "10086")
	if err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	msg := c.MustGet("msg").(*db.SphMsg)
	msg.Image = db.Image(imgRes.Data.ImgMsg)
	if err := db.Conn.Save(&msg).Error; err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	utils.ResOk(c)
}

func DelMsg(c *gin.Context) {
	u := c.MustGet("account").(*db.SphAccount)
	msgid := c.MustGet("msgid").(int64)
	if err := db.Conn.Delete(&db.SphMsg{}, "id = ?", msgid).Error; err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	if err := updateMsg(u); err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	utils.ResOk(c)
}

func updateMsg(account *db.SphAccount) error {
	var msgs []db.SphMsg
	if err := db.Conn.Find(&msgs).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}
	_, err := wx.New(account, &msgs)
	if err != nil {
		return err
	}
	return nil
}
