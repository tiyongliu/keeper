package modules

type SimpleSettingMysql struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
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
