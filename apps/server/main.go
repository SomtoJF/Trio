package main

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/somtojf/trio/controllers"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func main() {
	r := gin.Default()
	clientAddress := os.Getenv("CLIENT_ADDRESS")

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{clientAddress}
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	r.Use(cors.New(config))

	public := r.Group("/")
	{
		public.POST("/login", controllers.Login)
		public.POST("/signup", controllers.Signup)
	}

	authenticated := r.Group("/")
	authenticated.Use(middleware.CheckAuth())
	{
		authenticated.GET("/me", controllers.GetCurrentUser)
	}

	r.Run()
}
