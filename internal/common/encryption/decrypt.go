package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

var dataDecryptor *Decryptor

type Decryptor struct {
	privateKey *rsa.PrivateKey
}

func InitDecryptor(file string) error {
	key, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("err to read file: %w", err)
	}

	keyBlock, _ := pem.Decode([]byte(key))
	if keyBlock == nil {
		return fmt.Errorf("err to decode private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("err to parse privat key: %w", err)
	}

	dataDecryptor = &Decryptor{privateKey: privateKey}

	return nil
}

func (m *Decryptor) Decrypt(msg []byte) ([]byte, error) {
	msgLen := len(msg)
	hash := sha512.New()
	random := rand.Reader

	step := m.privateKey.Size()
	var decryptedBytes []byte

	for start := 0; start < msgLen; start += step {
		finish := start + step
		if finish > msgLen {
			finish = msgLen
		}

		decryptedBlockBytes, err := rsa.DecryptOAEP(hash, random, m.privateKey, msg[start:finish], nil)
		if err != nil {
			return nil, fmt.Errorf("decrypt part message process error: %w", err)
		}

		decryptedBytes = append(decryptedBytes, decryptedBlockBytes...)
	}

	return decryptedBytes, nil
}
