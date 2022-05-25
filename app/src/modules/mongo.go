package modules

type SimpleSettingMongoDB struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

// mongoUrl := "mongodb://" + user + ":" + password + "@" + url + "/" + dbname
//clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
