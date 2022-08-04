package main

import (
	"fmt"
	"github.com/SeregaSergo/Go_intensive/team_0/internal/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

const (
	DB_HOST     = "localhost"
	DB_PORT     = 5432
	DB_USER     = "postgres"
	DB_PASSWORD = "kek"
	DB_NAME     = "serega"
)

func main() {
	var frequencies []storage.FrequencyRecord
	db := connectDB()
	result := db.Find(&frequencies)
	fmt.Printf("Returned %d rows", result.RowsAffected)
	for i, v := range frequencies {
		fmt.Println(i, v.Frequency)
	}
}

func connectDB() *gorm.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Println("engine creation failed", err)
	}
	dbSql, _ := db.DB()
	err = dbSql.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected")
	db.AutoMigrate(&storage.FrequencyRecord{})
	return db
}
