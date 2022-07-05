package bridge

type DatabaseConnections struct {
}

func NewDatabaseConnections() *DatabaseConnections {
	return &DatabaseConnections{}
}

func (dc *DatabaseConnections) Refresh(conid string) {

}
