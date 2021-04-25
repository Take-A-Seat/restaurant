package main

import (
	"context"
	"github.com/Take-A-Seat/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func addManager(restaurantId primitive.ObjectID, userId primitive.ObjectID) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	managerCollection := client.Database(mongoDatabase).Collection("managers")
	_, err = managerCollection.InsertOne(context.Background(), bson.M{"restaurantId": restaurantId, "userId": userId})
	if err!=nil{
		return err
	}

	return nil
}