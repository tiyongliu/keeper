package internal

import (
	"encoding/json"
	"github.com/samber/lo"
	"io/ioutil"
	"keeper/app/pkg/logger"
	"keeper/app/utility"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

const (
	defaultEncryptionKey = "mQAUaXhavRGJDxDTXSCg7Ej0xMmGCrx6OKA07DIMBiDcYYkvkaXjTAzPUEHEHEf9"
	encryptionKeyKey     = "encryptionKey"
)

var _encryptionKey string

func LoadEncryptionKey() string {
	if _encryptionKey != "" {
		return _encryptionKey
	}
	defaultFile := utility.DataDirCore()
	keyFile := filepath.Join(defaultFile, ".key")
	encryptor := CreateEncryptor(defaultEncryptionKey)
	if !utility.IsExist(keyFile) {
		if err := os.MkdirAll(filepath.Dir(keyFile), os.ModePerm); err != nil {
			log.Fatalf("os.MkdirAll failed err: %v\n", err)
			return ""
		}

		generatedKey := randomBytes(32)
		newKey := string(generatedKey)
		result := map[string]string{
			encryptionKeyKey: newKey,
		}
		encrypt := encryptor.encrypt(result)

		if err := ioutil.WriteFile(keyFile, []byte(encrypt), os.ModePerm); err != nil {
			log.Fatalf("ioutil.WriteFile failed err: %v\n", err)
			return ""
		}
	}

	encryptedData, err := ioutil.ReadFile(keyFile)
	if err != nil {
		log.Fatalf("ioutil.ReadFile failed err: %v\n", err)
		return ""
	}

	decrypt := encryptor.decrypt(string(encryptedData))
	data := map[string]string{}
	err = json.Unmarshal([]byte(decrypt), &data)
	if err != nil {
		log.Fatalf("json.Unmarshal failed err: %v\n", err)
		return ""
	}
	_encryptionKey = data[encryptionKeyKey]
	return _encryptionKey
}

var _encryptor *SimpleEncryptor

func getEncryptor() *SimpleEncryptor {
	if _encryptor != nil {
		return _encryptor
	}

	_encryptor = CreateEncryptor(LoadEncryptionKey())
	return _encryptor
}

func encryptPasswordField(connection map[string]string, field string) map[string]string {
	if connection != nil &&
		connection[field] != "" &&
		!strings.HasPrefix(connection[field], "crypt:") &&
		connection["passwordMode"] != "saveRaw" {
		connection[field] = "crypt:" + getEncryptor().encrypt(connection[field])
	}
	return connection
}

/*func decryptPasswordField(connection map[string]string, field string) map[string]string {
	if connection != nil &&
		connection[field] != "" &&
		strings.HasPrefix(connection[field], "crypt:") {
		decrypt := getEncryptor().decrypt(strings.Split(connection[field], "crypt:")[1])
		//TODO "123456" 需要去掉前后的“”, 当前只是剔除前后空格
		connection[field] = decrypt[1 : len(decrypt)-1]
	}

	return connection
}*/

func decryptPasswordField(connection map[string]interface{}, field string) map[string]interface{} {
	if connection != nil && connection[field] != nil && reflect.ValueOf(connection[field]).Kind() == reflect.String {
		logger.Infof("password.reflect: %s", connection[field])
		value := connection[field].(string)
		if field != "" && strings.HasPrefix(value, "crypt:") {
			decrypt := getEncryptor().decrypt(strings.Split(value, "crypt:")[1])
			connection[field] = strings.Trim(decrypt[1:len(decrypt)-1], "")
		}
	}

	return connection
}

func EncryptConnection(connection map[string]string) map[string]string {
	connection = encryptPasswordField(connection, "password")
	connection = encryptPasswordField(connection, "sshPassword")
	connection = encryptPasswordField(connection, "sshKeyfilePassword")
	return connection
}

func MaskConnection(connection map[string]string) map[string]string {
	if len(connection) == 0 {
		return connection
	}
	return utility.MapOmit(connection, []string{"password", "sshPassword", "sshKeyfilePassword"})
}

func DecryptConnection(connection map[string]interface{}) map[string]interface{} {
	connection = decryptPasswordField(connection, "password")
	connection = decryptPasswordField(connection, "sshPassword")
	connection = decryptPasswordField(connection, "sshKeyfilePassword")
	return connection
}

func PickSafeConnectionInfo(connection map[string]interface{}) map[string]interface{} {
	return lo.MapValues(connection, func(v interface{}, k string) interface{} {
		if k == "engine" || k == "port" || k == "authType" || k == "sshMode" || k == "passwordMode" {
			return v
		}
		if v != nil {
			return "***"
		}
		return nil
	})
}
