package user

import (
	"testing"
	"time"
)

func TestNewUser(t *testing.T) {
	user := New("test-usr")
	if user.Name != "test-usr" {
		t.Errorf("New user name wrong. Expected=test-usr Got=%s", user.Name)
	}
	if !user.AccessKey.Valid() {
		t.Errorf("New user accesskey timestamp wrong, is before current timestamp")
	}
}

func TestNewAccessKey(t *testing.T) {
	accessKey := NewAccessKey()
	if !accessKey.Valid() {
		t.Errorf("AccessKey timestamp wrong, is before current timestamp")
	}
	if len(accessKey.Key) != accessKeyLength {
		t.Errorf("Incorrect AccessKey.Key length. Expected=%d Got=%d",
			accessKeyLength, len(accessKey.Key))
	}
}

func TestAccessKeyValid(t *testing.T) {
	accessKey := AccessKey{
		Key:     createNewKey(),
		ValidTo: time.Now().UTC().AddDate(0, 0, 1),
	}
	if !accessKey.Valid() {
		t.Errorf("accessKey.Valid wrong. Key should be valid: %v", accessKey)
	}
	accessKey = AccessKey{
		Key:     createNewKey(),
		ValidTo: time.Now().UTC().AddDate(0, 0, -1),
	}
	if accessKey.Valid() {
		t.Errorf("accessKey.Valid wrong. Key should not be valid: %v", accessKey)
	}
}
