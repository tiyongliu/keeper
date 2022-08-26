package plugins

type DatabaseInfoObjects struct {
	Tables      []map[string]interface{} `json:"tables"`
	Collections []map[string]interface{} `json:"collections"`
	Views       []map[string]interface{} `json:"views"`
	MatViews    []map[string]interface{} `json:"matviews"`
	Procedures  []map[string]interface{} `json:"procedures"`
	Functions   []map[string]interface{} `json:"functions"`
	Triggers    []map[string]interface{} `json:"triggers"`
}
