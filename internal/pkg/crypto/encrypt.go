package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"sync"
)

var (
	aead     cipher.AEAD
	aeadOnce sync.Once
	aeadErr  error
)

// Init 初始化 AES-256-GCM，key 为 32 字节 hex 字符串(64字符)
func Init(keyHex string) error {
	aeadOnce.Do(func() {
		key, err := hex.DecodeString(keyHex)
		if err != nil {
			aeadErr = fmt.Errorf("decode encryption key: %w", err)
			return
		}
		if len(key) != 32 {
			aeadErr = errors.New("encryption key must be 32 bytes (64 hex chars)")
			return
		}
		block, err := aes.NewCipher(key)
		if err != nil {
			aeadErr = fmt.Errorf("create aes cipher: %w", err)
			return
		}
		aead, aeadErr = cipher.NewGCM(block)
	})
	return aeadErr
}

// Encrypt AES-256-GCM 加密，nonce(12字节)前置于密文
func Encrypt(plaintext []byte) ([]byte, error) {
	if aead == nil {
		return nil, errors.New("crypto not initialized, call Init first")
	}
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}
	ciphertext := aead.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Decrypt AES-256-GCM 解密
func Decrypt(ciphertext []byte) ([]byte, error) {
	if aead == nil {
		return nil, errors.New("crypto not initialized, call Init first")
	}
	nonceSize := aead.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}
	nonce, ct := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aead.Open(nil, nonce, ct, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}
	return plaintext, nil
}
