package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

var host = os.Getenv("HOST")
var port = os.Getenv("PORT")
var user = os.Getenv("USER")
var password = os.Getenv("PASSWORD")
var dbname = os.Getenv("DBNAME")
var sslmode = os.Getenv("SSLMODE")

var dbInfo = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, user, password, dbname, sslmode)
var Db *gorm.DB

func Boot() (*gorm.DB, error) {
	fmt.Println("--- CALL BOOT FUNCTION ---")
	fmt.Println(dbInfo)

	db, err := gorm.Open(postgres.Open(dbInfo), &gorm.Config{})

	if err != nil {
		fmt.Println("ERROR CONNECTION")
		return nil, err
	}

	Db = db

	//Создаем таблицу
	if os.Getenv("CREATE_TABLE") == "yes" {

		if err := db.AutoMigrate(&User{}); err != nil {
			fmt.Println("GORM ERROR")
			fmt.Println(err)
		}
	}

	return db, nil
}
