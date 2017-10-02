package model

import (
	"database/sql"
	"github.com/VG-Tech-Dojo/treasure2017/mid/fono09/sorame_bot/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *sql.DB
var gdb *gorm.DB

func InitDB() {
	var err error
	gdb, err = gorm.Open("mysql", config.DB)
	db = gdb.DB()
	if err != nil {
		panic(err)
	}
	gdb.AutoMigrate(&User{}, &Sorame{})
}
