package main

import (
	"context"
	"fmt"
	"github.com/Take-A-Seat/storage"
	"github.com/Take-A-Seat/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func createMenu(menu models.Menu) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	menusCollection := client.Database(mongoDatabase).Collection("menus")
	menu.Id = primitive.NewObjectID()
	for indexPage, page := range menu.Pages {
		for indexSection, section := range page.Sections {
			for indexProduct, _ := range section.Products {
				menu.Pages[indexPage].Sections[indexSection].Products[indexProduct].Id = primitive.NewObjectID()
			}
		}
	}
	_, err = menusCollection.InsertOne(context.Background(), bson.M{
		"_id":          menu.Id,
		"restaurantId": menu.RestaurantId,
		"pages":        menu.Pages,
	})
	if err != nil {
		return err
	}
	return nil
}

func createOrUpdateMenu(menu models.Menu, restaurantId string) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	menusCollection := client.Database(mongoDatabase).Collection("menus")
	restaurantIdObj, err := primitive.ObjectIDFromHex(restaurantId)
	if err != nil {
		return err
	}

	filter := bson.M{"restaurantId": restaurantIdObj}
	numberMenus, err := menusCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		return err
	}

	if numberMenus == 0 {
		err = createMenu(menu)
		if err != nil {
			return err
		}
	} else {
		updateObject := bson.D{{"$set", bson.D{
			{"pages", menu.Pages},
		}}}

		_, err = menusCollection.UpdateOne(context.Background(), filter, updateObject)
		if err != nil {
			return err
		}
	}

	return nil
}

func getMenuByRestaurantId(restaurantId string) (models.Menu, error) {
	var menu models.Menu

	restaurantObjId, err := primitive.ObjectIDFromHex(restaurantId)
	if err != nil {
		return models.Menu{}, err
	}

	var filter = bson.M{"restaurantId": restaurantObjId}
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return models.Menu{}, err
	}

	fmt.Println("da", restaurantObjId)
	menusCollection := client.Database(mongoDatabase).Collection("menus")
	err = menusCollection.FindOne(context.Background(), filter).Decode(&menu)
	if err != nil {
		return models.Menu{}, err
	}

	return menu, nil
}
