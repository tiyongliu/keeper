package aesCbc

import (
	"bytes"
	"unsafe"
)

const (
	BLOCK_SIZE_32 = 32
)

var keySizes = []int{16, 24, 32}

type cbcBuffer struct {
	previousCiphertext []uint32
	previousCipher     []uint32
	blocksize          int
}

type RI struct {
	Nk   int
	Nb   int
	Nr   int
	fi   [24]byte
	ri   [24]byte
	fkey [120]uint32
	rkey [120]uint32
}

type AesCipher256 struct {
	/* Holds the algorithm's internal key */
	rinst  *RI
	buffer *cbcBuffer /* holds the mode's internal buffers */

	/* holds the key */
	key []byte

	/* These were included to speed up encryption/decryption proccess, so
	 * there is not need for resolving symbols every time.
	 */
	blockSize int
}

func NewAesCipher256(key, iv []byte) *AesCipher256 {
	td := mcryptOpen()
	mcryptGenericInit(td, key, iv)
	return td
}

func (aesCipher *AesCipher256) Encrypt() []byte {
	dataSize := ((len(origData)-1)/aesCipher.blockSize + 1) * aesCipher.blockSize
	newData := make([]byte, dataSize)
	copy(newData, origData)
	return mcrypt(aesCipher.buffer, newData, aesCipher.blockSize, aesCipher.rinst)
}

func (aesCipher *AesCipher256) Decrypt(origData []byte) []byte {
	return mdecrypt(aesCipher.buffer, origData, aesCipher.blockSize, aesCipher.rinst)
}

func (aesCipher *AesCipher256) BlockSize() int {
	return aesCipher.blockSize
}

//获取加密key长度
func getKeySize(size int) int {
	for _, val := range keySizes {
		if size <= val {
			return val
		}
	}
	return 32
}

func getBlockSize() int {
	return BLOCK_SIZE_32
}

func mcryptOpen() *AesCipher256 {
	td := &AesCipher256{}
	td.blockSize = getBlockSize()
	td.buffer = &cbcBuffer{}
	td.rinst = &RI{}
	return td
}

//初始化秘钥
func mcryptGenericInit(td *AesCipher256, key, iv []byte) int {
	keySize := len(key)
	if keySize == 0 || keySize > td.blockSize {
		return -1
	}
	if len(iv) < td.blockSize {
		newIv := make([]byte, td.blockSize)
		copy(newIv, iv)
		iv = newIv
	} else {
		iv = iv[:td.blockSize]
	}
	td.key = make([]byte, td.blockSize)
	copy(td.key, key)
	initMcrypt(td.buffer, iv, td.blockSize)
	mcryptSetKey(td.rinst, td.key, getKeySize(keySize))
	return 0
}

func initMcrypt(buf *cbcBuffer, iv []byte, size int) {
	buf.blocksize = size
	buf.previousCiphertext = make([]uint32, size/4)
	buf.previousCipher = make([]uint32, size/4)
	index := 0
	for i := 0; i < size; i += 4 {
		buf.previousCiphertext[index] = pack(iv[i:])
		index++
	}
}

//cbc加密
func mcrypt(buf *cbcBuffer, plaintext []byte, blocksize int, rinst *RI) []byte {
	var plain []uint32
	intSize := 4
	txtLen := len(plaintext)
	fplain := *((*[]uint32)(unsafe.Pointer(&plaintext)))
	onceLen := blocksize / intSize //8
	for j := 0; j < txtLen/blocksize; j++ {
		plain = fplain[j*blocksize/intSize:]
		for i := 0; i < onceLen; i++ {
			plain[i] ^= buf.previousCiphertext[i]
		}
		//每次加密32字节
		mcryptEncrypt(rinst, *((*[]byte)(unsafe.Pointer(&plain))))
		/* Copy the ciphertext to prev_ciphertext */
		copy(buf.previousCiphertext, plain)
	}
	return plaintext
}

//cbc解密
func mdecrypt(buf *cbcBuffer, ciphertext []byte, blocksize int, rinst *RI) []byte {
	var cipher []uint32
	fcipher := *((*[]uint32)(unsafe.Pointer(&ciphertext)))
	txtLen := len(ciphertext)
	intSize := 4
	for j := 0; j < txtLen/blocksize; j++ {
		cipher = fcipher[j*blocksize/intSize:]
		copy(buf.previousCipher, cipher)
		mcryptDecrypt(rinst, *((*[]byte)(unsafe.Pointer(&cipher))))
		for i := 0; i < blocksize/intSize; i++ {
			cipher[i] ^= buf.previousCiphertext[i]
		}
		/* Copy the ciphertext to prev_cipher */
		copy(buf.previousCiphertext, buf.previousCipher)
	}
	return bytes.Trim(ciphertext, "\x00")
}
