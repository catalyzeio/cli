package crypto

// ICrypto
type ICrypto interface {
	DecryptFile(encryptedFilePath, key, iv, outputFilePath string) error
	EncryptFile(plainFilePath string, key, iv []byte) (string, error)
	Hex(src []byte) []byte
	Unhex(src []byte) []byte
	Base64Encode(src []byte) []byte
	Base64Decode(src []byte) []byte
}

// SCrypto is an implementor of ICrypto
type SCrypto struct{}

// New creates a new instance of ICrypto
func New() ICrypto {
	return &SCrypto{}
}
