package internal

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	json2 "encoding/json"
	"fmt"
	"github.com/Luzifer/go-openssl/v4"
	"log"
	"math/rand"
)

func randomBytes(size int) []byte {
	iv := make([]byte, size)
	rand.Read(iv)
	return iv
}

const min_key_length = 16

type Option struct {
	Key        string
	VerifyHmac bool
	Debug      bool
	Reviver    interface{}
}

type SimpleEncryptor struct {
	SSL       *openssl.OpenSSL
	CryptoKey string
	Opts      *Option
}

func CreateEncryptor(opts interface{}) *SimpleEncryptor {
	var o *Option
	switch opts.(type) {
	case string:
		o = &Option{
			Key:        opts.(string),
			VerifyHmac: true,
			Debug:      false,
		}
	case Option:
		o = opts.(*Option)
	default:
		log.Fatalln("invalid opts")
	}

	if o == nil || o.Key == "" {
		log.Fatalln("a string key must be specified")
	}
	if len(o.Key) < min_key_length {
		log.Fatalln("key must be at least ' + MIN_KEY_LENGTH + ' characters long")
	}

	h := sha256.New()
	h.Write([]byte(o.Key))

	return &SimpleEncryptor{
		SSL:       openssl.New(),
		CryptoKey: fmt.Sprintf("%x\n", h.Sum(nil)),
		Opts:      o,
	}
}

func (simple *SimpleEncryptor) encrypt(obj interface{}) string {
	json, err := json2.Marshal(&obj)
	if err != nil {
		return ""
	}

	passphrase := randomBytes(16)

	cipher, err := simple.SSL.EncryptBinaryBytes(string(passphrase), json, openssl.PBKDF2SHA256)
	if err != nil {
		return ""
	}

	var buffer bytes.Buffer
	buffer.Write(passphrase)
	buffer.Write(cipher)

	encryptedJson := base64.RawStdEncoding.EncodeToString(buffer.Bytes())
	return hmacSha256(string(cipher), simple.CryptoKey) + encryptedJson
}

func (simple *SimpleEncryptor) decrypt(opensslEncrypted string) string {
	expectedHmac := opensslEncrypted[:64]

	buffer := opensslEncrypted[64:]
	encryptedJson, err := base64.RawStdEncoding.DecodeString(buffer)
	if err != nil {
		return ""
	}

	passphrase := encryptedJson[:16]
	cipher := encryptedJson[16:]
	if !hmacCheck(cipher, []byte(expectedHmac), []byte(simple.CryptoKey)) {
		//HMAC does not match
		return ""
	}

	resBytes, err := simple.SSL.DecryptBinaryBytes(string(passphrase), cipher, openssl.PBKDF2SHA256)
	if err != nil {
		return ""
	}

	if simple.Opts.Reviver != nil {
	}

	return fmt.Sprintf("%s", resBytes)
}

func hmacCheck(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha256.New, key)
	mac.Write(message)
	expectedMAC := hex.EncodeToString(mac.Sum([]byte(nil)))
	return hmac.Equal(messageMAC, []byte(expectedMAC))
}

func hmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
