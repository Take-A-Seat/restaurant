package main

import (
	"errors"
	"github.com/Take-A-Seat/storage"
	"github.com/Take-A-Seat/storage/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func handleGetTablesByAreaId(c *gin.Context) {
	var areaId = c.Param("areaId")
	tables, err := getTablesByAreaId(areaId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, tables)
	}
}

func handleGetTableById(c *gin.Context) {
	var tableId = c.Param("tableId")
	table, err := getTableById(tableId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, table)
	}
}

func handleDeleteTable(c *gin.Context) {

	var tableId = c.Param("tableId")
	err := deleteTable(tableId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"error": "Delete  table success"})
	}
}

func handleUpdateTable(c *gin.Context) {
	var table models.Table

	err := c.ShouldBindJSON(&table)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var tableId = c.Param("tableId")
	err = updateTable(table, tableId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"error": "Update  table success"})
	}
}

func handleCreateTable(c *gin.Context) {
	var table models.Table

	err := c.ShouldBindJSON(&table)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = createTable(table)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"error": "Success create table"})
	}
}

func handleGetAreasByRestaurantId(c *gin.Context) {
	var restaurantId = c.Param("id")
	areas, err := getAreasByRestaurantId(restaurantId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, areas)
	}
}

func handleGetAreaById(c *gin.Context) {
	var areaId = c.Param("areaId")
	area, err := getAreaById(areaId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, area)
	}
}

func handleDeleteArea(c *gin.Context) {
	var area models.Area

	err := c.ShouldBindJSON(&area)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var areaId = c.Param("areaId")
	err = deleteArea(areaId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"error": "Delete  area success"})
	}
}

func handleUpdateArea(c *gin.Context) {
	var area models.Area

	err := c.ShouldBindJSON(&area)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var areaId = c.Param("areaId")
	err = updateArea(area, areaId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"error": "Update  area success"})
	}
}

func handleCreateArea(c *gin.Context) {
	var area models.Area

	err := c.ShouldBindJSON(&area)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = createArea(area)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"error": "Success create area"})
	}
}

func handleCreateRestaurant(c *gin.Context) {
	var restaurant models.Restaurant

	loggedInUserId, err := storage.GetLoggedInUserId(c, apiUrl)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form, _ := c.MultipartForm()
	if err := createRestaurant(restaurant, loggedInUserId, form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	} else {
		c.JSON(http.StatusCreated, gin.H{"error": "Success create restaurant"})
	}
}

func handleUpdateRestaurant(c *gin.Context) {
	var restaurant models.Restaurant

	form, _ := c.MultipartForm()
	err := updateRestaurant(restaurant, form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func getRestaurantByManagerIdHandler(c *gin.Context) {
	managerId := c.Param("id")
	managerIdObject, err := primitive.ObjectIDFromHex(managerId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": errors.New("Error parse id into primitive objct")})
		return
	}

	restaurant, code, err := getRestaurantByManagerId(managerIdObject)
	if code == 200 {
		c.JSON(http.StatusOK, restaurant)
		return
	}

	if code == 404 {
		c.JSON(http.StatusNotFound, gin.H{"message": "This account has no restaurant"})
		return
	}

	if code == 400 || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err})
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