package db

import (
	"errors"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Database = connectDB()

type TFile struct {
	gorm.Model
	Url  string `gorm:"index"`
	Size int64
}

func connectDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{Logger: silentLogger})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&TFile{})
	return db
}

func IsinDb(file *TFile) bool {
	var tUrl TFile
	// retries := 0

	for {
		result := Database.Where("url = ? AND size = ?", file.Url, file.Size).First(&tUrl)

		if err := result.Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false
			} else if strings.Contains(err.Error(), "SQLITE_BUSY") {
				time.Sleep(time.Second)
				// retries += 1
				// log.Println("[E] DataBase Connection Error: Retry: ", retries)
				continue
			}
			panic(err)
		}

		return true
	}

}

func IsUrlInDatabase(url string) bool {
	var tUrl TFile
	// retries := 0

	for {
		result := Database.Where("url = ?", url).First(&tUrl)

		if err := result.Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return false
			} else if strings.Contains(err.Error(), "SQLITE_BUSY") {
				time.Sleep(time.Second)
				// retries += 1
				// log.Println("[E] DataBase Connection Error: Retry: ", retries)
				continue
			}
			panic(err)
		}

		return true
	}

}

func AddFileToDb(file *TFile) error {
	// retries := 0

	for {
		result := Database.Create(&file)

		if err := result.Error; err != nil {
			if strings.Contains(err.Error(), "SQLITE_BUSY") {
				time.Sleep(time.Second)
				// retries += 1
				// log.Println("[E] DataBase Connection Error: Retry: ", retries)
				continue
			}
			panic(err)
		}

		return nil
	}

}

var silentLogger = logger.New(
	nil, // Use your preferred io.Writer if you want to log somewhere
	logger.Config{
		LogLevel: logger.Silent,
	},
)

func DocumentsCount() int64 {
	var documentCount int64
	Database.Model(&TFile{}).Count(&documentCount)
	return documentCount
}
