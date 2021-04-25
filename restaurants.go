package main

import (
	"context"
	"github.com/Take-A-Seat/storage"
	"github.com/Take-A-Seat/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createRestaurant(restaurant models.Restaurant, userId primitive.ObjectID) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	restaurantsCollection := client.Database(mongoDatabase).Collection("restaurants")

	restaurantId := primitive.NewObjectID()
	_, err = restaurantsCollection.InsertOne(context.Background(), bson.M{
		"_id":         restaurantId,
		"name":        restaurant.Name,
		"description": restaurant.Description,
		"address":     restaurant.Address,
		"phone":       restaurant.Phone,
		"program":     restaurant.Program,
		"postCode":    restaurant.PostCode,
		"deleteAt":    restaurant.DeleteAt})

	if err != nil {
		return err
	}

	err = addManager(restaurantId,userId)
	if err != nil {
		return err
	}

	return nil
}
