package print

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/ghodss/yaml"
)

// Printfln formats according to a format specifier and writes to standard output and append a new line at the end.
// It returns the number of bytes written and any write error encountered.
func Printfln(format string, a ...interface{}) (n int, err error) {
	return Fprintfln(os.Stdout, format, a...)
}

// Fprintfln formats according to a format specifier and writes to w and append a new line at the end.
// It returns the number of bytes written and any write error encountered.
func Fprintfln(w io.Writer, format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(w, format+"\n", a...)
}

// JSON formats according to a format specifier and writes to standard output and formats it to pretty json.
// It returns the number of bytes written and any write error encountered.
func JSON(a interface{}) (n int, err error) {
	json, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		return 0, err
	}
	return fmt.Println(string(json))
}

// Yaml formats according to a format specifier and writes to standard output and formats it to pretty json.
// It returns the number of bytes written and any write error encountered.
func Yaml(a interface{}) (n int, err error) {
	return Fyaml(os.Stdout, a)
}

// Fyaml formats according to a format specifier and writes to w and formats it to yaml.
// It returns the number of bytes written and any write error encountered.
func Fyaml(w io.Writer, a interface{}) (n int, err error) {
	yaml, err := yaml.Marshal(&a)
	if err != nil {
		return 0, err
	}
	return fmt.Fprintln(w, string(yaml))
}
