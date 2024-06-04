package main

import (
	"log"
	"time"

	_ "ginguide/pkg/examples/app-engine/gophers"

	"github.com/gin-gonic/gin"
)

func main() {
	// hellowWorldServer()

	// customRouting()

	// routesGrouping()

	// separatingBusinessLogicFromControllers()
}

func LoggerMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		log.Printf("Req meth: %s | Status %d | Duration %d", c.Request.Method, c.Writer.Status(), duration)
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKey := c.GetHeader("X-API-Key")

		if apiKey == "" {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
		}

	}
}

func hellowWorldServer() {
	router := gin.Default()
	router.Use(LoggerMiddleWare())

	router.GET("/", func(c *gin.Context) { c.String(200, "Hello World!") })
	router.GET("bye", func(c *gin.Context) { c.String(200, "Goodbye!") })

	authGroup := router.Group("/api")
	authGroup.Use(AuthMiddleware())
	authGroup.GET("/data", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Authentificated and authorized!"})
	})

	router.Run(":8080")
}

func customRouting() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	router.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(200, "User ID"+id)
	})

	router.GET("/search", func(c *gin.Context) {
		query := c.DefaultQuery("q", "default-value")
		c.String(200, "Search query: "+query)
	})

	router.Run(":8080")
}

func routesGrouping() {
	router := gin.Default()

	public := router.Group("/public")
	{
		public.GET("/info", func(c *gin.Context) {
			c.String(200, "Public information")
		})

		public.GET("/products", func(c *gin.Context) {
			c.String(200, "Public product list")
		})
	}

	private := router.Group("/private")
	private.Use(AuthMiddleware())
	{
		private.GET("/data", func(c *gin.Context) {
			c.String(200, "Private data accessible after authentification")
		})

		private.POST("/create", func(c *gin.Context) {
			c.String(200, "Create a new resource")
		})
	}

	router.Run(":8080")
}

func separatingBusinessLogicFromControllers() {
	router := gin.Default()
	userController := &UserController{}
	router.GET("/users/:id", userController.GetUserInfo)

	router.Run()
}

type UserController struct{}

func (uc *UserController) GetUserInfo(c *gin.Context) {
	userID := c.Param("id")

	c.JSON(200, gin.H{"id": userID, "name": "John Doe", "email": "john@example.com"})
}
