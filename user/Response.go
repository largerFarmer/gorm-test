package user

import (
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RespondWithError(db *gorm.DB, c *gin.Context, code int, msg string, err error) {
	if err != nil {
		LogToDB(db, "error", msg+": "+err.Error())
	} else {
		LogToDB(db, "error", msg)
	}
	c.JSON(code, gin.H{"error": msg})
}

func LogToDB(db *gorm.DB, level, msg string) {
	log := LogEntry{
		Level:     level,
		Message:   msg,
		Timestamp: time.Now(),
	}
	db.Create(&log)
}
