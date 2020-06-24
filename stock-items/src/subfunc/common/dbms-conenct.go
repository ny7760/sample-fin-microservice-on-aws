package common

import (
	"fmt"
	"log"
	"os"

	// for mysql
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// ConnectDB ... DB接続メソッド
func ConnectDB(mode string) *gorm.DB {

	// Read DBMS Enviroment Variables
	Dbms := os.Getenv("DBMS")
	UserName := os.Getenv("DB_USERNAME")
	PassWord := os.Getenv("DB_PASSWORD")
	var Host string
	if mode == "read" {
		Host = os.Getenv("DB_HOST_READ")
	} else {
		Host = os.Getenv("DB_HOST_MASTER")
	}
	DBName := os.Getenv("DB_NAME")

	CONNECT := UserName + ":" + PassWord + "@" + Host + "/" + DBName + "?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(Dbms, CONNECT)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully Connected!")
	}
	return db
}
