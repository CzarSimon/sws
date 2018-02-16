package service

import (
	"testing"
)

func TestReadService(t *testing.T) {
	manifest, err := ReadService("../../resources/test-service.yml")
	if err != nil {
		t.Fatalf("ReadService returned error")
	}
	if manifest.ApiVersion != "v1" {
		t.Errorf("manifest.ApiVersion wrong. Expected=v1 Got=%s", manifest.ApiVersion)
	}
	s := manifest.Spec
	expectedImage := "czarsimon/sws/test-image:latest"
	if s.Image != expectedImage {
		t.Errorf("service.Image wrong. Expected=%s Got=%s", expectedImage, s.Image)
	}
}
