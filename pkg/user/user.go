package user

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"time"
)

const (
	accessKeyLength = 20
	validDays       = 30
)

// User holds user information.
type User struct {
	Name      string    `yaml:"name" json:"name"`
	AccessKey AccessKey `yaml:"accessKey" json:"accessKey"`
}

// AccessKey holds access key and validity information.
type AccessKey struct {
	Key     string    `yaml:"key" json:"key"`
	ValidTo time.Time `yaml:"validTo" json:"validTo"`
}

// New creats a new user and access key.
func New(username string) User {
	accessKey := NewAccessKey()
	return User{
		Name:      username,
		AccessKey: accessKey,
	}
}

// NewAccessKey creates a new access key.
func NewAccessKey() AccessKey {
	return AccessKey{
		Key:     createNewKey(),
		ValidTo: createValidTo(),
	}
}

func (key AccessKey) Valid() bool {
	return key.ValidTo.After(time.Now().UTC())
}

func createNewKey() string {
	bytes := make([]byte, accessKeyLength)
	_, err := rand.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(bytes))
	return hash[:accessKeyLength]
}

func createValidTo() time.Time {
	return time.Now().UTC().AddDate(0, 0, validDays)
}
