package router

import (
	"goServer/src/api"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func RouterInit() {
	var r = gin.Default()
	// r.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{
	// 		"message": "pong",
	// 	})
	// })
	r.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"https://foo.com"},
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type"},
		// ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		// MaxAge: 12 * time.Hour,
	}))

	r.POST("/api/stock/", api.AddStock)
	r.GET("/api/stocks", api.GetStocks)
	// r.DELETE(("/api/user/:id"), api.DeleteUser)
	// r.PUT("/api/user/:id", api.PutUser)
	// r.GET("/api/user/:id", api.GetUser)

	// r.POST("/api/users", api.CreateAccount)
	// r.DELETE(("/api/users"), api.DeleteUser)
	// r.PUT("/api/users", api.PutUser)
	// r.GET("api/users", api.GetUser)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
