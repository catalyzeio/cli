package ssl

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/zakjan/cert-chain-resolver/certUtil"
)

func CmdResolve(chainPath, privateKeyPath, hostname, outputPath string, force bool, is ISSL) error {
	if _, err := os.Stat(chainPath); os.IsNotExist(err) {
		return fmt.Errorf("A cert does not exist at path '%s'", chainPath)
	}
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("A private key does not exist at path '%s'", privateKeyPath)
	}
	if !force && outputPath != "" {
		if _, err := os.Stat(outputPath); err == nil {
			return fmt.Errorf("File already exists at path '%s'. Specify `--force` to overwrite", outputPath)
		}
	}
	err := is.Verify(chainPath, privateKeyPath, hostname, false)
	if err == nil {
		logrus.Println("Certificate chain and key are valid and complete")
		return nil
	} else if !IsIncompleteChainErr(err) {
		return err
	}
	data, err := is.Resolve(chainPath)
	file := os.Stdout
	if outputPath != "" {
		os.Remove(outputPath)
		file, err = os.OpenFile(outputPath, os.O_CREATE|os.O_RDWR, 0400)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	file.Write(data)
	if outputPath != "" {
		logrus.Printf("Resolved certificate chain and wrote the contents to '%s'", outputPath)
	}
	return nil
}

func (s *SSSL) Resolve(chainPath string) ([]byte, error) {
	logrus.Println("Incomplete certificate chain found, attempting to resolve this")
	b, err := ioutil.ReadFile(chainPath)
	if err != nil {
		return nil, err
	}

	cert, err := certUtil.DecodeCertificate(b)
	if err != nil {
		return nil, err
	}

	certs, err := certUtil.FetchCertificateChain(cert)
	if err != nil {
		return nil, err
	}

	return certUtil.EncodeCertificates(certs), nil
}
