package mongo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"keeper/app/db/standard/modules"
	"keeper/app/pkg/logger"
)

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

func (s *Source) Collections(databaseName string) ([]*modules.MongoDBCollection, error) {
	names, err := s.client.Database(databaseName).ListCollectionNames(s.ctx, bson.D{})
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

func (s *Source) Columns(databaseName, tableName string) (interface{}, error) {
	return nil, nil
}

func (s *Source) ListCollections(databaseName string) {
	cursor, err := s.client.Database(databaseName).ListCollections(s.ctx, bson.D{})
	if err != nil {
		fmt.Println(err)
		return
	}

	for cursor.Next(s.ctx) {
		elements, err := cursor.Current.Elements()
		//{"name": "open_order_scan_history","type": "collection","options": {},"info": {"readOnly": false,"uuid": {"$binary":{"base64":"1gNSWotySlmYnJ9ZXgx3nQ==","subType":"04"}}},"idIndex": {"v": {"$numberInt":"2"},"key": {"_id": {"$numberInt":"1"}},"name": "_id_"}}
		logger.Infof("ele %s, err: %v", elements, err)
	}
}
