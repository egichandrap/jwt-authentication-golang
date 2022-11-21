package main

import (
	"jwt-authentication-golang/controllers"
	"jwt-authentication-golang/database"
	"jwt-authentication-golang/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database
	database.Connect("root:12345678@tcp(localhost:3306)/JWT?parseTime=true")
	database.Migrate()
	// Initialize Router
	router := initRouter()
	router.Run(":8090")
}
func initRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/token", controllers.GenerateToken)
		api.POST("/user/register", controllers.RegisterUser)
		api.DELETE("/user/delete", controllers.DeleteUser)
		secured := api.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/ping", controllers.Ping)
		}
	}
	return router
}
