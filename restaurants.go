package main

import (
	"context"
	"github.com/Take-A-Seat/storage"
	"github.com/Take-A-Seat/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
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

	err = addManager(restaurantId, userId)
	if err != nil {
		return err
	}

	return nil
}

func getAllRestaurants() ([]models.Restaurant, error) {
	var listRestaurants []models.Restaurant

	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return nil, err
	}

	restaurantsCollection := client.Database(mongoDatabase).Collection("restaurants")
	cursor, err := restaurantsCollection.Find(context.Background(), bson.M{"deleteAt": time.Time{}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var restaurant models.Restaurant
		err := cursor.Decode(&restaurant)
		if err != nil {
			return nil, err
		}

		listRestaurants = append(listRestaurants, restaurant)
	}

	return listRestaurants, nil
}

func getRestaurantById(restaurantId primitive.ObjectID) (models.Restaurant, error) {
	var restaurant models.Restaurant
	var filter = bson.M{"_id": restaurantId}

	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return models.Restaurant{}, err
	}

	restaurantsCollection := client.Database(mongoDatabase).Collection("restaurants")
	err = restaurantsCollection.FindOne(context.Background(), filter).Decode(&restaurant)
	if err != nil {
		return models.Restaurant{}, err
	}

	return restaurant, nil
}

func updateRestaurant(restaurant models.Restaurant) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return  err
	}

	restaurantsCollection := client.Database(mongoDatabase).Collection("restaurants")

	update := bson.D{
		{"$set", bson.D{
			{"name",restaurant.Name},
			{"description",restaurant.Description},
			{"address",restaurant.Description},
			{"phone",restaurant.Phone},
			{"program",restaurant.Program},
			{"postCode",restaurant.PostCode},
		}},
	}

	errorUpdate :=restaurantsCollection.FindOneAndUpdate(context.Background(),bson.M{"_id": restaurant.Id},update)
	if errorUpdate.Err()!=nil{
		return errorUpdate.Err()
	}

	return nil
}