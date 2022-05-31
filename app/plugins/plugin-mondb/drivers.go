package plugin_mondb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"keeper/app/pkg/standard"
)

type MongoDBDrivers struct {
	DB *mongo.Client
}

func NewMongoDB(db *mongo.Client) standard.SqlStandard {
	return &MongoDBDrivers{db}
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

	return map[string]string{
		"version":     fmt.Sprintf("%s", buildInfoDoc["version"]),
		"versionText": fmt.Sprintf("MongoDB %s", buildInfoDoc["version"]),
	}, nil
	//return &standard.VersionMsg{
	//	Version:     fmt.Sprintf("%s", buildInfoDoc["version"]),
	//	VersionText: fmt.Sprintf("MongoDB %s", buildInfoDoc["version"]),
	//}, nil
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

func (mg *MongoDBDrivers) ListDatabases() (interface{}, error) {
	buildInfoCmd := bson.D{bson.E{Key: "listDatabases", Value: 1}}
	var buildInfoDoc List
	db := mg.DB.Database("admin")
	err := db.RunCommand(context.TODO(), buildInfoCmd).Decode(&buildInfoDoc)
	return buildInfoDoc, err
}

func (mg *MongoDBDrivers) Close() error {
	return mg.DB.Disconnect(context.Background())
}
