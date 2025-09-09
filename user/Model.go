package user

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

var jwtKey = []byte("19910405fanpy")

type User struct {
	gorm.Model `json:"grom_._model,omitempty" :"gorm_._model"`
	Username   string `json:"user_name,omitempty" :"user_name"`
	Password   string `json:"password,omitempty" :"password"`
	Email      string
}
type Post struct {
	gorm.Model `json:"gorm_._model,omitempty"`
	Title      *string `json:"title,omitempty"`
	Content    *string `json:"total,omitempty"`
	UserId     uint
}
type Comments struct {
	gorm.Model
	Content *string
	UserId  uint
	PostId  uint
}

type Claims struct {
	Username string `json:"username"`
	UserId   uint   `json:"userId"`
	jwt.RegisteredClaims
}

type LogEntry struct {
	ID        uint      `gorm:"primaryKey"`
	Level     string    `gorm:"not null"`
	Message   string    `gorm:"type:text;not null"`
	Timestamp time.Time `gorm:"autoCreateTime"`
}
