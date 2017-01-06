package gcm

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"os"
)

const (
	chunkSize = 1024 * 1024 // 1 MiB

	// AAD (Additional authenticated data) is to be used in the GCM algorithm
	AAD = "7f57c07ee9459ed704d5e403086f6503"
)

// EncryptFile encrypts the file at the specified path using GCM.
func EncryptFile(inFilePath, outFilePath string, key, iv, aad []byte) error {
	if _, err := os.Stat(inFilePath); os.IsNotExist(err) {
		return fmt.Errorf("A file does not exist at %s", inFilePath)
	}

	inFile, err := os.Open(inFilePath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer outFile.Close()

	r, err := NewEncryptReader(inFile, key, iv, aad)
	if err != nil {
		return err
	}

	_, err = io.Copy(outFile, r)
	return err
}

// DecryptFile decrypts the file at the specified path using GCM.
func DecryptFile(inFilePath, outFilePath string, key, iv, aad []byte) error {
	if _, err := os.Stat(inFilePath); os.IsNotExist(err) {
		return fmt.Errorf("A file does not exist at %s", inFilePath)
	}

	inFile, err := os.Open(inFilePath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer outFile.Close()

	w, err := NewDecryptWriteCloser(outFile, key, iv, aad)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, inFile)
	return err
}

// Wraps data from an io.Reader in an encrypted GCM data stream.
type EncryptReader struct {
	src io.Reader
	eof bool

	gcm cipher.AEAD
	iv  []byte
	aad []byte

	sealed []byte
	off    int

	buff []byte
}

func NewEncryptReader(src io.Reader, key, iv, aad []byte) (*EncryptReader, error) {
	// copy the IV since it will be incremented
	ivCopy := make([]byte, len(iv))
	copy(ivCopy, iv)

	aes, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCMWithNonceSize(aes, len(iv))
	if err != nil {
		return nil, err
	}
	return &EncryptReader{
		src: src,

		gcm: gcm,
		iv:  ivCopy,
		aad: aad,

		sealed: []byte{},

		buff: make([]byte, chunkSize),
	}, nil
}

func (r *EncryptReader) Read(p []byte) (int, error) {
	n := len(p)
	off := 0
	for off < n {
		// encrypt the next chunk if no data available
		if r.off >= len(r.sealed) {
			if r.eof {
				return off, io.EOF
			}
			if err := r.seal(); err != nil {
				return off, err
			}
		}
		// copy encrypted data from the current chunk
		read := copy(p[off:], r.sealed[r.off:])
		r.off += read
		off += read
	}
	return off, nil
}

func (r *EncryptReader) CalculateTotalSize(size int) int {
	parts := size / chunkSize
	if size%chunkSize > 0 {
		parts++
	}
	return r.gcm.Overhead()*parts + size
}

func (r *EncryptReader) seal() error {
	// pull in the next chunk from the reader
	n, err := io.ReadFull(r.src, r.buff)
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		// mark EOF reached for subsequent seal attempts
		r.eof = true
	} else if err != nil {
		// some other problem; abort
		return err
	}
	// encrypt this chunk
	r.sealed = r.gcm.Seal(nil, r.iv, r.buff[:n], r.aad)
	incrementIV(r.iv)
	r.off = 0
	return nil
}

// Unwraps an encrypted GCM data to the given io.Writer stream.
type DecryptWriteCloser struct {
	dst io.WriteCloser

	gcm cipher.AEAD
	iv  []byte
	aad []byte

	sealed []byte
	off    int
}

func NewDecryptWriteCloser(dst io.WriteCloser, key, iv, aad []byte) (*DecryptWriteCloser, error) {
	// copy the IV since it will be incremented
	ivCopy := make([]byte, len(iv))
	copy(ivCopy, iv)

	aes, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCMWithNonceSize(aes, len(iv))
	if err != nil {
		return nil, err
	}

	return &DecryptWriteCloser{
		dst: dst,

		gcm: gcm,
		iv:  ivCopy,
		aad: aad,

		sealed: make([]byte, chunkSize+gcm.Overhead()),
	}, nil
}

func (w *DecryptWriteCloser) Write(p []byte) (int, error) {
	n := len(p)
	off := 0
	for off < n {
		// copy encrypted data into the current chunk
		written := copy(w.sealed[w.off:], p[off:])
		w.off += written
		off += written
		// decrypt if chunk is finished
		if w.off == cap(w.sealed) {
			if err := w.open(); err != nil {
				return off, err
			}
		}
	}
	return off, nil
}

func (w *DecryptWriteCloser) Close() error {
	if w.off > 0 {
		if err := w.open(); err != nil {
			return err
		}
	}
	return w.dst.Close()
}

func (w *DecryptWriteCloser) open() error {
	opened, err := w.gcm.Open(nil, w.iv, w.sealed[:w.off], w.aad)
	if err != nil {
		return err
	}
	if _, err := w.dst.Write(opened); err != nil {
		return err
	}
	incrementIV(w.iv)
	w.off = 0
	return nil
}

func incrementIV(iv []byte) {
	for i := len(iv) - 1; i >= 0; i-- {
		iv[i]++
		if iv[i] != 0 {
			return
		}
	}
}
