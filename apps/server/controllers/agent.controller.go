package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/models"
	"gorm.io/gorm"
)

// DeleteAgent deletes an agent
func DeleteAgent(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("agentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID"})
		return
	}

	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userModel := currentUser.(models.User)

	var agent models.Agent
	if err := initializers.DB.Joins("Chat").
		Where("agents.external_id = ? AND chats.user_id = ?", agentID, userModel.ID).
		First(&agent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found or does not belong to the user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve agent"})
		}
		return
	}

	if err := initializers.DB.Delete(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete agent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Agent deleted successfully"})
}

// GetAgent retrieves an agent's details
func GetAgent(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("agentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID"})
		return
	}

	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userModel := currentUser.(models.User)

	var agent models.Agent
	if err := initializers.DB.Joins("Chat").
		Where("agents.external_id = ? AND chats.user_id = ?", agentID, userModel.ID).
		First(&agent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found or does not belong to the user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve agent"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"agent": agent})
}

// UpdateAgent updates an agent's details
func UpdateAgent(c *gin.Context) {
	agentID, err := uuid.Parse(c.Param("agentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid agent ID"})
		return
	}

	var body struct {
		Name   string   `json:"name" binding:"required,max=20"`
		Lingo  string   `json:"lingo" binding:"required,max=20"`
		Traits []string `json:"traits" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	currentUser, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userModel := currentUser.(models.User)

	var agent models.Agent
	if err := initializers.DB.Joins("Chat").
		Where("agents.external_id = ? AND chats.user_id = ?", agentID, userModel.ID).
		First(&agent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Agent not found or does not belong to the user"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve agent"})
		}
		return
	}

	agent.Name = body.Name
	agent.Lingo = body.Lingo
	agent.Traits = body.Traits

	if err := initializers.DB.Save(&agent).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update agent"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"agent": agent})
}
