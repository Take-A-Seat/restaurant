package main

import (
	"context"
	"github.com/Take-A-Seat/storage"
	"github.com/Take-A-Seat/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strconv"
	"time"
)

func createArea(area models.Area) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	areasCollection := client.Database(mongoDatabase).Collection("areas")
	area.Id = primitive.NewObjectID()
	_, err = areasCollection.InsertOne(context.Background(), bson.M{
		"_id":            area.Id,
		"name":           area.Name,
		"displayName":    area.DisplayName,
		"priority":       area.Priority,
		"restaurantId":   area.RestaurantId,
		"onlineCapacity": area.OnlineCapacity,
		"deleteAt":       time.Time{},
	})
	if err != nil {
		return err
	}

	return nil
}

func updateArea(updateArea models.Area, areaId string) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	areasCollection := client.Database(mongoDatabase).Collection("areas")
	areaIdObject, err := primitive.ObjectIDFromHex(areaId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": areaIdObject}
	updateObject := bson.D{{"$set", bson.D{
		{"name", updateArea.Name},
		{"displayName", updateArea.DisplayName},
		{"priority", updateArea.Priority},
		{"onlineCapacity", updateArea.OnlineCapacity},
	}}}

	_, err = areasCollection.UpdateOne(context.Background(), filter, updateObject)
	if err != nil {
		return err
	}

	return nil
}

func deleteArea(areaId string) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	areasCollection := client.Database(mongoDatabase).Collection("areas")

	tablesCollection := client.Database(mongoDatabase).Collection("tables")
	areaIdObject, err := primitive.ObjectIDFromHex(areaId)
	if err != nil {
		return err
	}

	filterArea := bson.M{"_id": areaIdObject}
	_, err = areasCollection.UpdateOne(context.Background(), filterArea, bson.D{{"$set", bson.D{{"deleteAt", time.Now()}}}})
	if err != nil {
		return err
	}

	filterTables := bson.M{"areaId": areaIdObject}
	_, err = tablesCollection.UpdateMany(context.Background(), filterTables, bson.D{{"$set", bson.D{{"deleteAt", time.Now()}}}})
	if err != nil {
		return err
	}

	return nil
}

func getAreaById(areaId string) (models.Area, error) {
	var area models.Area

	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return area, err
	}

	areasCollection := client.Database(mongoDatabase).Collection("areas")
	areaIdObject, err := primitive.ObjectIDFromHex(areaId)
	if err != nil {
		return area, err
	}

	filterArea := bson.M{"_id": areaIdObject, "deleteAt": time.Time{}}
	err = areasCollection.FindOne(context.Background(), filterArea).Decode(&area)
	if err != nil {
		return area, err
	}

	return area, nil
}

func getAreasByRestaurantId(restaurantId string) ([]models.Area, error) {
	var listAreas []models.Area
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return listAreas, err
	}

	areasCollection := client.Database(mongoDatabase).Collection("areas")
	restaurantObjId, err := primitive.ObjectIDFromHex(restaurantId)
	if err != nil {
		return listAreas, err
	}

	filterArea := bson.M{"restaurantId": restaurantObjId, "deleteAt": time.Time{}}
	cursor, err := areasCollection.Find(context.Background(), filterArea)
	if err != nil {
		return listAreas, err
	}

	for cursor.Next(context.TODO()) {
		var area models.Area
		err = cursor.Decode(&area)
		if err != nil {
			return listAreas, err
		}

		listTables, _ := getTablesByAreaId(area.Id.Hex())
		totalNumberPlace := 0
		availableFree := 0
		for _, table := range listTables {
			totalNumberPlace += table.MaxPeople
			if table.AvailableNow == true {
				availableFree += table.MaxPeople
			}
		}
		area.Capacity = strconv.Itoa(availableFree) + "/" + strconv.Itoa(totalNumberPlace)
		area.NumberTables = len(listTables)
		listAreas = append(listAreas, area)
	}

	return listAreas, nil
}
