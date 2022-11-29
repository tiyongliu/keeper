package mongo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"keeper/app/db/standard/modules"
)

func (s *Source) Dialect() string {
	return Adapter
}

func (s *Source) Ping() error {
	return s.client.Ping(s.ctx, nil)
}

func (s *Source) Version() (*modules.Version, error) {
	//todo 这里是写死，有问题，需要调整。
	db := s.client.Database("local")
	buildInfoCmd := bson.D{bson.E{Key: "buildInfo", Value: 1}}
	var buildInfoDoc bson.M
	if err := db.RunCommand(s.ctx, buildInfoCmd).Decode(&buildInfoDoc); err != nil {
		return nil, err
	}

	return &modules.Version{
		Version:     fmt.Sprintf("%s", buildInfoDoc["version"]),
		VersionText: fmt.Sprintf("MongoDB %s", buildInfoDoc["version"]),
	}, nil
}

func (s *Source) Close() error {
	if s.client != nil {
		return s.client.Connect(s.ctx)
	}
	return nil
}

func (s *Source) ListDatabases() (interface{}, error) {

	buildInfoCmd := bson.D{bson.E{Key: "listDatabases", Value: 1}}
	var buildInfoDoc modules.MongoDBDatabaseList
	db := s.client.Database("admin")
	err := db.RunCommand(s.ctx, buildInfoCmd).Decode(&buildInfoDoc)
	return buildInfoDoc.Databases, err
}

func (s *Source) Query(sql string) (interface{}, error) {
	return nil, nil
}
