package main

import (
	"context"
	"encoding/json"
	"github.com/Take-A-Seat/storage"
	"github.com/Take-A-Seat/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

//type RestaurantWithDetails struct {
//	RestaurantDetails models.Restaurant                   `json:"restaurantDetails"`
//	ListSpecifics     []models.SpecificRestaurantRelation `json:"listSpecifics"`
//	ListTypes         []models.TypeRestaurantRelation     `json:"listTypes"`
//}

func createRestaurant(restaurant models.Restaurant, userId primitive.ObjectID, form *multipart.Form) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	restaurantsCollection := client.Database(mongoDatabase).Collection("restaurants")

	restaurant.Email = strings.Join(form.Value["email"], "")
	restaurant.Country = strings.Join(form.Value["country"], "")
	restaurant.Facebook = strings.Join(form.Value["facebook"], "")
	restaurant.Instagram = strings.Join(form.Value["instagram"], "")
	restaurant.Twitter = strings.Join(form.Value["twitter"], "")
	restaurant.Website = strings.Join(form.Value["website"], "")
	restaurant.Name = strings.Join(form.Value["name"], "")
	restaurant.Phone = strings.Join(form.Value["phone"], "")
	restaurant.PostCode, _ = strconv.Atoi(strings.Join(form.Value["postCode"], ""))
	restaurant.Lat, _ = strconv.ParseFloat(strings.Join(form.Value["lat"], ""), 64)
	restaurant.Lng, _ = strconv.ParseFloat(strings.Join(form.Value["lng"], ""), 64)
	restaurant.Description = strings.Join(form.Value["description"], "")
	restaurant.City = strings.Join(form.Value["city"], "")
	restaurant.StreetAndNumber = strings.Join(form.Value["streetAndNumber"], "")
	restaurant.VisibleOnline, _ = strconv.ParseBool(strings.Join(form.Value["visibleOnline"], ""))
	restaurant.Province = strings.Join(form.Value["province"], "")

	programString := strings.Join(form.Value["program"], "")
	err = json.Unmarshal([]byte(programString), &restaurant.Program)
	if err != nil {
		return err
	}

	restaurantId := primitive.NewObjectID()

	filePrefix := restaurantId.Hex()
	if len(form.File["logo"]) > 0 {
		file := form.File["logo"][0]
		newFile, err := storage.HandleFile(file, filePrefix, hostname+"/files/"+restaurantId.Hex()+"/", directoryFiles+restaurantId.Hex()+"/")
		if err != nil {
			return err
		}

		restaurant.Logo = models.File(newFile)
	}

	_, err = restaurantsCollection.InsertOne(context.Background(), bson.M{
		"_id":             restaurantId,
		"name":            restaurant.Name,
		"description":     restaurant.Description,
		"phone":           restaurant.Phone,
		"program":         restaurant.Program,
		"postCode":        restaurant.PostCode,
		"country":         restaurant.Country,
		"email":           restaurant.Email,
		"website":         restaurant.Website,
		"facebook":        restaurant.Facebook,
		"instagram":       restaurant.Instagram,
		"twitter":         restaurant.Twitter,
		"logo":            restaurant.Logo,
		"streetAndNumber": restaurant.StreetAndNumber,
		"city":            restaurant.City,
		"lat":             restaurant.Lat,
		"lng":             restaurant.Lng,
		"province":        restaurant.Province,
		"visibleOnline":   restaurant.VisibleOnline,
		"deleteAt":        restaurant.DeleteAt})

	if err != nil {
		return err
	}

	err = addManager(restaurantId, userId)
	if err != nil {
		return err
	}

	return nil
}

func getAllRestaurants() ([]models.RestaurantWithDetails, error) {
	var listRestaurants []models.RestaurantWithDetails

	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return nil, err
	}

	restaurantsCollection := client.Database(mongoDatabase).Collection("restaurants")
	cursor, err := restaurantsCollection.Find(context.Background(), bson.M{"deleteAt": time.Time{}, "visibleOnline": true})
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var restaurant models.RestaurantWithDetails
		err := cursor.Decode(&restaurant.RestaurantDetails)
		if err != nil {
			return nil, err
		}

		restaurant.ListTypes, err = getTypesFromRestaurantId(restaurant.RestaurantDetails.Id.Hex())
		if err != nil {
			return nil, err
		}

		restaurant.ListSpecifics, err = getSpecificFromRestaurantId(restaurant.RestaurantDetails.Id.Hex())
		if err != nil {
			return nil, err
		}

		listRestaurants = append(listRestaurants, restaurant)
	}

	return listRestaurants, nil
}

