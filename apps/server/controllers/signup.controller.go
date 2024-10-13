package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
	"golang.org/x/crypto/bcrypt"
)

type signUpInput struct {
	Username string `json:"userName" binding:"required,max=20"`
	FullName string `json:"fullName" binding:"required,max=50"`
	Password string `json:"password" binding:"required,max=20,min=8"`
}

// Signup godoc
//	@Summary		Signup a new user
//	@Description	Creates a new user account
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			userInput	body		signUpInput				true	"User details"
//	@Success		201			{object}	map[string]interface{}	"Account created successfully"
//	@Failure		400			{object}	map[string]interface{}	"Bad request"
//	@Failure		500			{object}	map[string]interface{}	"Internal server error"
//	@Router			/signup [post]
func Signup(c *gin.Context) {
	var body signUpInput

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
