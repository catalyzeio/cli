package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"os"
)

// EncryptFile takes in an open plaintext file and encrypts it to a temporary
// location based on the key and IV. It is up to the caller to ensure the
// encrypted file is deleted after it's used. The passed in key and iv should
// *NOT* be base64 encoded or hex encoded.
func (c *SCrypto) EncryptFile(plainFilePath string, key, iv []byte) (string, error) {
	plainFile, err := os.Open(plainFilePath)
	defer plainFile.Close()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	file, err := os.OpenFile(fmt.Sprintf("%s.encr", plainFilePath), os.O_CREATE|os.O_RDWR, 0600)
	defer file.Close()
	if err != nil {
		return "", err
	}

	// pack the file size in the first 8 bytes
	// would be nice to figure out how to do this and the Decrypt at least in a
	// similar manner
	stat, _ := plainFile.Stat()
	var origSize = uint64(stat.Size())
	sizeBytes := make([]byte, 8)
	binary.PutUvarint(sizeBytes, origSize)
	binary.Write(file, binary.LittleEndian, sizeBytes)

	encrypter := cipher.NewCBCEncrypter(block, iv)
	chunkSize := 24 * 1024
	for {
		chunk := make([]byte, chunkSize)
		read, _ := plainFile.Read(chunk)
		if read == 0 {
			break
		}
		encrChunk := make([]byte, chunkSize)
		encrypter.CryptBlocks(encrChunk, chunk)
		file.Write(encrChunk)
	}
	return file.Name(), nil
}
