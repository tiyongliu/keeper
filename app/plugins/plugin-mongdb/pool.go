package plugin_mongdb

import (
	"bytes"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"keeper/app/modules"
	"keeper/app/pkg/standard"
)

func NewSimpleMongoDBPool(setting *modules.SimpleSettingMongoDB) (standard.SqlStandard, error) {
	if setting == nil {
		return nil, fmt.Errorf("setting is nil")
	}
	if setting.Port == "" {
		return nil, fmt.Errorf("lack of setting.Port")
	}
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

	err = db.Ping(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return NewMongoDB(db), nil
}
