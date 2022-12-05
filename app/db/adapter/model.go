package adapter

/*
type DatabaseInfoObjects struct {
	//Tables []*TableInfo `json:"tables"`
	//Collections []*CollectionInfo `json:"collections"`
	Tables      []*modules.MysqlTableSchema  `json:"tables"`
	Collections []*modules.MongoDBCollection `json:"collections"`
	Views       []*ViewInfo                  `json:"views"`

	MatViews   []*ViewInfo      `json:"matviews"`
	Procedures []*ProcedureInfo `json:"procedures"`
	Functions  []*FunctionInfo  `json:"functions"`
	Triggers   []*TriggerInfo   `json:"triggers"`
}

type DatabaseInfo struct {
	DatabaseInfoObjects
	Schemas       []*SchemaInfo `json:"schemas"`
	Engine        string        `json:"engine"`
	DefaultSchema string        `json:"defaultSchema"`
}

type SchemaInfo struct {
	ObjectId   string
	SchemaName string
}

type DatabaseObjectInfo struct {
	PairingId       string
	ObjectId        string
	CreateDate      string
	ModifyDate      string
	HashCode        string
	ObjectTypeField string
	ObjectComment   string
}

type NamedObjectInfo struct {
	PureName   string
	SchemaName string
}

type ConstraintInfo struct {
	PairingId      string
	constraintName string
	ConstraintType ConstraintType
}

type ConstraintType string

const (
	ConstraintType_PrimaryKey ConstraintType = "primaryKey"
	ConstraintType_ForeignKey ConstraintType = "foreignKey"
	ConstraintType_Index      ConstraintType = "index"
	ConstraintType_Check      ConstraintType = "check"
	ConstraintType_Unique     ConstraintType = "unique"
)

type ColumnInfo struct {
	NamedObjectInfo
	PairingId          string
	ColumnName         string
	NotNull            bool
	AutoIncrement      bool
	DataType           string
	Precision          int
	Scale              int
	Length             int
	ComputedExpression string
	IsPersisted        bool
	IsSparse           bool
	DefaultValue       string
	DefaultConstraint  string
	ColumnComment      string
	IsUnsigned         bool
	IsZerofill         bool
}

type ColumnsConstraintInfo struct {
	Columns []*ColumnReference
	ConstraintInfo
}

type ColumnReference struct {
	ColumnName       string
	RefColumnName    string
	IsIncludedColumn bool
	IsDescending     bool
}

type PrimaryKeyInfo struct {
	ColumnsConstraintInfo
}

type ForeignKeyInfo struct {
	RefSchemaName string
	RefTableName  string
	UpdateAction  string
	DeleteAction  string
	ColumnsConstraintInfo
}

type IndexInfo struct {
	IsUnique  bool
	IndexType IndexType
}

type IndexType string

const (
	IndexType_Normal    IndexType = "normal"
	IndexType_Clustered IndexType = "clustered"
	IndexType_Xml       IndexType = "xml"
	IndexType_Spatial   IndexType = "spatial"
	IndexType_Fulltext  IndexType = "fulltext"
)

type UniqueInfo struct {
	ColumnsConstraintInfo
}

type CheckInfo struct {
	Definition string
	ConstraintInfo
}

type TableInfo struct {
	Columns      []*ColumnInfo
	PrimaryKey   PrimaryKeyInfo
	ForeignKeys  []*ForeignKeyInfo
	Dependencies []*ForeignKeyInfo
	Indexes      []*IndexInfo
	Uniques      []*UniqueInfo
	Checks       CheckInfo
	// preloadedRows?: any[];
	PreloadedRowsKey        []string
	preloadedRowsInsertOnly []string
	//tableRowCount?: number | string;

	TableRowCount      int
	IsDynamicStructure bool
}

type CollectionInfo struct {
	DatabaseObjectInfo
}

type ViewInfo struct {
	Columns []*ColumnInfo `json:"columns"`
}

type ProcedureInfo struct {
	SqlObjectInfo
}

type FunctionInfo struct {
	SqlObjectInfo
}

type TriggerInfo struct {
	SqlObjectInfo
}

type SqlObjectInfo struct {
	DatabaseObjectInfo
	CreateSql      string `json:"createSql"`
	RequiresFormat bool   `json:"requiresFormat"`
}
*/
