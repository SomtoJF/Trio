package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Signup(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"message": "account created successfully",
	})
}
