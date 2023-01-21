package modules

type MongoDBDatabase struct {
	Name       string `bson:"name" json:"name"`
	SizeOnDisk int    `bson:"sizeOnDisk" json:"sizeOnDisk"`
	Empty      bool   `bson:"empty" json:"empty"`
}

type MongoDBDatabaseList struct {
	Databases []*MongoDBDatabase `bson:"databases" json:"databases"`
	TotalSize int                `bson:"totalSize" json:"totalSize"`
	Ok        int                `bson:"ok" json:"ok"`
}

type MongoDBCollection struct {
	PureName string `json:"pureName"`
	//Engine   string `json:"engine"`
}

type CollectionDataOptions struct {
	PureName       string                 `json:"pureName"`
	CountDocuments bool                   `json:"countDocuments"`
	Limit          int64                  `json:"limit"`
	Skip           int64                  `json:"skip"`
	Aggregate      interface{}            `json:"aggregate"`
	Condition      map[string]interface{} `json:"condition"`
	Sort           map[string]int         `json:"sort"`
}
