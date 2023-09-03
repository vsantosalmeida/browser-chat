package config

import (
	"fmt"

	"github.com/vsantosalmeida/browser-chat/entity"

	"github.com/apex/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const dsnPattern = "%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local"

// InitDB create the connection with DB and migrate the user, room and message tables.
func InitDB() *gorm.DB {
	dsn := fmt.Sprintf(
		dsnPattern,
		GetStingEnvVarOrPanic(MySqlUser),
		GetStingEnvVarOrPanic(MySqlPass),
		GetStingEnvVarOrPanic(MySqlHost),
		GetStingEnvVarOrPanic(MySqlDB),
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.WithError(err).Fatal("failed to open db connection")
	}

	if err = db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&entity.User{}); err != nil {
		log.WithError(err).Fatal("failed to migrate user table")
	}

	if err = db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&entity.Room{}); err != nil {
		log.WithError(err).Fatal("failed to migrate room table")
	}

	if err = db.Set("gorm:table_options", "ENGINE=InnoDB").AutoMigrate(&entity.Message{}); err != nil {
		log.WithError(err).Fatal("failed to migrate message table")
	}

	return db
}
