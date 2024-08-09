package utils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ResOkWithData(c *gin.Context, data interface{}) {
	ResOkWithAll(c, "成功", data)
}
func ResOkWithMsg(c *gin.Context, msg string) {
	ResOkWithAll(c, msg, nil)
}

func ResOkWithAll(c *gin.Context, msg string, data interface{}) {
	Res(c, 0, msg, data)
}
func ResOk(c *gin.Context) {
	ResOkWithMsg(c, "成功")
}
func Res(c *gin.Context, status int, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": status,
		"msg":    msg,
		"data":   data,
	})
}

func ResErrWithMsg(c *gin.Context, msg string) {
	ResErrWithAll(c, msg, nil)
}
func ResErrWithAll(c *gin.Context, msg string, data interface{}) {
	Res(c, 404, msg, data)
}
func ResErr(c *gin.Context) {
	ResErrWithMsg(c, "失败")
}
