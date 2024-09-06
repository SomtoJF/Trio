package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	aihelpers "github.com/somtojf/trio/ai-helpers"
	"github.com/somtojf/trio/models"
	"github.com/somtojf/trio/types"
)

type completionsRequest struct {
	Text string `json:"text" binding:"required"`
}

func GetCompletion(c *gin.Context) {
	var body completionsRequest
	user, ok := c.Value("currentUser").(models.User)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	completionsRequest := types.GeminiCompletionsRequest{
		Prompt:     body.Text,
		SenderID:   user.ID,
		SenderType: types.SenderTypeUser,
	}

	resp, err := aihelpers.GetGeminiCompletions(c, completionsRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": resp.Candidates[0].Content.Parts[0]})
}
