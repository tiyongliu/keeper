package modules

type MysqlConnection struct {
	Server     string `json:"server"`
	Engine     string `json:"engine"`
	SshMode    string `json:"sshMode"`
	SshPort    string `json:"sshPort"`
	SshKeyfile string `json:"sshKeyfile"`
	User       string `json:"username"`
	Password   string `json:"password"`
}

type SimpleSettingMysql struct {
	Host            string `yaml:"host" json:"host"`
	Username        string `yaml:"username" json:"username"`
	Password        string `yaml:"password" json:"password"`
	DBName          string `yaml:"dBName" json:"dbName"`
	Charset         string `yaml:"charset" json:"charset"`
	MaxIdle         int    `yaml:"maxIdle" json:"maxIdle"`
	MaxOpen         int    `yaml:"maxOpen" json:"maxOpen"`
	Loc             string `yaml:"loc" json:"loc"`
	MultiStatements bool   `yaml:"multiStatements" json:"multiStatements"`
	ParseTime       bool   `yaml:"parseTime" json:"parseTime"`
	ShowSql         bool   `yaml:"showSql" json:"showSql"`
}

/*
[0] {
[0]   server: 'localhost',
[0]   engine: 'mysql@dbgate-plugin-mysql',
[0]   sshMode: 'userPassword',
[0]   sshPort: '22',
[0]   sshKeyfile: '/Users/liuliutiyong/.ssh/id_rsa',
[0]   user: 'root',
[0]   password: '123456'
[0] }
*/
