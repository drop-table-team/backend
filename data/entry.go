package data

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Entry struct {
	Title         string   `json:"title"`
	Tags          []string `json:"tags"`
	Short         string   `json:"short"`
	Transcription string   `json:"transcription"`
	Uuid          string   `json:"uuid"`
}

func AddEntry(database *mongo.Database, entry Entry) {
	entries := database.Collection("entries")

	_, err := entries.InsertOne(context.Background(), entry)
	if err != nil {
		log.Fatal(err)
	}

}

func ViewEntries(database *mongo.Database) ([]*Entry, error) {
	entries := database.Collection("entries")

	findOptions := options.Find()

	var results []*Entry

	cur, err := entries.Find(context.Background(), bson.D{}, findOptions)
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
