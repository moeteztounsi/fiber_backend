package database

import (
	"log"

	"github.com/alpha_batta/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstence struct {
	Db *gorm.DB
}

var Database DbInstence

func Connect() {

	dsn := "root:25865406mT?@tcp(127.0.0.1:3306)/batta?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Could not connect to DataBase")
	}

	log.Println("Connected to the database successfuly")

	db.Logger = logger.Default.LogMode(logger.Info)

	log.Println("Running migrations")
	// here we're gonna run our migrations to create tables in the database
	db.AutoMigrate(&models.User{})

	Database = DbInstence{Db: db}

}
