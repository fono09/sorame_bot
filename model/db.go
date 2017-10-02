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
    sdsn := strings.Split(config.DB, "/")
    dsn := sdsn[0] + "/"
    dbn := sdsn[1]

    mysql, err := sql.Open("mysql", dsn)
    if err != nil {
        return err
    }
    defer mysql.Close()

    _, err = mysql.Exec("CREATE DATABASE IF NOT EXISTS ?;", dbn)
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
