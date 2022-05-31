package utility

import (
	"fmt"
	"testing"
)

func Test_EncryptPasswordField(t *testing.T) {
	connection := map[string]string{
		"server":     "localhost",
		"engine":     "mysql@dbgate-plugin-mysql",
		"sshMode":    "userPassword",
		"sshPort":    "22",
		"sshKeyfile": "/Users/liuliutiyong/.ssh/id_rsa",
		"user":       "root",
		"password":   "123456",
	}

	pwd := encryptPasswordField(connection, "password")["password"]
	fmt.Println(pwd)
	connection["password"] = pwd
	_ = decryptPasswordField(connection, "password")["password"]

}

func Test_LoadEncryptionKey(t *testing.T) {
	key2 := LoadEncryptionKey()
	fmt.Println(key2)
}

func Test_MaskConnection(t *testing.T) {
	fmt.Println(MaskConnection(map[string]string{
		"server":     "localhost",
		"engine":     "mysql@dbgate-plugin-mysql",
		"sshMode":    "userPassword",
		"sshPort":    "22",
		"sshKeyfile": "/Users/liuliutiyong/.ssh/id_rsa",
		"user":       "root",
		"password":   "123456",
	}))
}

func Test_PickSafeConnectionInfo(t *testing.T) {
	fmt.Println(PickSafeConnectionInfo(map[string]string{
		"server":     "localhost",
		"engine":     "mysql@dbgate-plugin-mysql",
		"sshMode":    "userPassword",
		"sshPort":    "22",
		"sshKeyfile": "/Users/liuliutiyong/.ssh/id_rsa",
		"user":       "root",
		"password":   "123456",
	}))
}
