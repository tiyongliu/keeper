package standard

type SqlStandard interface {
	Dialect() string
	Connect() interface{}
	GetPoolInfo() interface{}
	GetVersion() (interface{}, error)
	ListDatabases() (interface{}, error)
	Close() error
	Tables(databaseName, tableName string) (interface{}, error)
	Columns(databaseName, tableName string) (interface{}, error)
}
