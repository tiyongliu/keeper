package standard

const (
	MYSQLALIAS = "mysql"
	MONGOALIAS = "mongo"
)

type SqlStandard interface {
	Dialect() string
	Ping() error
	Connect() interface{}
	GetPoolInfo() interface{}
	GetVersion() (*VersionMsg, error)
	ListDatabases() (interface{}, error)
	Close() error
	//Tables(...string) (interface{}, error)
	//Columns(databaseName, tableName string) (interface{}, error)
}
