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
	Password string `json:"password" binding:"required,max=20,min=8"`
}

// Login godoc
//	@Summary		Login user
//	@Description	Logs in a user and returns an access token
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			loginInput	body		loginInput				true	"Login credentials"
//	@Success		200			{object}	map[string]interface{}	"success message"
//	@Failure		400			{object}	map[string]interface{}	"error message"
//	@Failure		500			{object}	map[string]interface{}	"internal server error"
//	@Router			/login [post]
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userFound.PasswordHash), []byte(body.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username or password is incorrect"})
		return
	}

	generateToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userFound.ExternalID.String(),
		"username": userFound.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	token, err := generateToken.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "An error occured"})
	}

	c.SetCookie("Access_Token", token, 604800, "/", domain, false, true)

	c.JSON(200, gin.H{
		"message": "success",
	})
}
