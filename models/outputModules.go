package models

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type OutputModule struct {
	ModuleName    string `json:"moduleName"`
	ModuleAddress string `json:"moduleAddress"`
}

func AddModule(client *mongo.Client, module OutputModule) {
	modules := client.Database("mongo_data").Collection("outputModules")

	_, err := modules.InsertOne(context.TODO(), module)
	if err != nil {
		log.Fatal(err)
	}

}

func RemoveModule(client *mongo.Client, name string) {
	modules := client.Database("mongo_data").Collection("outputModules")
	filter := OutputModule{ModuleName: name}

	_, err := modules.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}
}
