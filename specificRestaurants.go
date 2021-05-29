package main

import (
	"context"
	"fmt"
	"github.com/Take-A-Seat/storage"
	"github.com/Take-A-Seat/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getAllSpecific() ([]models.SpecificRestaurant, error) {
	var listSpecific []models.SpecificRestaurant

	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return nil, err
	}

	specificRestaurantCollection := client.Database(mongoDatabase).Collection("specificRestaurant")
	cursor, err := specificRestaurantCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var specific models.SpecificRestaurant
		err := cursor.Decode(&specific)
		if err != nil {
			return nil, err
		}

		listSpecific = append(listSpecific, specific)
	}

	return listSpecific, nil
}

func getSpecificFromRestaurantId(restaurantId string) ([]models.SpecificRestaurantRelation, error) {
	var listSpecific []models.SpecificRestaurantRelation

	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return listSpecific, err
	}

	specificRestaurantCollection := client.Database(mongoDatabase).Collection("specificRestaurantRelations")
	restaurantIdObject, err := primitive.ObjectIDFromHex(restaurantId)
	if err != nil {
		return listSpecific, err
	}

	filter := bson.M{"restaurantId": restaurantIdObject}
	cursor, err := specificRestaurantCollection.Find(context.Background(), filter)
	if err != nil {
		return listSpecific, err
	}

	for cursor.Next(context.TODO()) {
		var specific models.SpecificRestaurantRelation
		err := cursor.Decode(&specific)
		if err != nil {
			return listSpecific, err
		}

		listSpecific = append(listSpecific, specific)
	}

	return listSpecific, nil
}

func updateSpecificsRestaurant(listReceived []primitive.ObjectID, restaurantId string) error {
	var listToAdd []primitive.ObjectID
	var listToRemove []primitive.ObjectID
	listRelationRestaurant, _ := getSpecificFromRestaurantId(restaurantId)

	for _, item := range listReceived {
		if checkContainSpecificRelation(listRelationRestaurant, item) == false {
			listToAdd = append(listToAdd, item)
		}
	}

	for _, item := range listRelationRestaurant {
		if ContainId(listReceived, item.SpecificRestaurantId) == false {
			listToRemove = append(listToRemove, item.SpecificRestaurantId)
		}
	}

	fmt.Println("listRemove:", listToRemove)
	fmt.Println("listToAdd:", listToAdd)

	restaurantIdObject, err := primitive.ObjectIDFromHex(restaurantId)

	for _, item := range listToAdd {
		specificRelationId := primitive.NewObjectID()
		err = createSpecificRelation(models.SpecificRestaurantRelation{
			Id:                   specificRelationId,
			RestaurantId:         restaurantIdObject,
			SpecificRestaurantId: item,
		})
		if err != nil {
			return err
		}
	}

	for _, item := range listToRemove {
		err = deleteSpecificRelation(restaurantIdObject, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkContainSpecificRelation(list []models.SpecificRestaurantRelation, id primitive.ObjectID) bool {
	for _, item := range list {
		if item.SpecificRestaurantId == id {
			return true
		}
	}

	return false
}

func ContainId(list []primitive.ObjectID, id primitive.ObjectID) bool {
	for _, item := range list {
		if item == id {
			return true
		}
	}
	return false
}

func createSpecificRelation(specificRelation models.SpecificRestaurantRelation) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	specificRestaurantRelationsCollection := client.Database(mongoDatabase).Collection("specificRestaurantRelations")
	_, err = specificRestaurantRelationsCollection.InsertOne(context.Background(), bson.M{
		"_id":                  specificRelation.Id,
		"restaurantId":         specificRelation.RestaurantId,
		"specificRestaurantId": specificRelation.SpecificRestaurantId,
	})
	if err != nil {
		return err
	}

	return nil
}

func deleteSpecificRelation(restaurantId primitive.ObjectID, specificId primitive.ObjectID) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	specificRestaurantRelationsCollection := client.Database(mongoDatabase).Collection("specificRestaurantRelations")

	filter := bson.D{{"restaurantId", restaurantId}, {"specificRestaurantId", specificId}}
	_, err = specificRestaurantRelationsCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}
