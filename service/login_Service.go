package service

import (
	"LeakInfo/bean"
	"LeakInfo/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

// Đăng ký người dùng
func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user bean.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
			return
		}

		if user.Username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Username không được để trống"})
			return
		}

		user.Status = 1
		user.Role = "user"

		// Mã hóa mật khẩu
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi mã hóa mật khẩu"})
			return
		}
		user.Password = string(hashedPassword)

		// Lưu vào DB
		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi tạo tài khoản"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Đăng ký thành công"})
	}
}

// Đăng nhập và tạo token
func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input bean.User
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu không hợp lệ"})
			return
		}

		var user bean.User
		if err := db.Where("username = ? and status = 1 ", input.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Tài khoản không tồn tại"})
			return
		}

		// Kiểm tra mật khẩu
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Sai mật khẩu"})
			return
		}

		tokenString, err := utils.GenerateJWTToken(user.ID, user.Username, user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi tạo token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"token": tokenString})
	}
}
