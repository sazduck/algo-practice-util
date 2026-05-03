package testutil

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

type TestCase struct {
	Name     string `json:"name"`
	Input    string `json:"input"`
	Expected string `json:"expected"`
}

func Run(t *testing.T, tc []TestCase, f func(r io.Reader, w io.Writer) error) {
	for _, tt := range tc {
		t.Run(tt.Name, func(t *testing.T) {
			r := strings.NewReader(tt.Input)

			var sb strings.Builder

			if err := f(r, &sb); err != nil {
				t.Fatal(err)
			}

			got := strings.TrimSpace(sb.String())
			want := strings.TrimSpace(tt.Expected)

			if got != want {
				t.Errorf("\ninput: %s\nwant: %s\ngot: %s\n", tt.Input, want, got)
			}
		})
	}
}

func LoadTestCases(t *testing.T, path string) []TestCase {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to open file %s: %v", path, err)
	}
	var tests []TestCase
	err = json.Unmarshal(data, &tests)
	if err != nil {
		t.Fatalf("failed to parse JSON %v", err)
	}
	return tests
}

func RunWithDefaultTestCasesPath(t *testing.T, f func(r io.Reader, w io.Writer) error) {
	tests := LoadTestCases(t, "./test_cases.json")
	Run(t, tests, f)
}
