package main

import (
	"wx_video_help/config"
	"wx_video_help/db"
	"wx_video_help/server"
	"wx_video_help/wx"
)

func main() {
	config.Load()
	db.InitSqlLite()
	wx.Run()
	server.Run()
}
