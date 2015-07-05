package storage

import (
	"testing"
)

func TestCrud(t *testing.T) {
	s := New()
	s.Set("hi", "hello")

	if s.Get("hi").(string) != "hello" {
		t.Errorf("Failed to test set and get. Actual Data: %v", s.Get("/free").(string))
	}
}
