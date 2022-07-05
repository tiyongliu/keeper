package modules

type SimpleSettingMysql struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type MysqlDatabase struct {
	Name string `json:"name"`
}
