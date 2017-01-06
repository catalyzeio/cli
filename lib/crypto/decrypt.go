package crypto

import (
	"io"

	"github.com/catalyzeio/gcm/gcm"
)

// DecryptFile takes in an ecrypted file and decrypts it to the given
// output path based on the Key and IV. The Key and IV should be the hex and
// base64 encoded version
func (c *SCrypto) DecryptFile(encryptedFilePath, key, iv, outputFilePath string) error {
	return gcm.DecryptFile(encryptedFilePath, outputFilePath, c.Unhex([]byte(key), KeySize), c.Unhex([]byte(iv), IVSize), c.Unhex([]byte(gcm.AAD), AADSize))
}

// NewDecryptWriteCloser takes a io.WriteCloser and wraps it in a
// type that will decrypt Writes to the io.WriteCloser as they are written.
// The passed in key and iv should *NOT* be base64 encoded or hex encoded.
func (c *SCrypto) NewDecryptWriteCloser(writeCloser io.WriteCloser, key, iv string) (*gcm.DecryptWriteCloser, error) {
	return gcm.NewDecryptWriteCloser(writeCloser, c.Unhex([]byte(key), KeySize), c.Unhex([]byte(iv), IVSize), c.Unhex([]byte(gcm.AAD), AADSize))
}
