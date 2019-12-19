package database

import (
	"github.com/jinzhu/gorm"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db * gorm.DB

func init(){
	var DBNAME = "mydb"
	var DBUSERNAME = "root"
	var DBPASSWORD = "1234"
	var DBADDRESS = "localhost"
	var DBPORT = "3306"
	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8", DBUSERNAME, DBPASSWORD, DBADDRESS, DBPORT, DBNAME)
	var err error
	db, err = gorm.Open("mysql", url)
	if err != nil {
		log.Println(err)
		return
	}
}