func getRestaurantById(restaurantId primitive.ObjectID) (models.RestaurantWithDetails, error) {
	var restaurant models.RestaurantWithDetails
	var filter = bson.M{"_id": restaurantId}

	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return models.RestaurantWithDetails{}, err
	}

	restaurantsCollection := client.Database(mongoDatabase).Collection("restaurants")
	err = restaurantsCollection.FindOne(context.Background(), filter).Decode(&restaurant.RestaurantDetails)

	if err != nil {
		return models.RestaurantWithDetails{}, err
	}
	restaurant.ListTypes, err = getTypesFromRestaurantId(restaurant.RestaurantDetails.Id.Hex())
	if err != nil {
		return models.RestaurantWithDetails{}, err
	}

	restaurant.ListSpecifics, err = getSpecificFromRestaurantId(restaurant.RestaurantDetails.Id.Hex())
	if err != nil {
		return models.RestaurantWithDetails{}, err
	}

	return restaurant, nil
}

func getRestaurantByManagerId(managerId primitive.ObjectID) (models.RestaurantWithDetails, int, error) {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return models.RestaurantWithDetails{}, 400, err
	}

	var managerRelation ManagerRelation
	restaurantsCollection := client.Database(mongoDatabase).Collection("managers")
	count, err := restaurantsCollection.CountDocuments(context.Background(), bson.M{"userId": managerId})
	if count == 0 || err != nil {
		return models.RestaurantWithDetails{}, 404, err
	} else {
		err = restaurantsCollection.FindOne(context.Background(), bson.M{"userId": managerId}).Decode(&managerRelation)
		if err != nil {
			return models.RestaurantWithDetails{}, 400, err
		} else {
			restaurant, err := getRestaurantById(managerRelation.RestaurantId)
			if err != nil {
				return models.RestaurantWithDetails{}, 400, err
			}

			return restaurant, 200, nil
		}
	}

}

func updateRestaurant(restaurant models.Restaurant, form *multipart.Form) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	restaurant.Email = strings.Join(form.Value["email"], "")
	restaurant.Country = strings.Join(form.Value["country"], "")
	restaurant.Facebook = strings.Join(form.Value["facebook"], "")
	restaurant.Instagram = strings.Join(form.Value["instagram"], "")
	restaurant.Twitter = strings.Join(form.Value["twitter"], "")
	restaurant.Website = strings.Join(form.Value["website"], "")
	restaurant.Name = strings.Join(form.Value["name"], "")
	restaurant.Phone = strings.Join(form.Value["phone"], "")
	restaurant.PostCode, _ = strconv.Atoi(strings.Join(form.Value["postCode"], ""))
	restaurant.Lat, _ = strconv.ParseFloat(strings.Join(form.Value["lat"], ""), 64)
	restaurant.Lng, _ = strconv.ParseFloat(strings.Join(form.Value["lng"], ""), 64)
	restaurant.Description = strings.Join(form.Value["description"], "")
	restaurant.City = strings.Join(form.Value["city"], "")
	restaurant.StreetAndNumber = strings.Join(form.Value["streetAndNumber"], "")
	restaurant.Province = strings.Join(form.Value["province"], "")
	restaurant.VisibleOnline, _ = strconv.ParseBool(strings.Join(form.Value["visibleOnline"], ""))
	filePrefix := strings.Join(form.Value["id"], "")
	changeLogo, _ := strconv.ParseBool(strings.Join(form.Value["changeLogo"], ""))
	programString := strings.Join(form.Value["program"], "")

	err = json.Unmarshal([]byte(programString), &restaurant.Program)
	if err != nil {
		return err
	}

	restaurantIdObject, err := primitive.ObjectIDFromHex(filePrefix)
	if err != nil {
		return err
	}

	restaurantFromDB, err := getRestaurantById(restaurantIdObject)
	if err != nil {
		return err
	}

	if changeLogo == true {
		if len(form.File["logo"]) > 0 {
			file := form.File["logo"][0]

			newFile, err := storage.HandleFile(file, filePrefix, hostname+"/files/"+filePrefix+"/", directoryFiles+filePrefix+"/")
			if err != nil {
				return err
			}
			restaurant.Logo = models.File(newFile)
		}

	} else {
		restaurant.Logo = restaurantFromDB.RestaurantDetails.Logo
	}

	restaurantsCollection := client.Database(mongoDatabase).Collection("restaurants")

	update := bson.D{
		{"$set", bson.D{
			{"name", restaurant.Name},
			{"description", restaurant.Description},
			{"phone", restaurant.Phone},
			{"program", restaurant.Program},
			{"postCode", restaurant.PostCode},
			{"logo", restaurant.Logo},
			{"email", restaurant.Email},
			{"country", restaurant.Country},
			{"city", restaurant.City},
			{"streetAndNumber", restaurant.StreetAndNumber},
			{"province", restaurant.Province},
			{"facebook", restaurant.Facebook},
			{"instagram", restaurant.Instagram},
			{"twitter", restaurant.Twitter},
			{"website", restaurant.Website},
			{"lat", restaurant.Lat},
			{"lng", restaurant.Lng},
			{"visibleOnline", restaurant.VisibleOnline},
		}},
	}

	errorUpdate := restaurantsCollection.FindOneAndUpdate(context.Background(), bson.M{"_id": restaurantIdObject}, update)
	if errorUpdate.Err() != nil {
		return errorUpdate.Err()
	}

	return nil
}