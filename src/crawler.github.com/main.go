package main

import (
	"fmt"
	"log"

	"crawler.github.com/wandoujia"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// mysql配置
const (
	username = ""
	password = ""
	dbname   = ""
)

func main() {
	db, err := gorm.Open("mysql", fmt.Sprintf("%s:%s@/%s?charset=utf8&parseTime=True&loc=Local", username, password, dbname))
	if err != nil {
		log.Panic("open database error:", err)
	}
	// db.CreateTable(models.APP{})
	defer db.Close()
	wandoujia.Work(db)
}
