package config

import (
	"os"
	"testing"
)

func TestLoadEmptyFile(t *testing.T) {
	f, _ := os.CreateTemp("", "empty*.json")
	f.Close()
	defer os.Remove(f.Name())
	_, err := Init(f.Name())
	if err == nil {
		t.Error("expected error for empty file")
	}
}

func TestLoadNegativeUint(t *testing.T) {
	content := `{"Floors": -1, "Monsters": 2, "OpenAt": "14:05:00", "Duration": 2}`
	f, _ := os.CreateTemp("", "neg*.json")
	f.WriteString(content)
	f.Close()
	defer os.Remove(f.Name())
	_, err := Init(f.Name())
	if err == nil {
		t.Error("expected error for negative Floors")
	}
}

func TestLoadMissingField(t *testing.T) {
	content := `{"Monsters": 2, "OpenAt": "14:05:00", "Duration": 2}`
	f, _ := os.CreateTemp("", "missing*.json")
	f.WriteString(content)
	f.Close()
	defer os.Remove(f.Name())
	_, err := Init(f.Name())
	if err == nil {
		t.Error("expected error for missing Floors")
	}
}
