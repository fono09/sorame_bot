package model

import (
	//	"fmt"
	//"github.com/VG-Tech-Dojo/treasure2017/mid/fono09/sorame_bot/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	//	"log"
)

func Init(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Sorame{})
}
