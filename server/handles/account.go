package handles

import (
	"wx_video_help/db"
	"wx_video_help/utils"
	"wx_video_help/wx"

	"github.com/gin-gonic/gin"
)

func Save(c *gin.Context) {
	u := c.MustGet("account").(*db.SphAccount)
	var params struct {
		TimeSleep int `json:"timeSleep"`
	}
	if err := c.ShouldBind(&params); err != nil {
		utils.ResErrWithMsg(c, "参数错误，"+err.Error())
		return
	}

	u.TimeSleep = params.TimeSleep
	if err := db.Conn.Save(&u).Error; err != nil {
		utils.ResErrWithMsg(c, "保存失败，"+err.Error())
		return
	}

	index := wx.GetUserIndex(u.ID)
	if index != -1 {
		if _, err := wx.New(u, nil); err != nil {
			utils.ResErrWithMsg(c, err.Error())
			return
		}
	}
	utils.ResOk(c)
}
func GetAccountInfo(c *gin.Context) {
	u := c.MustGet("account").(*db.SphAccount)
	if u.Switch {
		index := wx.GetUserIndex(u.ID)
		if index == -1 {
			u.Switch = false
		} else {
			w, _ := wx.New(u, nil)
			if w.Close {
				u.Switch = false
			}
		}
	}
	var msgs []db.SphMsg
	if err := db.Conn.Find(&msgs, "aid = ?", u.ID).Error; err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	utils.ResOkWithData(c, gin.H{
		"info": u,
		"msgs": msgs,
	})
}

func SetServer(c *gin.Context) {
	u := c.MustGet("account").(*db.SphAccount)
	var params struct {
		Status bool `json:"status"`
	}
	if err := c.ShouldBind(&params); err != nil {
		utils.ResErrWithMsg(c, "参数错误，"+err.Error())
		return
	}
	if params.Status {
		var msgs []db.SphMsg
		if err := db.Conn.Find(&msgs).Error; err != nil {
			utils.ResErrWithMsg(c, err.Error())
			return
		}
		if err := wx.Add(u, &msgs); err != nil {
			utils.ResErrWithMsg(c, err.Error())
			return
		}
		utils.ResOkWithMsg(c, "运行成功")
	} else {
		if err := wx.Stop(u); err != nil {
			utils.ResErrWithMsg(c, err.Error())
			return
		}
		utils.ResOkWithMsg(c, "关闭成功")
	}
	db.Conn.Model(&db.SphAccount{}).Where("id = ?", u.ID).Update("switch", params.Status)
}

func GetAccountList(c *gin.Context) {
	var params struct {
		Page     int    `json:"page"`
		PerPage  int    `json:"perPage"`
		Sph      string `json:"sph"`
		NickName string `json:"nickname"`
	}
	if err := c.ShouldBind(&params); err != nil {
		utils.ResErrWithMsg(c, "参数错误，"+err.Error())
		return
	}
	user := c.MustGet("user").(*db.SphUser)
	var accounts []db.SphAccount
	var total int64 = 0
	if params.Sph != "" {
		if err := db.Conn.Find(&accounts, "uniqid like ?", params.Sph).Error; err != nil {
			utils.ResErrWithMsg(c, err.Error())
			return
		}
	} else if params.NickName != "" {
		if err := db.Conn.Find(&accounts, "nick_name like %?%", params.NickName).Error; err != nil {
			utils.ResErrWithMsg(c, err.Error())
			return
		}
	} else {
		var as []db.SphAccount
		if err := db.Conn.Find(&as).Count(&total).Error; err != nil {
			utils.ResErrWithMsg(c, err.Error())
			return
		}
		offset := (params.Page - 1) * params.PerPage
		if err := db.Conn.Limit(params.PerPage).Offset(offset).Find(&accounts, "uid = ?", user.ID).Error; err != nil {
			utils.ResErrWithMsg(c, err.Error())
			return
		}
	}

	var panels []*PanelItem
	for _, v := range accounts {
		err := "未知"
		if !v.Switch {
			err = "未运行"
		}
		panels = append(panels, &PanelItem{
			ID:       v.ID,
			NickName: v.NickName,
			Close:    true,
			Error:    err,
			Count:    0,
		})
	}

	l := wx.List()
	if len(l) > 0 {
		for _, v := range l {
			err := ""
			for _, p := range panels {
				if v.ID == p.ID {
					if v.Err != nil {
						err = v.Err.Error()
					}
					p.Close = v.Close
					p.Error = err
					p.Count = v.Count
					break
				}
			}
		}
	}
	utils.ResOkWithData(c, gin.H{
		"rows":  panels,
		"total": total,
	})
}
func DelAccount(c *gin.Context) {
	account := c.MustGet("account").(*db.SphAccount)
	if err := wx.Stop(account); err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	if err := db.Conn.Delete(&db.SphAccount{}, "id = ?", account.ID).Error; err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	if err := db.Conn.Delete(&db.SphMsg{}, "aid = ?", account.ID).Error; err != nil {
		utils.ResErrWithMsg(c, err.Error())
		return
	}
	utils.ResOk(c)
}

type PanelItem struct {
	ID       int64  `json:"id"`
	NickName string `json:"nickname"`
	Close    bool   `json:"close"`
	Error    string `json:"error"`
	Count    int    `json:"count"`
}
