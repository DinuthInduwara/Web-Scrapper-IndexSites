package db

import (
	"errors"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var Database = connectDB()

type TFile struct {
	gorm.Model
	Url  string
	Size int64
}

func connectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&TFile{})
	return db
}

func IsinDb(file *TFile) bool {
	var tUrl TFile
	result := Database.Where("url = ? AND size = ?", file.Url, file.Size).First(&tUrl)

	if err := result.Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false
		}
		// Handle other errors
		panic(err)
		return false
	}

	return true
}

func AddFileToDb(file *TFile) error {
	result := Database.Create(&file)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
