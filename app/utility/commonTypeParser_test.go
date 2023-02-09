package utility

import (
	"fmt"
	"github.com/samber/lo"
	"keeper/app/db/standard/modules"
	"keeper/app/pkg/logger"
	"regexp"
	"testing"
)

func TestIsTypeString(t *testing.T) {
	reg := regexp.MustCompile("(?i)char|binary")
	fmt.Println(reg.MatchString(" 报告字符 bInary"))

	array := []*modules.ForeignKeys{
		{
			ConstraintName: "qrtz_simple_triggers_ibfk_1",
			PureName:       "qrtz_simple_triggers",
			UpdateAction:   "NO ACTION",
			DeleteAction:   "NO ACTION",
			RefTableName:   "qrtz_triggers",
			ColumnName:     "TRIGGER_GROUP",
			RefColumnName:  "TRIGGER_GROUP",
		},
	}

	all := lo.GroupBy[*modules.ForeignKeys, string](array, func(t *modules.ForeignKeys) string {
		return t.ConstraintName
	})

	logger.Infof("%s", ToJsonStr(all))
}
