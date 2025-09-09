package user

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//func Init() *gorm.DB {
//	db, err := gorm.Open(mysql.Open("root:root@123@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"))
//
//	if err != nil {
//		panic("数据库连接失败")
//	}
//
//	return db
//}

func GetDBConn() (*gorm.DB, error) {

	db, err := gorm.Open(mysql.Open("root:root@123@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"))

	if err != nil {
		panic("数据库连接失败")
	}
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Comments{})
	db.AutoMigrate(&Post{})
	return db, err
}
