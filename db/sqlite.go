package db

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitSqlLite() {
	db, err := gorm.Open(sqlite.Open("wx.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	Conn = db
	createDb()
}

func createDb() {
	err := Conn.AutoMigrate(&SphUser{}, &SphAccount{}, &SphMsg{})
	if err != nil {
		panic(err)
	}
}
