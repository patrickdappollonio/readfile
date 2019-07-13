package readfile

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestParseFile(t *testing.T) {
	file, err := ioutil.TempFile("", "prefix")
	if err != nil {
		t.Fatalf("unable to create file to parse: %s", err.Error())
	}
	defer os.Remove(file.Name())

	if _, err := file.Write([]byte("first\nsecond\n       third  \n# fourth")); err != nil {
		t.Fatalf("unable to write to test file: %s", err.Error())
	}

	file.Close()

	expected := []string{"first", "second", "third"}

	rd := New(file.Name())

	returned, err := rd.Parse()
	if err != nil {
		t.Fatalf("not expecting an error when calling parse, got: %s", err.Error())
	}

	if a, b := len(expected), len(returned); a != b {
		t.Fatalf("expecting compared slices to be equal in size, got %d vs %d", a, b)
	}

	for i := 0; i < len(expected); i++ {
		if expected[i] != returned[i] {
			t.Fatalf("expecting values in position %d to be equal, got %q vs %q", i, expected[i], returned[i])
		}
	}
}
