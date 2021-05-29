package main

import (
	"context"
	"github.com/Take-A-Seat/storage"
	"github.com/Take-A-Seat/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getAllTypesRestaurant() ([]models.TypeRestaurant, error) {
	var typesRestaurant []models.TypeRestaurant

	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return nil, err
	}

	typesRestaurantsCollection := client.Database(mongoDatabase).Collection("typesRestaurants")
	cursor, err := typesRestaurantsCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var typeRestaurant models.TypeRestaurant
		err := cursor.Decode(&typeRestaurant)
		if err != nil {
			return nil, err
		}

		typesRestaurant = append(typesRestaurant, typeRestaurant)
	}

	return typesRestaurant, nil
}

func getTypesFromRestaurantId(restaurantId string) ([]models.TypeRestaurantRelation, error) {
	var typesRestaurant []models.TypeRestaurantRelation

	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return typesRestaurant, err
	}

	restaurantIdObject, err := primitive.ObjectIDFromHex(restaurantId)
	if err != nil {
		return typesRestaurant, err
	}

	filter := bson.M{"restaurantId": restaurantIdObject}

	typesRestaurantsCollection := client.Database(mongoDatabase).Collection("typesRestaurantsRelations")
	cursor, err := typesRestaurantsCollection.Find(context.Background(), filter)
	if err != nil {
		return typesRestaurant, err
	}

	for cursor.Next(context.TODO()) {
		var typeRestaurant models.TypeRestaurantRelation
		err := cursor.Decode(&typeRestaurant)
		if err != nil {
			return typesRestaurant, err
		}

		typesRestaurant = append(typesRestaurant, typeRestaurant)
	}

	return typesRestaurant, nil
}

func updateTypesRestaurant(listReceived []primitive.ObjectID, restaurantId string) error {
	var listToAdd []primitive.ObjectID
	var listToRemove []primitive.ObjectID
	listRelationRestaurant, _ := getTypesFromRestaurantId(restaurantId)

	for _, item := range listReceived {
		if checkContainTypeRelation(listRelationRestaurant, item) == false {
			listToAdd = append(listToAdd, item)
		}
	}

	for _, item := range listRelationRestaurant {
		if ContainId(listReceived, item.TypeRestaurantId) == false {
			listToRemove = append(listToRemove, item.TypeRestaurantId)
		}
	}

	restaurantIdObject, err := primitive.ObjectIDFromHex(restaurantId)

	for _, item := range listToAdd {
		specificRelationId := primitive.NewObjectID()
		err = createTypeRelation(models.TypeRestaurantRelation{
			Id:               specificRelationId,
			RestaurantId:     restaurantIdObject,
			TypeRestaurantId: item,
		})
		if err != nil {
			return err
		}
	}

	for _, item := range listToRemove {
		err = deleteTypeRelation(restaurantIdObject, item)
		if err != nil {
			return err
		}
	}

	return nil
}

func checkContainTypeRelation(list []models.TypeRestaurantRelation, id primitive.ObjectID) bool {
	for _, item := range list {
		if item.TypeRestaurantId == id {
			return true
		}
	}

	return false
}

func createTypeRelation(typeRelation models.TypeRestaurantRelation) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	specificRestaurantRelationsCollection := client.Database(mongoDatabase).Collection("typesRestaurantsRelations")
	_, err = specificRestaurantRelationsCollection.InsertOne(context.Background(), bson.M{
		"_id":              typeRelation.Id,
		"restaurantId":     typeRelation.RestaurantId,
		"typeRestaurantId": typeRelation.TypeRestaurantId,
	})
	if err != nil {
		return err
	}

	return nil
}

func deleteTypeRelation(restaurantId primitive.ObjectID, typeId primitive.ObjectID) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	specificRestaurantRelationsCollection := client.Database(mongoDatabase).Collection("typesRestaurantsRelations")

	filter := bson.D{{"restaurantId", restaurantId}, {"typeRestaurantId", typeId}}
	_, err = specificRestaurantRelationsCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return nil
}
