package schema

type DatabaseAnalyser struct {
}

func NewDatabaseAnalyser() {

}

func CreateEmptyStructure() *DatabaseInfo {
	return &DatabaseInfo{}
}

func (d *DatabaseAnalyser) incrementalAnalysis(structure *DatabaseInfo) {
	if d._getFastSnapshot() != nil {
		return
	}
}

func (d *DatabaseAnalyser) _getFastSnapshot() error {
	return nil
}
