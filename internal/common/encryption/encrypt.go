// Package encryption - пакет, в котором происходит шифрование данных
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

// Encryptor - структура в которой лежит публичный клюс
type Encryptor struct {
	openKey *rsa.PublicKey
}

// InitEncryptor - функция создает публичный ключ
func InitEncryptor(file string) (*Encryptor, error) {
	key, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("InitEncryptor: failed to read encryption file: %w", err)
	}

	keyBlock, _ := pem.Decode(key)
	if keyBlock == nil {
		return nil, fmt.Errorf("InitEncryptor: empty key")
	}

	pubKey, err := x509.ParsePKCS1PublicKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("InitEncryptor: err to parse key: %w", err)
	}

	return &Encryptor{openKey: pubKey}, nil
}

// Encrypt - метод для зашифровки данных
func (m *Encryptor) Encrypt(msg []byte) ([]byte, error) {

	hash := sha512.New()
	random := rand.Reader

	step := m.openKey.Size() - 2*hash.Size() - 2
	var encryptedBytes []byte

	for start := 0; start < len(msg); start += step {
		finish := start + step
		if finish > len(msg) {
			finish = len(msg)
		}

		encryptedBlockBytes, err := rsa.EncryptOAEP(hash, random, m.openKey, msg[start:finish], nil)
		if err != nil {
			return nil, fmt.Errorf("Encrypt: encrypt part message process error: %w ", err)
		}

		encryptedBytes = append(encryptedBytes, encryptedBlockBytes...)
	}
	return encryptedBytes, nil
}
