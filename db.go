package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Email struct {
	Id        string
	Recipient string
	Sender    string
	Subject   string
	Body      string
}

var MongoClient *mongo.Client
var MongoCollection *mongo.Collection

func CreateMongoClient() {
	var err error

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(Config.MongoURI).SetServerAPIOptions(serverAPI)

	MongoClient, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	MongoCollection = MongoClient.Database("lasagnamail").Collection("emails")
}

func GetEmailById(id string) Email {
	var result Email

	MongoCollection.FindOne(context.TODO(), bson.D{{Key: "id", Value: id}}).Decode(&result)

	return result
}

func GetInbox(address string) []Email {
	var results []Email

	cursor, err := MongoCollection.Find(context.TODO(), bson.D{{Key: "recipient", Value: address}})
	if err != nil {
		return []Email{}
	}

	defer func() {
		cursor.Close(context.Background())
	}()

	for cursor.Next(context.Background()) {
		var result Email

		cursor.Decode(&result)

		results = append(results, result)
	}

	return results
}

func CreateEmail(email Email) {
	_, err := MongoCollection.InsertOne(context.TODO(), email)
	if err != nil {
		panic(err)
	}
}
