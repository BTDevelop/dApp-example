/*
* Main API file for the server
* @author: Ayan Banerjee
* @Organization: Math & Cody
 */
package main

import (
	"ethential/dapp/controllers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	config := cors.DefaultConfig()

	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Authorization", "Content-Type"}
	r.Use(cors.New(config))
	r.Use(static.Serve("/", static.LocalFile("./public", true)))
	api := r.Group("/api/v0")
	{
		api.GET("/genToken/:clientID", controllers.GenTokenController)
		api.POST("/transferToken", controllers.TransferTokenController)
		api.POST("/getTokenBalance", controllers.TokenBalanceController)
		api.POST("/swap", controllers.SwapTokenController)
		api.POST("/approve", controllers.ApproveController)
	}
	err := r.Run(":4551")
	if err != nil {
		panic(err)
	}
}
