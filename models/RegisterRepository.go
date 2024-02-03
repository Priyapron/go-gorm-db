// models/RegisterRepository.go
package models

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterRepository struct {
	Db *gorm.DB
}

func NewRegisterRepository(db *gorm.DB) *RegisterRepository {
	return &RegisterRepository{Db: db}
}

func (r *RegisterRepository) RegisterUser(c *gin.Context) {
	var newUser Register
	c.BindJSON(&newUser)

	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" {
		c.JSON(400, gin.H{"message": "Name, Email, and Password are required"})
		return
	}

	existingUser := Register{}
	result := r.Db.Where("email = ?", newUser.Email).First(&existingUser)
	if result.RowsAffected > 0 {
		c.JSON(400, gin.H{"message": "Email already registered"})
		return
	}

	newUser.Hash = GeneratePasswordHash(newUser.Password)

	if err := r.Db.Create(&newUser).Error; err != nil {
		c.JSON(500, gin.H{"message": "Failed to register user"})
		return
	}

	newUser.Password = ""
	c.JSON(201, newUser)
}

func (r *RegisterRepository) AuthenticateUser(email, password string) bool {
	var user Register
	result := r.Db.Where("email = ?", email).First(&user)

	if result.RowsAffected == 0 {
		return false
	}

	return CheckPasswordHash(password, user.Hash)
}
