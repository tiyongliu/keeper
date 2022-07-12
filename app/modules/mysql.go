package modules

type SimpleSettingMysql struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type MysqlDatabase struct {
	Name string `json:"name"`
}

type MysqlTableSchema struct {
	PureName      string `json:"pureName"`
	ObjectType    string `json:"objectType"`
	TableRowCount int    `json:"tableRowCount"`
	ModifyDate    string `json:"modifyDate"`
}

type MysqlTableColumn struct {
	PureName         string      `json:"pureName"`
	ColumnName       string      `json:"columnName"`
	IsNullable       string      `json:"isNullable"`
	DAtaType         string      `json:"dataType"`
	CharMaxLength    *int        `json:"charMaxLength,omitempty"`
	NumericPrecision *int        `json:"numericPrecision,omitempty"`
	NumericScale     *int        `json:"numericScale,omitempty"`
	DefaultValue     interface{} `json:"defaultValue"`
	ColumnComment    interface{} `json:"columnComment"`
	ColumnType       string      `json:"columnType"`
	EXTRA            string      `json:"extra"`
}

type MysqlPrimaryKey struct {
	ConstraintName string `json:"constraintName"`
	PureName       string `json:"pureName"`
	ColumnName     string `json:"columnName"`
}

type MysqlForeignKeys struct {
	ConstraintName string `json:"constraintName"`
	PureName       string `json:"pureName"`
	UpdateAction   string `json:"updateAction"`
	DeleteAction   string `json:"deleteAction"`
	RefTableName   string `json:"refTableName"`
	ColumnName     string `json:"columnName"`
	RefColumnName  string `json:"refColumnName"`
}
