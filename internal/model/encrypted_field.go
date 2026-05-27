package model

import (
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/miaomiaopu/ipmp-server/internal/pkg/crypto"
)

// EncryptedField 加密字段类型
// 写入时 AES-256-GCM 加密 → hex 编码存储
// 读取时 hex 解码 → AES-256-GCM 解密
// JSON 序列化时自动掩码（邮箱 a***@d***.com / 手机 138****8000）
type EncryptedField string

// Scan 实现 sql.Scanner，从数据库读取时自动解密
func (e *EncryptedField) Scan(value interface{}) error {
	if value == nil {
		*e = ""
		return nil
	}
	var cipherHex string
	switch v := value.(type) {
	case []byte:
		cipherHex = string(v)
	case string:
		cipherHex = v
	default:
		return fmt.Errorf("unsupported type for EncryptedField: %T", value)
	}
	if cipherHex == "" {
		*e = ""
		return nil
	}
	cipherBytes, err := hex.DecodeString(cipherHex)
	if err != nil {
		return fmt.Errorf("decode encrypted field: %w", err)
	}
	plaintext, err := crypto.Decrypt(cipherBytes)
	if err != nil {
		return fmt.Errorf("decrypt encrypted field: %w", err)
	}
	*e = EncryptedField(plaintext)
	return nil
}

// Value 实现 driver.Valuer，写入数据库时自动加密为 hex 字符串
func (e EncryptedField) Value() (driver.Value, error) {
	if e == "" {
		return "", nil
	}
	ciphertext, err := crypto.Encrypt([]byte(e))
	if err != nil {
		return nil, fmt.Errorf("encrypt encrypted field: %w", err)
	}
	return hex.EncodeToString(ciphertext), nil
}

// MarshalJSON JSON 序列化时自动掩码
func (e EncryptedField) MarshalJSON() ([]byte, error) {
	s := string(e)
	if s == "" {
		return []byte(`""`), nil
	}
	masked := maskValue(s)
	return []byte(fmt.Sprintf(`"%s"`, masked)), nil
}

// maskValue 对敏感值进行掩码
// len≤4: 仅保留首字符 → "张***"
// len≤8: 保留首尾各2字符 → "ab***cd"
// len>8: 保留前3后4字符 → "138****8000"
func maskValue(s string) string {
	runes := []rune(s)
	if len(runes) <= 4 {
		return string(runes[0:1]) + "***"
	}
	if len(runes) <= 8 {
		return string(runes[0:2]) + "***" + string(runes[len(runes)-2:])
	}
	return string(runes[0:3]) + "****" + string(runes[len(runes)-4:])
}

var _ driver.Valuer = EncryptedField("")
var _ interface{ Scan(interface{}) error } = (*EncryptedField)(nil)

// ErrCryptoNotReady 加密模块未初始化
var ErrCryptoNotReady = errors.New("crypto module not initialized")
