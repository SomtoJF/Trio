package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/generative-ai-go/genai"
	"github.com/somtojf/trio/clients"
	"github.com/somtojf/trio/controllers"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/middleware"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
}

func SetContext(geminiClient *genai.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("GeminiClient", geminiClient)
		c.Next()
	}
}

func main() {
	r := gin.Default()
	clientAddress := os.Getenv("CLIENT_ADDRESS")

	geminiClient, err := clients.CreateGeminiClient(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	defer geminiClient.Close()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{clientAddress}
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}

	r.Use(cors.New(config))
	r.Use(SetContext(geminiClient))

	public := r.Group("/")
	{
		public.POST("/login", controllers.Login)
		public.POST("/signup", controllers.Signup)
	}

	authenticated := r.Group("/")
	authenticated.Use(middleware.CheckAuth())
	{

		authenticated.POST("/logout", controllers.Logout)
		authenticated.POST("/reset-password", controllers.ResetPassword)

		authenticated.GET("/completions", controllers.GetCompletion)

		// Chat related endpoints
		chats := authenticated.Group("/chats")
		{
			chats.POST("", controllers.CreateChat)
			chats.GET("", controllers.GetUserChats)
			chats.GET("/:chatId", controllers.GetChatInfo)
			chats.DELETE("/:chatId", controllers.DeleteChat)
			chats.PUT("/:chatId", controllers.UpdateChat)
			chats.POST("/:chatId/messages", controllers.AddMessageToChat)
			chats.POST("/:chatId/agents", controllers.AddAgentToChat)
			chats.POST("/create-with-agents", controllers.CreateChatWithAgents)
			chats.POST("/:chatId/messages/reflection", controllers.PostReflectionMessage)
		}

		user := authenticated.Group("/me")
		{
			user.GET("", controllers.GetCurrentUser)
		}

		// Agent related endpoints
		agents := authenticated.Group("/agents")
		{
			agents.GET("/:agentId", controllers.GetAgent)
			agents.PUT("/:agentId", controllers.UpdateAgent)
			agents.DELETE("/:agentId", controllers.DeleteAgent)
		}
	}

	r.Run()
}
