package modules

type MysqlConnection struct {
	Server     string `json:"server"`
	Engine     string `json:"engine"`
	SshMode    string `json:"sshMode"`
	SshPort    string `json:"sshPort"`
	SshKeyfile string `json:"sshKeyfile"`
	User       string `json:"user"`
	Password   string `json:"password"`
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
