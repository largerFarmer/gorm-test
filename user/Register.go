package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 注册
func Register(c *gin.Context) {

	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		//打印错误日志
		c.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		return

	}
	//密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return
	}
	user.Password = string(hashedPassword)
	db, err := GetDBConn()

	//这里限制了用户名和邮箱唯一
	if err := db.Create(&user).Error; err != nil {
		RespondWithError(db, c, http.StatusBadRequest, "用户名或邮箱已存在", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"msg": "注册成功"})

}
