package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func Init(db *gorm.DB) {
	db.AutoMigrate(&User{}, &Sorame{})
}
