package plugin_mondb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"system/pkg/standard"
)

type MongoDBDrivers struct {
}

func NewMongoDB() standard.SqlStandard {
	return &MongoDBDrivers{}
}

func (mongdb *MongoDBDrivers) GetVersion() (interface{}, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())
	db := client.Database("local")

	buildInfoCmd := bson.D{bson.E{Key: "buildInfo", Value: 1}}
	var buildInfoDoc bson.M
	if err := db.RunCommand(context.TODO(), buildInfoCmd).Decode(&buildInfoDoc); err != nil {
		return nil, err
	}
	return &standard.VersionMsg{
		Version:     fmt.Sprintf("%s", buildInfoDoc["version"]),
		VersionText: fmt.Sprintf("MongoDB %s", buildInfoDoc["version"]),
	}, nil
}

func (mongdb *MongoDBDrivers) ListDatabases() (interface{}, error) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(context.Background())

	_, err = client.ListDatabaseNames(context.Background(), bson.M{})
	//databases, err := client.ListDatabaseNames(context.Background(), bson.M{})

	return nil, err
}
