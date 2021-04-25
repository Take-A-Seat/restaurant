package main

import (
	"github.com/Take-A-Seat/storage"
	"github.com/Take-A-Seat/storage/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func handleCreateRestaurant(c *gin.Context)  {
	var restaurant models.Restaurant
	if err := c.ShouldBindJSON(&restaurant); err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}

	loggedInUserId, err := storage.GetLoggedInUserId(c, apiUrl)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}


	if err := createRestaurant(restaurant,loggedInUserId); err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{"error":err.Error()})
		return
	}else{
		c.JSON(http.StatusCreated, gin.H{"error":"Success create restaurant"})
	}
}

func handleUpdateRestaurant(c *gin.Context)  {

}

func handleGetRestaurantById(c *gin.Context)  {

}

func handleGetRestaurantsByUserId(c *gin.Context){

}

func handleDeleteRestaurant(c *gin.Context)  {

}