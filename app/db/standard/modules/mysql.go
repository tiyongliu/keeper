package modules

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

type PrimaryKey struct {
	ConstraintName string `json:"constraintName"`
	PureName       string `json:"pureName"`
	ColumnName     string `json:"columnName"`
}

type ForeignKeys struct {
	ConstraintName string `json:"constraintName"`
	PureName       string `json:"pureName"`
	UpdateAction   string `json:"updateAction"`
	DeleteAction   string `json:"deleteAction"`
	RefTableName   string `json:"refTableName"`
	ColumnName     string `json:"columnName"`
	RefColumnName  string `json:"refColumnName"`
}

type MysqlRowsResult struct {
	Rows    interface{} `json:"rows"`
	Columns []*Column   `json:"columns"`
}

type UniqueName struct {
	ConstraintName string `json:"constraintName"`
}

type Column struct {
	ColumnName string `json:"columnName"`
}

type Indexe struct {
	ConstraintName string `json:"constraintName"`
	TableName      string `json:"tableName"`
	ColumnName     string `json:"columnName"`
	IndexType      string `json:"indexType"`
	NonUnique      bool   `json:"nonUnique"`
}

type Table struct {
	PureName      string `json:"pureName"`
	TableRowCount int    `json:"tableRowCount"`
	ModifyDate    string `json:"modifyDate"`
}

type TableColumn struct {
	PureName         string      `json:"pureName"`
	ColumnName       string      `json:"columnName"`
	IsNullable       string      `json:"isNullable"`
	DataType         string      `json:"dataType"`
	CharMaxLength    *int        `json:"charMaxLength"`
	NumericPrecision *int        `json:"numericPrecision"`
	NumericScale     *int        `json:"numericScale"`
	DefaultValue     interface{} `json:"defaultValue"`
	ColumnComment    string      `json:"columnComment"`
	ColumnType       string      `json:"columnType"`
	Extra            interface{} `json:"extra"`
}

type TransformColumnInfo struct {
	NotNull       bool        `json:"notNull"`
	AutoIncrement bool        `json:"autoIncrement"`
	ColumnName    string      `json:"columnName"`
	ColumnComment string      `json:"columnComment"`
	DataType      string      `json:"dataType"`
	DefaultValue  interface{} `json:"defaultValue"`
	IsUnsigned    bool        `json:"isUnsigned"`
	IsZerofill    bool        `json:"isZerofill"`
}

type View struct {
	PureName   string `json:"pureName"`
	ModifyDate string `json:"modifyDate"`
}

type Programmable struct {
	PureName          string      `json:"pureName"`
	ObjectType        string      `json:"objectType"`
	ModifyDate        string      `json:"modifyDate"`
	ReturnDataType    string      `json:"returnDataType"`
	RoutineDefinition interface{} `json:"routineDefinition"`
	IsDeterministic   bool        `json:"isDeterministic"`
}

type ViewTexts struct {
	PureName string `json:"pureName"`
}

// todo 需要先写页面
type MysqlSelectRequest struct {
	CommandType string `json:"commandType"`
	Form        *Form  `json:"form"`
}

type Form struct {
	Name  ViewTexts `json:"name"`
	Alias string    `json:"alias"`
}

type MysqlSelectColumn struct{}
