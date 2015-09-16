package test_helpers

import (
	"io/ioutil"
	"os"
	"testing"
)

func CreateTempfile(t *testing.T, data string) string {
	dir := os.TempDir()

	f, err := ioutil.TempFile(dir, "test")
	if err != nil {
		t.Fatalf("Tempfile %s: %v", f.Name(), err)
	}

	if err := ioutil.WriteFile(f.Name(), []byte(data), 0644); err != nil {
		t.Fatalf("WriteFile %s: %v", f.Name(), err)
	}

	return f.Name()
}
