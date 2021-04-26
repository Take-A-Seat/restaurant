package main

import (
	"errors"
	"github.com/Take-A-Seat/storage"
	"github.com/Take-A-Seat/storage/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func handleUpdateRestaurant(c *gin.Context) {
	var restaurant models.Restaurant

	errBindJson := c.ShouldBindJSON(&restaurant)
	if errBindJson != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errBindJson.Error()})
		return
	}

	err := updateRestaurant(restaurant)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Update restaurant successfully"})

	}
}

func handleGetRestaurantById(c *gin.Context) {
	restaurantId := c.Param("id")
	restaurantObjId, err := primitive.ObjectIDFromHex(restaurantId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("Error parse id into primitive objct")})
		return
	}

	restaurant, err := getRestaurantById(restaurantObjId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("Error get restaurant by id")})
		return
	} else {
		c.JSON(http.StatusOK, restaurant)
	}
}

func handleGetAllRestaurants(c *gin.Context) {
	listRestaurants, err := getAllRestaurants()

	if err == nil {
		c.JSON(http.StatusOK, listRestaurants)
		return
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
}

func handleGetRestaurantsByUserId(c *gin.Context) {

}

func handleDeleteRestaurant(c *gin.Context) {

}