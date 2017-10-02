package model

import (
    "strings"
	"database/sql"
	"github.com/fono09/sorame_bot/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var db *sql.DB
var gdb *gorm.DB

func createTable() error {
    mysql, err := sql.Open("mysql", strings.Split(config.DB, "/")[0] + "/")
    if err != nil {
        return err
    }
    defer mysql.Close()

    _, err = mysql.Exec("CREATE DATABASE ? IF NOT EXISTS", strings.Split(config.DB, "/")[1])
    if err != nil {
        return err
    }
    return nil
}

func InitDB() {
    err := createTable()
    if err != nil {
        panic(err)
    }

	gdb, err = gorm.Open("mysql", config.DB)
	db = gdb.DB()
	if err != nil {
		panic(err)
	}
	gdb.AutoMigrate(&User{}, &Sorame{})
}
