package readfile

import (
	"fmt"
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
