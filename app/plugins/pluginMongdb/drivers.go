package pluginMongdb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"keeper/app/code"
	"keeper/app/modules"
	"keeper/app/pkg/standard"
)

type MongoDBDrivers struct {
	DB *mongo.Client
}

func NewMongoDB(db *mongo.Client) standard.SqlStandard {
	return &MongoDBDrivers{db}
}

func (mg *MongoDBDrivers) Dialect() string {
	return code.MONGOALIAS
}

func (mg *MongoDBDrivers) Connect() interface{} {
	return nil
}

func (mg *MongoDBDrivers) GetPoolInfo() interface{} {
	return mg.DB
}

func (mg *MongoDBDrivers) GetVersion() (interface{}, error) {
	db := mg.DB.Database("local")
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

func (mg *MongoDBDrivers) ListDatabases() (interface{}, error) {
	buildInfoCmd := bson.D{bson.E{Key: "listDatabases", Value: 1}}
	var buildInfoDoc modules.MongoDBDatabaseList
	db := mg.DB.Database("admin")
	err := db.RunCommand(context.TODO(), buildInfoCmd).Decode(&buildInfoDoc)
	return buildInfoDoc.Databases, err
}

func (mg *MongoDBDrivers) Close() error {
	return mg.DB.Disconnect(context.Background())
}

func (mg *MongoDBDrivers) Tables(databaseName, tableName string) (interface{}, error) {
	names, err := mg.DB.Database("auth").ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	dialect := mg.Dialect()
	var collections []*modules.MongoDBCollection
	for _, name := range names {
		collections = append(collections, &modules.MongoDBCollection{
			PureName: name,
			Engine:   dialect,
		})
	}

	return collections, nil
}

func (mg *MongoDBDrivers) Columns(databaseName, tableName string) (interface{}, error) {
	return nil, nil
}
