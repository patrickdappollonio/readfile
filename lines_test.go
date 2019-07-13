package readfile

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

const contents = `
Example
Demo
Abc
Def
# ignore this`

func BenchmarkParser(b *testing.B) {
	var r1 = strings.NewReader(strings.Repeat(contents, 500))
	var results []string
	var err error

	b.ReportAllocs()

	for n := 0; n < b.N; n++ {
		results, err = reader(r1)
		if err != nil {
			b.Fatal(err.Error())
		}
	}

	for p, v := range results {
		fmt.Fprint(ioutil.Discard, p, v)
	}
}

func TestLineParser(t *testing.T) {
	cases := []struct {
		value    string
		expected string
	}{
		{"abc", "abc"},
		{"def", "def"},
		{"# abc", ""},
		{"     ", ""},
		{"###", ""},
	}

	for i := 0; i < len(cases); i++ {
		t.Run(fmt.Sprintf("case-%d", i+1), func(tt *testing.T) {
			returned := handleLine(cases[i].value)

			if returned != cases[i].expected {
				t.Fatalf("on case %d, expecting %q, but got %q", i+1, cases[i].expected, returned)
			}
		})
	}
}

func TestParsingFile(t *testing.T) {
	cases := []struct {
		Sent     string
		Expected []string
	}{
		{
			Sent:     "Example\nSecond\nThird",
			Expected: []string{"Example", "Second", "Third"},
		},
		{
			Sent:     "Example\nSecond\nThird\n\nFourth",
			Expected: []string{"Example", "Second", "Third", "Fourth"},
		},
		{
			Sent:     "First\n# ignored",
			Expected: []string{"First"},
		},
		{
			Sent:     "    First\n      # ignored",
			Expected: []string{"First"},
		},
	}

	for pos, v := range cases {
		t.Run(fmt.Sprintf("testing-case-%d", pos+1), func(tt *testing.T) {

			rd := strings.NewReader(v.Sent)

			recv, err := reader(rd)
			if err != nil {
				tt.Fatalf("not expecting an error on case %d, got: %s", pos+1, err.Error())
			}

			if len(recv) != len(v.Expected) {
				tt.Fatalf("not expecting different sizes (%d vs %d), got %#v vs %#v", len(recv), len(v.Expected), recv, v.Expected)
			}

			for i := 0; i < len(recv); i++ {
				if e, r := recv[i], v.Expected[i]; e != r {
					tt.Fatalf("found mismatch: expecting %q, got %q on position %d for test %d", e, r, i+1, pos+1)
				}
			}
		})
	}
}

type failedReader struct{}

func (r *failedReader) Read(p []byte) (int, error) {
	return 0, errors.New("failed")
}

func TestFailReader(t *testing.T) {

	var rd io.Reader = new(failedReader)

	_, err := reader(rd)
	if err == nil {
		t.Fatalf("expecting an error on call, but got nothing")
	}
}
