package standard

type SqlStandard interface {
	GetVersion() (interface{}, error)
	ListDatabases() (interface{}, error)
}
