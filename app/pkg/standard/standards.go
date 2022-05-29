package standard

type SqlStandard interface {
	GetPoolInfo() interface{}
	GetVersion() (interface{}, error)
	ListDatabases() (interface{}, error)
	Close() error
}
