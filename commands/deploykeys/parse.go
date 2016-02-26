package deploykeys

import (
	"crypto/rsa"
	"errors"

	"golang.org/x/crypto/ssh"
)

func (d *SDeployKeys) ParsePrivateKey(b []byte) (*rsa.PrivateKey, error) {
	in, err := ssh.ParseRawPrivateKey(b)
	if err != nil {
		return nil, errors.New("Invalid RSA private key format")
	}
	privKey := in.(*rsa.PrivateKey)
	return privKey, nil

}

func (d *SDeployKeys) ParsePublicKey(b []byte) (ssh.PublicKey, error) {
	s, _, _, _, err := ssh.ParseAuthorizedKey(b)
	if err != nil {
		return nil, errors.New("Invalid RSA public key format")
	}
	return s, nil
}

func (d *SDeployKeys) ExtractPublicKey(privKey *rsa.PrivateKey) (ssh.PublicKey, error) {
	s, err := ssh.NewPublicKey(&privKey.PublicKey)
	if err != nil {
		return nil, errors.New("Invalid RSA public key format derived from private key")
	}
	return s, nil
}
