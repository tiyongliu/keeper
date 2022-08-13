package schema

type DatabaseAnalyser struct {
}

func NewDatabaseAnalyser() {

}

func CreateEmptyStructure() *DatabaseInfoObjects {
	return &DatabaseInfoObjects{
		Tables:      []*TableInfo{},
		Collections: []*CollectionInfo{},
		Views:       []*ViewInfo{},
		MatViews:    []*ViewInfo{},
		Procedures:  []*ProcedureInfo{},
		Functions:   []*FunctionInfo{},
		Triggers:    []*TriggerInfo{},
	}
}

//func (d *DatabaseAnalyser) incrementalAnalysis(structure *DatabaseInfo) {
//	if d._getFastSnapshot() != nil {
//		return
//	}
//}
//
//func (d *DatabaseAnalyser) _getFastSnapshot() error {
//	return nil
//}
