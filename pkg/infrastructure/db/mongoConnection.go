package db_chatNcallSvc

import (
	"context"
	"fmt"
	"log"
	"time"

	config_chatNcallSvc "github.com/ashkarax/ciao_socialMedia_chatNcallService/pkg/infrastructure/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDbCollections struct {
	OneToOneChats  *mongo.Collection
	OneToManyChats *mongo.Collection
}

func ConnectDatabaseMongo(config *config_chatNcallSvc.MongoDataBase) (*MongoDbCollections, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := fmt.Sprintf("mongodb://%s:%s@%s", config.MongoUsername, config.MongoPassword, config.MongoDbURL)
	fmt.Println("----------connection uri:", uri)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI))
	if err != nil {
		return nil, err
	}

	// defer func() {
	// 	if err = client.Disconnect(ctx); err != nil {
	// 		//defer function to close the client connection to mongodb at localhost:27017,if any error occur while closing,it will panic
	// 		panic(err)
	// 	}
	// }()

	ctx, cancel = context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		fmt.Println("can't ping to db:,err:", err)
		return nil, err
	}

	fmt.Printf("\nconnected to mongodb,on databse %s\n", config.DataBase)

	var mongoCollections MongoDbCollections
	mongoCollections.OneToOneChats = client.Database(config.DataBase).Collection("OneToOneChats")
	mongoCollections.OneToManyChats = client.Database(config.DataBase).Collection("OneToManyChats")

	// Insert the string into the collection
	_, err = mongoCollections.OneToManyChats.InsertOne(ctx, bson.D{
		{Key: "message", Value: "Your string data here"}})
	if err != nil {
		log.Fatal(err)
	}

	return &mongoCollections, nil
}
