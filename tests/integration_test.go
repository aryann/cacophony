package tests

import (
	"cacophony/evaluator"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func validateFileName(name string, t *testing.T) {
	t.Helper()
	if !strings.HasSuffix(name, ".in") &&
		!strings.HasSuffix(name, ".want") {
		t.Fatalf("found unexpected file: %s", name)
	}
}

type testCase struct {
	in   string
	want string
}

func TestPrograms(t *testing.T) {
	testCases := make([]testCase, 0)

	files, err := ioutil.ReadDir("testdata")
	if err != nil {
		t.Fatalf("could not read testdata dir: %v", err)
	}
	for _, file := range files {
		validateFileName(file.Name(), t)
		if strings.HasSuffix(file.Name(), ".in") {
			testCases = append(testCases, testCase{
				in:   file.Name(),
				want: strings.TrimSuffix(file.Name(), ".in") + ".want",
			})
		}
	}

	for _, testCase := range testCases {
		t.Run(testCase.in, func(t *testing.T) {
			in, err := os.Open(filepath.Join("testdata", testCase.in))
			if err != nil {
				t.Fatalf("could not read input file: %v", err)
			}

			bytes, err := ioutil.ReadFile(filepath.Join("testdata", testCase.want))
			if err != nil {
				t.Fatalf("could not read output file: %v", err)
			}
			want := string(bytes)

			var got strings.Builder
			if _, err := evaluator.Evaluate(in, &got); err != nil {
				t.Fatalf("evaluation failed: %v", err)
			}

			if got.String() != want {
				t.Fatalf("want output <%s>, got <%s>", want, got.String())
			}
		})
	}
}
