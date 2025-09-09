package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {

	var req User
	if err := c.ShouldBindJSON(&req); err != nil {
		//打印错误日志
		c.JSON(http.StatusBadRequest, gin.H{"err": "参数错误"})
		return
	}

	db, err := GetDBConn()
	var user User
	if err != nil {
		RespondWithError(db, c, http.StatusUnauthorized, "数据库连接失败", err)
	}
	if err := db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		RespondWithError(db, c, http.StatusUnauthorized, "用户不存在", err)
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		RespondWithError(db, c, http.StatusUnauthorized, "密码错误", err)

		return
	}
	// 生成 JWT

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	c.JSON(http.StatusOK, gin.H{"token": tokenString})

}

// JwtToken jwt中间件
func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization") // 获取请求头中的Authorization字段
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header binding failed"})
			c.Abort() // 中断请求处理流程
			return
		}

		// 解析Token并验证其有效性
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证Token的方法类型是否是我们所期望的HS256（HMAC SHA256）
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtKey, nil // 返回用于签名Token的密钥
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort() // 中断请求处理流程
			return
		}

		if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid { // 检查Token是否有效并且包含Claims信息
			// 将Claims信息存储在Context中，供后续处理使用
			//c.Set("claims", claims)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			c.Abort() // 中断请求处理流程
		}

		c.Next() // 继续执行后续的处理函数
	}
}
func ParseAccessToken(accessToken string) (*Claims, error) {
	parsedAccessToken, _ := jwt.ParseWithClaims(accessToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if claim, ok := parsedAccessToken.Claims.(*Claims); ok && parsedAccessToken.Valid {
		return claim, nil
	}
	return nil, fmt.Errorf("parsedAccessToken return invalid")
}

func ParseRefreshToken(refreshToken string) *Claims {
	parsedRefreshToken, _ := jwt.ParseWithClaims(refreshToken, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	return parsedRefreshToken.Claims.(*Claims)
}
