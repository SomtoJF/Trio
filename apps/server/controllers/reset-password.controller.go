package controllers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
	"golang.org/x/crypto/bcrypt"
)

type passwordResetRequest struct {
	Password    string `json:"password" binding:"required,max=20"`
	NewPassword string `json:"newPassword" binding:"required,max=20"`
}

// ResetPassword godoc
//	@Summary		Reset user password
//	@Description	Resets the password for the authenticated user
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			passwordResetRequest	body		passwordResetRequest	true	"Password reset details"
//	@Success		200						{object}	map[string]interface{}	"Password updated successfully"
//	@Failure		400						{object}	map[string]interface{}	"Bad request"
//	@Failure		401						{object}	map[string]interface{}	"Unauthorized"
//	@Failure		500						{object}	map[string]interface{}	"Internal server error"
//	@Router			/reset-password [post]
func ResetPassword(c *gin.Context) {
	var body passwordResetRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the current user from the context
	user, ok := c.Value("currentUser").(models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Verify the current password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	// Update the password in the database
	user.PasswordHash = string(hashedPassword)
	if err := initializers.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.SetCookie("Access_Token", "", -1, "/", os.Getenv("DOMAIN"), false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully. Please login with new password"})
}
