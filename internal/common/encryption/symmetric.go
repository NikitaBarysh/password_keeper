// Package encryption - пакет, в котором происходит шифрование данных
package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"fmt"
)

// SymmetricEncrypt - симметричное шифрование данных
func SymmetricEncrypt(src []byte, hashKey string) ([]byte, error) {
	if len(hashKey) == 0 {
		return nil, fmt.Errorf("SymmetricEncrypt: hashKey is required")
	}

	key := sha256.Sum256([]byte(hashKey))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("SymmetricEncrypt: error creating AES cipher: %v", err)
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, fmt.Errorf("SymmetricEncrypt: error creating GCM: %v", err)
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	dst := aesgcm.Seal(nil, nonce, src, nil)

	return dst, nil
}

// SymmetricDecrypt - симметричная расшифровка данных
func SymmetricDecrypt(data []byte, hashKey string) ([]byte, error) {
	key := sha256.Sum256([]byte(hashKey))

	aesblock, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, fmt.Errorf("SymmetricDecrypt: error creating AES cipher: %v", err)
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, fmt.Errorf("SymmetricDecrypt: error creating GCM: %v", err)
	}

	nonce := key[len(key)-aesgcm.NonceSize():]

	decryptedData, err := aesgcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, fmt.Errorf("SymmetricDecrypt: error decrypting data: %v", err)
	}

	return decryptedData, nil
}
