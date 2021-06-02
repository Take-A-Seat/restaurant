package main

import (
	"github.com/Take-A-Seat/auth/validatorAuth"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"time"
)

var mongoHost = "takeaseat.knilq.mongodb.net"
var mongoUser = "admin"
var mongoPass = "p4r0l4"
var mongoDatabase = "TakeASeat"
var apiUrl = "https://api.takeaseat.site"
var hostname = "https://api.takeaseat.site"
var directoryFiles = "/home/takeaseat/manager/web/files/"

func main() {
	port := os.Getenv("TAKEASEAT_RESTAURANTS_PORT")
	if port == "" {
		port = "9215"
	}

	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "DELETE", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accepts", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           1 * time.Minute,
		AllowCredentials: true,
	}))

	//privateRoutesUsers need Authorization token in header
	protectedUsers := router.Group("/restaurants")
	protectedUsers.Use(validatorAuth.AuthMiddleware(apiUrl + "/auth/isAuthenticated"))
	{

		//restaurant
		protectedUsers.POST("/", handleCreateRestaurant)
		protectedUsers.GET("/managerId/:id", getRestaurantByManagerIdHandler)
		protectedUsers.PUT("/id/:id", handleUpdateRestaurant)

		//area
		protectedUsers.POST("/id/:id/area", handleCreateArea)
		protectedUsers.PUT("/id/:id/area/:areaId", handleUpdateArea)
		protectedUsers.DELETE("/id/:id/area/:areaId", handleDeleteArea)
		protectedUsers.GET("/id/:id/area/:areaId", handleGetAreaById)
		protectedUsers.GET("/id/:id/areas", handleGetAreasByRestaurantId)

		//table
		protectedUsers.POST("/id/:id/area/:areaId/table", handleCreateTable)
		protectedUsers.PUT("/id/:id/area/:areaId/table/:tableId", handleUpdateTable)
		protectedUsers.DELETE("/id/:id/area/:areaId/table/:tableId", handleDeleteTable)
		protectedUsers.GET("/id/:id/area/:areaId/table/:tableId", handleGetTableById)
		protectedUsers.GET("/id/:id/areas/:areaId/tables", handleGetTablesByAreaId)

		//menu
		protectedUsers.POST("/id/:id/menu", handleCreateOrUpdateMenu)

		//specificsRestaurant
		protectedUsers.GET("/id/:id/specificsRestaurant",handleGetSpecificsFromRestaurant)
		protectedUsers.POST("/id/:id/specificsRestaurant",handleUpdateSpecificsRestaurant)

		//typesRestaurant
		protectedUsers.GET("/id/:id/typesRestaurant",handleGetTypesFromRestaurant)
		protectedUsers.POST("/id/:id/typesRestaurant",handleUpdateTypesRestaurant)

	}

	freeRoute := router.Group("/restaurants")
	{
		freeRoute.GET("/", handleGetAllRestaurants)

		freeRoute.GET("/id/:id/menu", handleGetMenuByRestaurantId)

		freeRoute.GET("/specificsRestaurant",handleGetAllSpecificsRestaurant)

		freeRoute.GET("/typesRestaurant",handleGetAllTypesRestaurant)

		freeRoute.GET("/id/:id", handleGetRestaurantById)



	}

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Port already in use!")
	}

}
