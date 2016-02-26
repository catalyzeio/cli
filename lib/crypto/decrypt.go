package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"fmt"
	"os"
)

// DecryptFile takes in an ecrypted file and decrypts it to the given
// output path based on the Key and IV. The Key and IV should be the hex and
// base64 encoded version
func (c *SCrypto) DecryptFile(encryptedFilePath, key, iv, outputFilePath string) error {
	encryptedFile, err := os.Open(encryptedFilePath)
	defer encryptedFile.Close()
	if err != nil {
		return err
	}
	sizeBytes := make([]byte, 8)
	binary.Read(encryptedFile, binary.LittleEndian, sizeBytes)
	var origSize int64
	binary.Read(bytes.NewBuffer(sizeBytes), binary.LittleEndian, &origSize)
	block, err := aes.NewCipher(c.Unhex(c.Base64Decode([]byte(key))))
	if err != nil {
		return err
	}

	file, err := os.OpenFile(outputFilePath, os.O_RDWR, 0600)
	if os.IsNotExist(err) {
		file, err = os.OpenFile(outputFilePath, os.O_CREATE|os.O_RDWR, 0600)
	}
	if err != nil {
		return err
	}
	defer file.Close()

	decrypter := cipher.NewCBCDecrypter(block, c.Unhex(c.Base64Decode([]byte(iv))))
	chunkSize := 24 * 1024
	for {
		chunk := make([]byte, chunkSize)
		read, _ := encryptedFile.Read(chunk)
		if read%aes.BlockSize != 0 {
			return fmt.Errorf("Logs unavailable for this job")
		}
		if read == 0 {
			break
		}
		plainChunk := make([]byte, read)
		decrypter.CryptBlocks(plainChunk, chunk[:read]) // only decrypt the amount we read
		file.Write(plainChunk)
	}
	file.Truncate(origSize)
	return nil
}
