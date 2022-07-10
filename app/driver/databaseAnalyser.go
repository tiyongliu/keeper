package driver

type DatabaseAnalyser struct {
}

func NewDatabaseAnalyser() {

}

func (d *DatabaseAnalyser) createEmptyStructure() {

}

func (d *DatabaseAnalyser) incrementalAnalysis(structure *Structure) {
	if d._getFastSnapshot() != nil {
		return
	}
}

func (d *DatabaseAnalyser) _getFastSnapshot() error {
	return nil
}
