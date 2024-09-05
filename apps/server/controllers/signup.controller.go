package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
	"golang.org/x/crypto/bcrypt"
)

type SignUpInput struct {
	Username string `json:"userName" binding:"required,max=20"`
	FullName string `json:"fullName" binding:"required,max=50"`
	Password string `json:"password" binding:"required,max=20"`
}

func Signup(c *gin.Context) {
	var body SignUpInput

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	initializers.DB.Where("username=?", body.Username).Find(&userFound)

	if userFound.ID != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username taken"})
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Username:     body.Username,
		FullName:     body.FullName,
		PasswordHash: string(passwordHash),
	}

	initializers.DB.Create(&user)

	c.JSON(http.StatusCreated, gin.H{
		"message": "account created successfully",
		"data":    body,
	})
}
