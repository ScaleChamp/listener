package models

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql/driver"
	"encoding/base64"
	"encoding/hex"
	uuid "github.com/satori/go.uuid"
	"os"
)

type Instance struct {
	Id             uuid.UUID `json:"id"`
	Name           string    `json:"name"`

	Kind           string    `json:"kind"`

	State          int       `json:"state"`
	Whitelist      []string  `json:"whitelist"`
	EvictionPolicy string    `json:"eviction_policy"`
	LicenseKey     string    `json:"license_key"`
	Database       string    `json:"database"`
	Username       string    `json:"username"`
	Password       string    `json:"password"`
	Secret         string    `json:"secret"`
	PlanId         uuid.UUID `json:"plan_id"`
	ProjectId      uuid.UUID `json:"project_id"`
}

type EncryptedString string

func (m EncryptedString) Value() (driver.Value, error) {
	key, err := hex.DecodeString(os.Getenv("LOCKBOX_MASTER_KEY"))
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return aesgcm.Seal(nil, nonce, []byte(m), nil), nil
}

func (m *EncryptedString) Scan(src interface{}) error {
	key, err := hex.DecodeString(os.Getenv("LOCKBOX_MASTER_KEY"))
	if err != nil {
		return err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	ciphertext, err := base64.RawStdEncoding.DecodeString(src.(string))
	if err != nil {
		return err
	}
	plaintext, err := gcm.Open(nil, ciphertext[:12], ciphertext[12:], nil)
	if err != nil {
		return err
	}
	*m = EncryptedString(plaintext)
	return nil
}

//fmt.Println(string(plaintext), base64.RawStdEncoding.EncodeToString(aesgcm2.Seal(ciphertext[:12], ciphertext[:12], plaintext, nil)))
