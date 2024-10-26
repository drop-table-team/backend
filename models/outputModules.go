package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type OutputModule struct {
	ModuleName    string `json:"moduleName"`
	ModuleAddress string `json:"moduleAddress"`
}

type UnregisterName struct {
	Name string `json:"name"`
}

func AddModule(client *mongo.Client, module OutputModule) {
	modules := client.Database("mongo_data").Collection("outputModules")

	_, err := modules.InsertOne(context.TODO(), module)
	if err != nil {
		log.Fatal(err)
	}

}

func RemoveModule(client *mongo.Client, name UnregisterName) {
	modules := client.Database("mongo_data").Collection("outputModules")
	filter := bson.D{{"modulename", name.Name}}

	_, err := modules.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
}
