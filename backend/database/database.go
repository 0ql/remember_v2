package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Db  *gorm.DB
	err error
)

type Question struct {
	gorm.Model
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Wrong0   string `json:"wrong0"`
	Wrong1   string `json:"wrong1"`
	Wrong2   string `json:"wrong2"`
}

type User struct {
	gorm.Model
	Email        string `json:"email"`
	Username     string `json:"username"`
	Pwdhash      string `json:"pwdhash"`
	Sessiontoken string `json:"sessiontoken"`
	Expires      int64  `json:"expires"`
	Rank         string `json:"rank"`
	Elo          int    `json:"elo"`
}

func init() {
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,          // Disable color
		},
	)

	Db, err = gorm.Open(sqlite.Open("main.db"), &gorm.Config{
		Logger: newLogger,
	})

	if err != nil {
		panic(err)
	}

	Db.AutoMigrate(&Question{})
	Db.AutoMigrate(&User{})
}
