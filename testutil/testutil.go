package testutil

import (
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
)

/*

// coderun
copy([...document.querySelectorAll('.io-sample>div')].map((d, i) => ({
  name: `example ${i}`,
  input: d.querySelector('.input-snippet code').innerText,
  expected: d.querySelector('.output-snippet code').innerText,
})));

// яндекс контест мб другие будут классы TODO
[...document.querySelectorAll('.StatementSamples_StatementSamples-Sample__0nXC_')]
.map((d, i) => ({
  name: `example ${i}`,
  input: d.querySelector('.CodeSnippet_CodeSnippet-Code__CBRAd.IOSnippet_IOSnippet-Code__En89R code').innerText,
  expected: d.querySelectorAll('.CodeSnippet_CodeSnippet-Code__CBRAd.IOSnippet_IOSnippet-Code__En89R code')[1].innerText
}))

*/

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
				t.Errorf("Ввод: %s\nОжидалось: %s\nПолучено: %s", tt.Input, want, got)
			}
		})
	}
}

func LoadTestCases(t *testing.T, path string) []TestCase {
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Не удалось прочитать файл %s: %v", path, err)
	}
	var tests []TestCase
	err = json.Unmarshal(data, &tests)
	if err != nil {
		t.Fatalf("Ошибка при десериализации JSON %v", err)
	}
	return tests
}

func RunWithDefaultTestCasesPath(t *testing.T, f func(r io.Reader, w io.Writer) error) {
	tests := LoadTestCases(t, "./test_cases.json")
	Run(t, tests, f)
}
