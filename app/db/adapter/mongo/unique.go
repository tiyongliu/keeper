package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (s *Source) ReadCollection(database string, opt *modules.CollectionDataOptions) (interface{}, error) {
	ctx := context.Background()
	collection := s.client.Database(database).Collection(opt.PureName)
	if opt.CountDocuments {
		count, err := countDocuments(ctx, collection, opt)
		if err != nil {
			logger.Errorf("exec countDocuments [database: %s, collection: %s] failed %v", database, opt.PureName, err)
			return 0, err
		}
		return count, nil
	} else if opt.Aggregate != nil {
		rows, err := aggregate(ctx, collection, opt)
		if err != nil {
			logger.Errorf("exec aggregate [database: %s, collection: %s] failed %v", database, opt.PureName, err)
			return nil, err
		}
		return rows, nil
	} else {
		rows, err := find(ctx, collection, opt)
		if err != nil {
			logger.Errorf("exec find [database: %s, collection: %s] failed %v", database, opt.PureName, err)
			return nil, err
		}
		return rows, nil
	}
}

func countDocuments(ctx context.Context, collection *mongo.Collection, opt *modules.CollectionDataOptions) (int64, error) {
	return collection.CountDocuments(ctx, opt.Condition)
}

func aggregate(ctx context.Context, collection *mongo.Collection, opt *modules.CollectionDataOptions) ([]bson.M, error) {
	results := make([]bson.M, 0)
	cursor, err := collection.Aggregate(ctx, []bson.D{})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &results)
	return results, err
}

func find(ctx context.Context, collection *mongo.Collection, opt *modules.CollectionDataOptions) ([]bson.M, error) {
	results := make([]bson.M, 0)
	cursor, err := collection.Find(ctx, opt.Condition, &options.FindOptions{
		Limit: &opt.Limit,
		Skip:  &opt.Skip,
		Sort:  opt.Sort,
	})
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &results)
	return results, err
}
