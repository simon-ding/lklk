package models

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)
import _ "github.com/jinzhu/gorm/dialects/sqlite"

func OpenDB() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./config/data.db")
	if err != nil {
		panic(err)
	}
	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(&TVSeries{}, &Episode{}, &Task{})
}

func GetDB(c *gin.Context) *gorm.DB {
	return c.MustGet("db").(*gorm.DB)
}