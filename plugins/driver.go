package plugins

var RequiredPlugin *DriversPlugin

type DriversPlugin struct {
	PackageName string
	Drivers     []*Driver
}

type Driver struct {
	Dialect                 *Dialect
	DatabaseEngineTypes     []string
	SupportedCreateDatabase bool
	Engine                  string
	Title                   string
	EditorMode              string
	DefaultPort             int
	SupportsDatabaseUrl     bool
	databaseUrlPlaceholder  string
	importExportArgs        []*ImportExportArg
}

type Dialect struct {
	LimitSelect            bool
	RangeSelect            bool
	OffsetFetchRangeSyntax bool
	StringEscapeChar       string
	FallbackDataType       string
}

type ImportExportArg struct {
	Type      string
	Name      string
	Label     string
	ApiName   string
	Direction string
}

func init() {
	RequiredPlugin = &DriversPlugin{
		PackageName: "dbgate-plugin-mongo",
		Drivers: []*Driver{{
			Dialect: &Dialect{
				LimitSelect:            true,
				RangeSelect:            true,
				OffsetFetchRangeSyntax: true,
				StringEscapeChar:       "'",
				FallbackDataType:       "nvarchar(max)",
			},
			DatabaseEngineTypes:     []string{"document"},
			SupportedCreateDatabase: true,
			Engine:                  "mongo@dbgate-plugin-mongo",
			Title:                   "MongoDB",
			EditorMode:              "javascript",
			DefaultPort:             27017,
			SupportsDatabaseUrl:     true,
			databaseUrlPlaceholder:  "e.g. mongodb://username:password@mongodb.mydomain.net/dbname",
			importExportArgs: []*ImportExportArg{{
				Type:      "checkbox",
				Name:      "createStringId",
				Label:     "Create string _id attribute",
				ApiName:   "createStringId",
				Direction: "target",
			}},
		}},
	}
}
