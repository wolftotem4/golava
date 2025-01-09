package env

import (
	"testing"
)

func TestSetEnvVar(t *testing.T) {
	t.Run("Test SetEnvVar", func(t *testing.T) {
		content := []byte("APP_KEY=base64:oldkey\n")
		key := "APP_KEY"
		value := "base64:newkey"

		newcontent, appends := SetEnvVar(content, key, value, "base64:oldkey")
		if appends {
			t.Errorf("Expected false, got %v", appends)
		}
		if string(newcontent) != "APP_KEY=base64:newkey\n" {
			t.Errorf("Expected APP_KEY=base64:newkey, got %s", string(newcontent))
		}
	})

	t.Run("Test SetEnvVar", func(t *testing.T) {
		content := []byte("APP_NAME=Golava\nAPP_KEY=base64:oldkey\n")
		key := "APP_KEY"
		value := "base64:newkey"

		newcontent, appends := SetEnvVar(content, key, value, "base64:oldkey")
		if appends {
			t.Errorf("Expected false, got %v", appends)
		}
		if string(newcontent) != "APP_NAME=Golava\nAPP_KEY=base64:newkey\n" {
			t.Errorf("Expected APP_NAME=Golava\nAPP_KEY=base64:newkey, got %s", string(newcontent))
		}
	})

	t.Run("Test SetEnvVar", func(t *testing.T) {
		content := []byte("APP_NAME=Golava\n")
		key := "APP_KEY"
		value := "base64:newkey"

		newcontent, appends := SetEnvVar(content, key, value, "")
		if !appends {
			t.Errorf("Expected true, got %v", appends)
		}
		if string(newcontent) != "\nAPP_KEY=base64:newkey\n" {
			t.Errorf("Expected \nAPP_KEY=base64:newkey\n, got %s", string(newcontent))
		}
	})

	t.Run("Test SetEnvVar", func(t *testing.T) {
		content := []byte("APP_NAME=Golava\nAPP_KEY=\nLISTEN_ADDR=127.0.0.1:8080\n")
		key := "APP_KEY"
		value := "base64:newkey"

		newcontent, appends := SetEnvVar(content, key, value, "")
		if appends {
			t.Errorf("Expected false, got %v", appends)
		}
		if string(newcontent) != "APP_NAME=Golava\nAPP_KEY=base64:newkey\nLISTEN_ADDR=127.0.0.1:8080\n" {
			t.Errorf("Expected APP_NAME=Golava\nAPP_KEY=base64:newkey\nLISTEN_ADDR=127.0.0.1:8080\n, got %s", string(newcontent))
		}
	})
}
