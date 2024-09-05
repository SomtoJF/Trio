package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
	"golang.org/x/crypto/bcrypt"
)

type loginInput struct {
	Username string `json:"userName" binding:"required,max=20"`
	Password string `json:"password" binding:"required,max=20"`
}

func Login(c *gin.Context) {
	domain := os.Getenv("DOMAIN")
	var body loginInput

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userFound models.User
	initializers.DB.Where("username=?", body.Username).Find(&userFound)

	if userFound.ID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.PasswordHash), []byte(body.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  userFound.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to generate token"})
	}

	c.SetCookie("Access_Token", token, 315360000, "/", domain, false, true)

	c.JSON(200, gin.H{
		"message": "Login successful",
	})
}
