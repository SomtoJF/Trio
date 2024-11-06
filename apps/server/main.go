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
	"github.com/somtojf/trio/qdrantpackage"

	docs "github.com/somtojf/trio/docs"
	"github.com/somtojf/trio/initializers"
	"github.com/somtojf/trio/middleware"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	var collections = []qdrantpackage.CollectionName{qdrantpackage.Messages}

	initializers.LoadEnvVariables()
	initializers.ConnectToDb()
	qdrantpackage.ConnectToQdrant()

	qdrantpackage.CreateQdrantCollections(qdrantpackage.QdrantClient, collections)
}

func SetContext(geminiClient *genai.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("GeminiClient", geminiClient)
		c.Next()
	}
}

// @title			Trio API
// @Schemes
// @version		1.0
// @description	Trio API Server
// @contact.name	Somtochukwu Francis
// @contact.email	somtofrancis5@gmail.com
// @host		localhost:4000
// @BasePath	/
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

	docs.SwaggerInfo.BasePath = "/"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
			chats.POST("/:chatId/messages", controllers.NewMessage)
			chats.POST("/:chatId/agents", controllers.AddAgentToChat)
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
