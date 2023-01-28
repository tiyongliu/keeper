package schema

type CommandTypeEnum string

const (
	Command_Type_Insert                CommandTypeEnum = "insert"
	Command_Type_Select                CommandTypeEnum = "select"
	Command_Type_Update                CommandTypeEnum = "update"
	Command_Type_Delete                CommandTypeEnum = "delete"
	Command_Type_Allow_Identity_Insert CommandTypeEnum = "allowIdentityInsert"
)

type Select struct {
	Columns     []*Column
	CommandType CommandTypeEnum `json:"commandType"`
	From        *From           `json:"from"`
	OrderBy     *OrderBy        `json:"orderBy"`
	Range       *Range          `json:"range"`
}

type Column struct {
	Alias         string      `json:"alias"`
	AutoIncrement bool        `json:"autoIncrement"`
	ColumnComment string      `json:"columnComment"`
	ColumnName    string      `json:"columnName"`
	DataType      string      `json:"dataType"`
	DefaultValue  interface{} `json:"defaultValue"`
	ExprType      string      `json:"exprType"`
	IsUnsigned    bool        `json:"isUnsigned"`
	IsZerofill    bool        `json:"isZerofill"`
	NotNull       bool        `json:"notNull"`
	PureName      string      `json:"pureName"`
}

type ColumnSource struct {
	Alias string `json:"alias"`
}

type From struct {
	Alias string `json:"alias"`
	Name  *Name
}

type Name struct {
	PureName string `json:"pureName"`
}

type OrderBy struct {
	ColumnName string `json:"columnName"`
	Direction  string `json:"direction"`
	ExprType   string `json:"exprType"`
}

type Range struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func (s *Select) ColumnsFileAlias() (filesAlias []string) {
	for _, column := range s.Columns {
		filesAlias = append(filesAlias, column.Alias)
	}

	return filesAlias
}
