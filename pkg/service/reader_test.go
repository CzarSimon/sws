package service

import (
	"testing"
)

func TestReadServiceYml(t *testing.T) {
	manifest, err := ReadService("../../resources/test/service.yml")
	if err != nil {
		t.Fatalf("ReadService returned error: %s", err.Error())
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

func TestReadServiceYaml(t *testing.T) {
	manifest, err := ReadService("../../resources/test/service.yaml")
	if err != nil {
		t.Fatalf("ReadService returned error: %s", err.Error())
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

func TestReadServiceJson(t *testing.T) {
	manifest, err := ReadService("../../resources/test/service.json")
	if err != nil {
		t.Fatalf("ReadService returned error: %s", err.Error())
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

func TestReadServiceTxt(t *testing.T) {
	_, err := ReadService("../../resources/test/service.txt")
	if err == nil {
		t.Fatalf("ReadService should have returned error for filetype .txt")
	}
}
