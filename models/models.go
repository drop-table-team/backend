package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Entry struct {
	Title                   string
	Tags                    []string
	Short                   string
	UniversalRepresentative string
}

func AddEntry(client *mongo.Client, entry Entry) {
	entries := client.Database("mongo_data").Collection("entries")

	_, err := entries.InsertOne(context.TODO(), entry)
	if err != nil {
		log.Fatal(err)
	}

}

func ViewEntries(client *mongo.Client) ([]*Entry, error) {
	entries := client.Database("mongo_data").Collection("entries")

	findOptions := options.Find()

	var results []*Entry

	cur, err := entries.Find(context.TODO(), bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var elem Entry
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}

		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}
	return results, nil

}
