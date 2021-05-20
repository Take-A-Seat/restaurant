package main

import (
	"context"
	"github.com/Take-A-Seat/storage"
	"github.com/Take-A-Seat/storage/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func createTable(table models.Table) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	tablesCollection := client.Database(mongoDatabase).Collection("tables")
	table.Id = primitive.NewObjectID()
	_, err = tablesCollection.InsertOne(context.Background(), bson.M{
		"_id":             table.Id,
		"number":          table.Number,
		"areaId":          table.AreaId,
		"tableGroupId":    table.TableGroupId,
		"priority":        table.Priority,
		"availableOnline": table.AvailableOnline,
		"minPeople":       table.MinPeople,
		"maxPeople":       table.MaxPeople,
		"deleteAt":        time.Time{},
	})
	if err != nil {
		return err
	}

	return nil
}

func updateTable(updateTable models.Table, tableId string) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	tablesCollection := client.Database(mongoDatabase).Collection("tables")
	tableIdObject, err := primitive.ObjectIDFromHex(tableId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": tableIdObject}
	updateObject := bson.D{{"$set", bson.D{
		{"tableGroupId", updateTable.TableGroupId},
		{"areaId", updateTable.AreaId},
		{"number", updateTable.Number},
		{"priority", updateTable.Priority},
		{"availableOnline", updateTable.AvailableOnline},
		{"minPeople", updateTable.MinPeople},
		{"maxPeople", updateTable.MaxPeople},
	}}}

	_, err = tablesCollection.UpdateOne(context.Background(), filter, updateObject)
	if err != nil {
		return err
	}

	return nil
}

func deleteTable(tableId string) error {
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return err
	}

	tablesCollection := client.Database(mongoDatabase).Collection("tables")
	tableIdObject, err := primitive.ObjectIDFromHex(tableId)
	if err != nil {
		return err
	}

	filterTables := bson.M{"_id": tableIdObject}
	_, err = tablesCollection.UpdateOne(context.Background(), filterTables, bson.D{{"deleteAt", time.Now()}})
	if err != nil {
		return err
	}

	return nil
}

func getTableById(tableId string) (models.Table, error) {
	var table models.Table

	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return table, err
	}

	tablesCollection := client.Database(mongoDatabase).Collection("tables")
	tableIdObject, err := primitive.ObjectIDFromHex(tableId)
	if err != nil {
		return table, err
	}

	filterArea := bson.M{"_id": tableIdObject, "deleteAt": time.Time{}}
	err = tablesCollection.FindOne(context.Background(), filterArea).Decode(&table)
	if err != nil {
		return table, err
	}

	return table, nil
}

func getTablesByAreaId(areaId string) ([]models.Table, error) {
	var listTable []models.Table
	client, err := storage.ConnectToDatabase(mongoUser, mongoPass, mongoHost, mongoDatabase)
	defer storage.DisconnectFromDatabase(client)
	if err != nil {
		return nil, err
	}

	tablesCollection := client.Database(mongoDatabase).Collection("tables")
	areaObjectId, err := primitive.ObjectIDFromHex(areaId)
	if err != nil {
		return nil, err
	}

	filterTable := bson.M{"areaId": areaObjectId, "deleteAt": time.Time{}}
	cursor, err := tablesCollection.Find(context.Background(), filterTable)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.TODO()) {
		var table models.Table
		err = cursor.Decode(&table)
		if err != nil {
			return nil, err
		}

		listTable = append(listTable, table)

	}

	return listTable, nil
}
