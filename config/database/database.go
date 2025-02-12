package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	server "bbscout/config"
)

var db *gorm.DB

func NewDatabaseConnection() {

	if db != nil {
		return
	}

	dbHost := server.Envs.DBHost
	dbPort := server.Envs.DBPort
	dbUser := server.Envs.DBUser
	dbPassword := server.Envs.DBPassword
	dbName := server.Envs.DBName

	// dsn := "afyarekod_user:4bUC4v%246xFM7hd8@tcp(127.0.0.1:32771)/data_platform_db?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser, dbPassword, dbHost, dbPort, dbName)
	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db = database

	fmt.Println("Finished setting database successfully")

}

func GetDB() *gorm.DB {
	return db
}
