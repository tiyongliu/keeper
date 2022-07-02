package standard

type SqlStandard interface {
	Dialect() string
	Connect() interface{}
	GetPoolInfo() interface{}
	GetVersion() (interface{}, error)
	ListDatabases() (interface{}, error)
	Close() error
}
