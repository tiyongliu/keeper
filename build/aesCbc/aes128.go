package aesCbc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

const (
	BLOCK_SIZE_16 = 16
)

func NewAesCipher128(key, iv []byte) *AesCipher128 {
	if len(key) == 0 || len(key) > BLOCK_SIZE_16 {
		return nil
	}
	if len(iv) < BLOCK_SIZE_16 {
		newIv := make([]byte, BLOCK_SIZE_16)
		copy(newIv, iv)
		iv = newIv
	} else {
		iv = iv[:BLOCK_SIZE_16]
	}
	newKey := make([]byte, BLOCK_SIZE_16)
	copy(newKey, key)

	block, err := aes.NewCipher(newKey)
	if err != nil {
		return nil
	}
	return &AesCipher128{
		key:   newKey,
		iv:    iv,
		block: block,
	}
}

type AesCipher128 struct {
	key   []byte
	iv    []byte
	block cipher.Block
}

func (aesCipher *AesCipher128) BlockSize() int {
	return BLOCK_SIZE_16
}

func (aesCipher *AesCipher128) Encrypt(origData []byte) []byte {
	encodeBytes := []byte(origData)
	blockSize := aesCipher.BlockSize()
	encodeBytes = padding(encodeBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(aesCipher.block, aesCipher.iv)
	crypted := make([]byte, len(encodeBytes))
	blockMode.CryptBlocks(crypted, encodeBytes)

	return crypted
}

func (aesCipher *AesCipher128) Decrypt(encrData []byte) []byte {
	blockMode := cipher.NewCBCDecrypter(aesCipher.block, aesCipher.iv)
	result := make([]byte, len(encrData))
	blockMode.CryptBlocks(result, encrData)
	return bytes.Trim(result, "\x00")
}

func padding(ciphertext []byte, blockSize int) []byte {
	dataSize := ((len(ciphertext)-1)/blockSize + 1) * blockSize
	if dataSize == len(ciphertext) {
		return ciphertext
	}
	newData := make([]byte, dataSize)
	copy(newData, ciphertext)
	return newData
}
