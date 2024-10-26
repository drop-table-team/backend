package main

import (
	"backend/models"
	"context"
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"os"
)

type Env struct {
	client *mongo.Client
}

func getRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got / request\n")
	io.WriteString(w, "This is my website!\n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /hello request\n")
	io.WriteString(w, "Hello, HTTP!\n")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	credential := options.Credential{
		AuthMechanism: "SCRAM-SHA-256",
		Username:      os.Getenv("MONGO_INITDB_ROOT_USERNAME"),
		Password:      os.Getenv("MONGO_INITDB_ROOT_PASSWORD"),
	}

	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017").SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	env := &Env{client: client}

	testEntry := models.Entry{Title: "TestTitle", Tags: []string{"Tag1", "Tag2"}, Short: "ShortDesc", UniversalRepresentative: "hhh"}

	models.AddEntry(env.client, testEntry)

	res, err := models.ViewEntries(env.client)

	if err != nil {
		log.Fatal(err)
	}

	for _, result := range res {
		res, _ := bson.MarshalExtJSON(result, false, false)
		fmt.Println(string(res))
	}

	http.HandleFunc("/", getRoot)
	http.HandleFunc("/hello", getHello)
	err = http.ListenAndServe(":8081", nil)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
