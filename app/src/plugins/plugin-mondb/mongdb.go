package plugin_mondb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func Connect() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)

	databases, err := client.ListDatabaseNames(context.TODO(), bson.M{})

	fmt.Println(databases)
	fmt.Println(err)

	db := client.Database("your-db-name")
	buildInfoCmd := bson.D{bson.E{Key: "buildInfo", Value: 1}}
	var buildInfoDoc bson.M
	if err := db.RunCommand(context.TODO(), buildInfoCmd).Decode(&buildInfoDoc); err != nil {
		log.Printf("Failed to run buildInfo command: %v", err)
		return
	}
	log.Println("Database version:", buildInfoDoc)
	//client.Ping(context.Background(), nil)
	listDatabases()
}

type List struct {
	Databases []*struct {
		Name       string `bson:"name" json:"name"`
		SizeOnDisk int    `bson:"sizeOnDisk" json:"sizeOnDisk"`
		Empty      bool   `bson:"empty" json:"empty"`
	} `bson:"databases" json:"databases"`
	TotalSize int `bson:"totalSize" json:"totalSize"`
	Ok        int `bson:"ok" json:"ok"`
}

func listDatabases() {
	//db.adminCommand('listDatabases');
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())
	db := client.Database("admin")
	buildInfoCmd := bson.D{bson.E{Key: "listDatabases", Value: 1}}
	var buildInfoDoc List
	command := db.RunCommand(context.TODO(), buildInfoCmd).Decode(&buildInfoDoc)
	fmt.Println(command)
}
