package plugin_mondb

import (
	"bytes"
	"context"
	"dbbox/app/src/modules"
	"dbbox/app/src/pkg/standard"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBDrivers struct {
	DB *mongo.Client
}

func NewMongoDB() standard.SqlStandard {
	return &MongoDBDrivers{}
}

func (mongdb *MongoDBDrivers) GetPoolInfo() interface{} {
	return mongdb.DB
}

func NewSimpleMongoDBPool(setting *modules.SimpleSettingMongoDB) (standard.SqlStandard, error) {
	if setting == nil {
		return nil, fmt.Errorf("setting is nil")
	}
	if setting.Username == "" {
		//return nil, fmt.Errorf("lack of setting.Username")
	}
	if setting.Password == "" {
		//return nil, fmt.Errorf("lack of setting.Password")
	}
	if setting.Host == "" {
		return nil, fmt.Errorf("lack of setting.Host")
	}
	if setting.Port == "" {
		return nil, fmt.Errorf("lack of setting.Port")
	}

	// mongoUrl := "mongodb://" + user + ":" + password + "@" + url + "/" + dbname
	//clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var buf bytes.Buffer
	buf.WriteString("mongodb://")
	if setting.Username != "" {
		buf.WriteString(setting.Username)
		buf.WriteString(":")
	}

	if setting.Password != "" {
		buf.WriteString(setting.Password)
		buf.WriteString("@")
	}

	buf.WriteString(setting.Host)
	buf.WriteString(":")
	buf.WriteString(setting.Port)

	db, err := mongo.Connect(context.Background(), options.Client().ApplyURI(buf.String()))
	if err != nil {
		return nil, err
	}

	//err = db.Ping(context.Background(), nil)
	//if err != nil {
	//	return nil, err
	//}

	return &MongoDBDrivers{db}, nil
}

func (mongdb *MongoDBDrivers) GetVersion() (interface{}, error) {
	defer mongdb.DB.Disconnect(context.Background())
	db := mongdb.DB.Database("local")

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
