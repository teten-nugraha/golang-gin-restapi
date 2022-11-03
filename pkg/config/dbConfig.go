package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/teten-nugraha/golang-gin-project/pkg/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func SetupDBConnection() *gorm.DB {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln("Failed to load env")
	}

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to create a connection to database")
	}

	log.Println("Connected to database")

	// migrate entities
	db.AutoMigrate(&entity.Book{}, &entity.User{})
	return db
}

func CloseDBConnection(db *gorm.DB) {
	dbSQL, err := db.DB()
	if err != nil {
		panic("Failed to close connection from db")
	}
	dbSQL.Close()
}
