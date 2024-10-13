package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// Logout godoc
//	@Summary		Logout user
//	@Description	Logs out the user by clearing the access token
//	@Tags			auth
//	@Success		200	{object}	map[string]interface{}	"Logout successful"
//	@Router			/logout [post]
func Logout(c *gin.Context) {
	c.SetCookie("Access_Token", "", -1, "/", os.Getenv("DOMAIN"), false, true)
	c.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
