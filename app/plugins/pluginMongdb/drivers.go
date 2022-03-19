package pluginMongdb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"keeper/app/pkg/logger"
	"keeper/app/pkg/standard"
	"keeper/app/plugins/modules"
)

type MongoDBDrivers struct {
	DB *mongo.Client
}

func NewMongoDB(db *mongo.Client) standard.SqlStandard {
	return &MongoDBDrivers{db}
}

func (mg *MongoDBDrivers) Dialect() string {
	return standard.MONGOALIAS
}

func (mg *MongoDBDrivers) Connect() interface{} {
	return nil
}

func (mg *MongoDBDrivers) GetPoolInfo() interface{} {
	return mg.DB
}

func (mg *MongoDBDrivers) GetVersion() (*standard.VersionMsg, error) {
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

//func (mg *MongoDBDrivers) Tables(args ...string) (interface{}, error) {
//	databaseName := args[0]
//	names, err := mg.DB.Database(databaseName).ListCollectionNames(context.Background(), bson.D{})
//	if err != nil {
//		return nil, err
//	}
//
//	dialect := mg.Dialect()
//	var collections []*modules.MongoDBCollection
//	for _, name := range names {
//		collections = append(collections, &modules.MongoDBCollection{
//			PureName: name,
//			Engine:   dialect,
//		})
//	}
//
//	return collections, nil
//}

func (mg *MongoDBDrivers) Collections(databaseName string) ([]*modules.MongoDBCollection, error) {
	names, err := mg.DB.Database(databaseName).ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	//dialect := mg.Dialect()
	var collections []*modules.MongoDBCollection

	for _, name := range names {
		collections = append(collections, &modules.MongoDBCollection{
			PureName: name,
			//Engine:   dialect,
		})
	}

	return collections, nil
}

func (mg *MongoDBDrivers) Columns(databaseName, tableName string) (interface{}, error) {
	return nil, nil
}

func (mg *MongoDBDrivers) ListCollections(databaseName string) {
	ctx := context.Background()
	cursor, err := mg.DB.Database(databaseName).ListCollections(ctx, bson.D{})
	if err != nil {
		fmt.Println(err)
		return
	}

	for cursor.Next(ctx) {
		elements, err := cursor.Current.Elements()
		//{"name": "open_order_scan_history","type": "collection","options": {},"info": {"readOnly": false,"uuid": {"$binary":{"base64":"1gNSWotySlmYnJ9ZXgx3nQ==","subType":"04"}}},"idIndex": {"v": {"$numberInt":"2"},"key": {"_id": {"$numberInt":"1"}},"name": "_id_"}}
		logger.Infof("ele %s, err: %v", elements, err)
	}
}

func (mg *MongoDBDrivers) Ping() error {
	return mg.DB.Ping(context.Background(), nil)
}
