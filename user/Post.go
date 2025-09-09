package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreatePost(c *gin.Context) {
	var postInfo Post
	if err := c.ShouldBindJSON(&postInfo); err != nil {
		fmt.Printf("ShouldBindJSON err postInfo: %v\n", postInfo)
		c.JSON(http.StatusBadRequest, gin.H{"error Post": err.Error()})
		return
	}
	fmt.Printf("postInfo: %v\n", postInfo)

	if postInfo.UserId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UserID can not be null"})
	}
	if postInfo.Title == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Title can not be null"})
	}
	if postInfo.Content == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Content can not be null"})
	}

	tokenString := c.GetHeader("Authorization")
	claims, ok := ParseAccessToken(tokenString)
	if ok != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ParseAccessToken error"})
	}
	if claims.UserId != postInfo.UserId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token user not"})
	}
	var storedUser User
	db, err := GetDBConn()
	if err != nil {
		fmt.Printf("Failed to get DB connection: %v\n", err)
		return
	}
	if err := db.First(&storedUser, postInfo.UserId).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username "})
		return
	}
	//log.Fatal("postInfo.Content=: ", postInfo.Content)
	fmt.Println("postInfo.Content=", postInfo.Content)
	post := Post{Title: postInfo.Title, Content: postInfo.Content, UserId: storedUser.ID}
	if err := db.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Post"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Post successfully"})
}

func ListPosts(c *gin.Context) {
	var postId uint

	db, err := GetDBConn()
	if err != nil {
		fmt.Printf("Failed to get DB connection: %v\n", err)
		return
	}
	if postId == 0 {

		var results []Post
		err := db.Find(&results).Error
		if err != nil {
			// 处理错误
			c.JSON(http.StatusUnauthorized, gin.H{"error": "batch queryPost err"})
			return
		} else {
			c.JSON(http.StatusOK, results)
		}
	} else {
		var postRet Post
		err := db.Find(&postRet, postId).Error
		if err != nil {
			// 处理错误
			c.JSON(http.StatusUnauthorized, gin.H{"error": "queryPost one err"})
			return
		} else {
			c.JSON(http.StatusOK, postRet)
		}
	}

}

func UpdatePost(c *gin.Context) {
	var postInfo Post
	if err := c.ShouldBindJSON(&postInfo); err != nil {
		fmt.Printf("ShouldBindJSON err postInfo: %v\n", postInfo)
		c.JSON(http.StatusBadRequest, gin.H{"error Post": err.Error()})
		return
	}
	fmt.Printf("postInfo: %v\n", postInfo)
	if postInfo.ID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Post ID can not be null"})
	}
	if postInfo.UserId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "UserID can not be null"})
	}
	if postInfo.Title == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Title can not be null"})
	}
	if postInfo.Content == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Content can not be null"})
	}

	tokenString := c.GetHeader("Authorization")
	claims, ok := ParseAccessToken(tokenString)
	if ok != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "ParseAccessToken error"})
	}
	if claims.UserId != postInfo.UserId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "token user not"})
	}

	db, err := GetDBConn()
	if err != nil {
		//fmt.Printf("Failed to get DB connection: %v\n", err)
		c.JSON(http.StatusInternalServerError, "Failed to get DB connection")
		return
	}
	var postDst Post
	result := db.Find(&postDst, postInfo.ID)
	if result.RowsAffected != 1 {
		// 处理错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "query postInfo.ID err"})
		return
	} else {

		result := db.Model(&postDst).Where("id = ?", postDst.ID).Select("Title", "Content").Updates(Post{Title: postInfo.Title, Content: postInfo.Content})
		if result.RowsAffected == 1 {
			c.JSON(http.StatusOK, "Post Update OK")
		} else {
			c.JSON(http.StatusBadRequest, "Post Update FAIL")
		}

	}

}
