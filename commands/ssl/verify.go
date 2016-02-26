package ssl

import (
	"crypto/dsa"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
)

func CmdVerify(chainPath, privateKeyPath, hostname string, selfSigned bool, is ISSL) error {
	if _, err := os.Stat(chainPath); os.IsNotExist(err) {
		return fmt.Errorf("A cert does not exist at path '%s'", chainPath)
	}
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("A private key does not exist at path '%s'", privateKeyPath)
	}
	err := is.Verify(chainPath, privateKeyPath, hostname, selfSigned)
	if err != nil {
		return err
	}
	logrus.Println("Certificate chain and key are valid")
	return nil
}

func IsIncompleteChainErr(err error) bool {
	switch err.(type) {
	case nil:
		return false
	case *IncompleteChainError:
		return true
	default:
		return false
	}
}

func IsHostnameMismatchErr(err error) bool {
	switch err.(type) {
	case nil:
		return false
	case *HostnameMismatchError:
		return true
	default:
		return false
	}
}

// Verify takes a chain and ensures it is a full chain optionally ensuring
// the given private key matches that chain.
func (s *SSSL) Verify(chainPath, privateKeyPath, hostname string, selfSigned bool) error {
	cert, err := tls.LoadX509KeyPair(chainPath, privateKeyPath)
	if err != nil {
		return err
	}
	if !selfSigned {
		x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
		if err != nil {
			return err
		}
		certPool := x509.NewCertPool()
		for i := 1; i < len(cert.Certificate); i++ {
			c, err := x509.ParseCertificate(cert.Certificate[i])
			if err != nil {
				return err
			}
			certPool.AddCert(c)
		}

		// verify we can make a chain from cert down through the intermediates to a root
		if _, err := x509Cert.Verify(x509.VerifyOptions{
			Intermediates: certPool,
		}); err != nil {
			return &IncompleteChainError{
				Err:     err,
				Message: "Failed to verify certificate chain",
			}
		}
		// verify the cert we pulled out matches the hostname specified
		if err := x509Cert.VerifyHostname(hostname); err != nil {
			return &HostnameMismatchError{
				Err:     err,
				Message: "Certificate hostname mismatch",
			}
		}
		outputCertInfo(x509Cert)
		warnOnExpired(x509Cert)
	}
	return nil
}

func outputCertInfo(cert *x509.Certificate) {
	logrus.Printf("Issued by: %s", cert.Issuer.CommonName)
	logrus.Printf("Subject: %s", cert.Subject.CommonName)
	switch cert.SignatureAlgorithm {
	case x509.UnknownSignatureAlgorithm:
		logrus.Println("Signature Algorithm: Unknown")
	case x509.MD2WithRSA:
		logrus.Println("Signature Algorithm: MD2 with RSA")
	case x509.MD5WithRSA:
		logrus.Println("Signature Algorithm: MD5 with RSA")
	case x509.SHA1WithRSA:
		logrus.Println("Signature Algorithm: SHA 1 with RSA")
	case x509.SHA256WithRSA:
		logrus.Println("Signature Algorithm: SHA 256 with RSA")
	case x509.SHA384WithRSA:
		logrus.Println("Signature Algorithm: SHA 384 with RSA")
	case x509.SHA512WithRSA:
		logrus.Println("Signature Algorithm: SHA 512 with RSA")
	case x509.DSAWithSHA1:
		logrus.Println("Signature Algorithm: DSA with SHA 1")
	case x509.DSAWithSHA256:
		logrus.Println("Signature Algorithm: DSA with SHA 256")
	case x509.ECDSAWithSHA1:
		logrus.Println("Signature Algorithm: ECDSA with SHA 1")
	case x509.ECDSAWithSHA256:
		logrus.Println("Signature Algorithm: ECDSA with SHA 256")
	case x509.ECDSAWithSHA384:
		logrus.Println("Signature Algorithm: ECDSA with SHA 384")
	case x509.ECDSAWithSHA512:
		logrus.Println("Signature Algorithm: ECDSA with SHA 512")
	}
	switch cert.PublicKeyAlgorithm {
	case x509.UnknownPublicKeyAlgorithm:
	case x509.RSA:
		logrus.Println("Public Key Algorithm: RSA")
		publicKey := cert.PublicKey.(*rsa.PublicKey)
		logrus.Printf("Key Size: %d", publicKey.N.BitLen())
	case x509.DSA:
		logrus.Println("Public Key Algorithm: DSA")
		publicKey := cert.PublicKey.(*dsa.PublicKey)
		logrus.Printf("Key Size: %d", publicKey.Y.BitLen())
	case x509.ECDSA:
		logrus.Println("Public Key Algorithm: ECDSA")
		publicKey := cert.PublicKey.(*ecdsa.PublicKey)
		logrus.Printf("Key Size: %d", publicKey.Y.BitLen())
	}
	logrus.Printf("Not Valid Before: %s", cert.NotBefore.Local().String())
	logrus.Printf("Not Valid After: %s", cert.NotAfter.Local().String())
	logrus.Println()
}

func warnOnExpired(cert *x509.Certificate) {
	if time.Since(cert.NotBefore) < 0 {
		logrus.Println("WARNING! This certificate is not yet valid!")
	}
	if time.Since(cert.NotAfter) > 0 {
		logrus.Println("WARNING! This certificate is expired!")
	}
}
