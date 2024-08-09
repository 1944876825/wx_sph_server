package server

import (
	"bytes"
	"fmt"
	"io/fs"
	"net/http"
	"wx_video_help/config"
	"wx_video_help/public"
	"wx_video_help/server/handles"
	"wx_video_help/server/middlewares"

	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 禁用 CORS
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*") // 允许所有域名访问
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	r.Any("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	dist, err := fs.Sub(public.Public, "dist/assets")
	if err != nil {
		panic("can't find folder: dist")
	}
	r.StaticFS("/assets", http.FS(dist))

	RawIndexHtml, err := public.Public.ReadFile("dist/index.html")
	if err != nil {
		panic("can't find folder: dist")
	}
	replayMap := map[string]string{
		"Vite + React + TS": config.Conf.Title,
		"./vite.svg":        config.Conf.Icon,
	}
	for k, v := range replayMap {
		if v != "" {
			RawIndexHtml = bytes.Replace(RawIndexHtml, []byte(k), []byte(v), 1)
		}
	}
	r.GET("/", func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", RawIndexHtml)
	})
	setWxRouter(r)

	fmt.Println("程序启动成功: ", fmt.Sprintf("http://127.0.0.1:%d", config.Conf.Port))
	if err := r.Run(fmt.Sprintf(":%d", config.Conf.Port)); err != nil {
		panic("web服务启动失败，" + err.Error())
	}
}

var RawIndexHtml []byte

func setWxRouter(r *gin.Engine) {
	r.POST("/getLoginUrl", handles.GetLoginUrl)
	r.GET("/loginStatus", handles.GetLoginStatus)
	r.POST("/login", middlewares.GetUser, handles.Login)

	r.POST("/addMsg", middlewares.AuthAccount, handles.AddMsg)
	r.POST("/msgInfo", middlewares.AuthMsg, handles.GetMsgInfo)
	r.POST("/saveMsg", middlewares.AuthMsg, handles.SaveMsg)
	r.POST("/img", middlewares.AuthMsg, handles.GetImg)
	r.POST("/upload", middlewares.AuthMsg, handles.UploadImg)
	r.POST("/delMsg", middlewares.AuthMsg, handles.DelMsg)

	r.POST("/save", middlewares.AuthAccount, handles.Save)
	r.POST("checkCookie", middlewares.AuthAccount, handles.CheckCookie)
	r.POST("/info", middlewares.AuthAccount, handles.GetAccountInfo)
	r.POST("/setServer", middlewares.AuthAccount, handles.SetServer)
	r.POST("/account/list", middlewares.AuthUser, handles.GetAccountList)
	r.POST("/delAccount", middlewares.AuthAccount, handles.DelAccount)
}
