package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Logout(c *gin.Context) {
	c.SetCookie("Access_Token", "", -1, "/", os.Getenv("DOMAIN"), false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
