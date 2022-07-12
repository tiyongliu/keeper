package modules

type SimpleSettingMongoDB struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

// mongoUrl := "mongodb://" + user + ":" + password + "@" + url + "/" + dbname
//clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

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
	Engine   string `json:"engine"`
}
