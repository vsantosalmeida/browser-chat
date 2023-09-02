package config

import (
	"github.com/vsantosalmeida/browser-chat/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "root:3luproec@tcp(127.0.0.1:3306)/chat?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	if err = db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&entity.User{}); err != nil {
		panic(err.Error())
	}

	return db
}
