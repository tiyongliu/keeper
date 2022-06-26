package standard

type SqlStandard interface {
	Connect() interface{}
	GetPoolInfo() interface{}
	GetVersion() (interface{}, error)
	ListDatabases() (interface{}, error)
	Close() error
}
