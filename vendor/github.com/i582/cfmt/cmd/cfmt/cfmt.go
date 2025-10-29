package cfmt

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/i582/cfmt/internal"
)

var disableColors int32

// DisableColors globally disables colors.
// After parsing, no styles will be applied and the output will be clean text.
func DisableColors() {
	atomic.AddInt32(&disableColors, 1)
}

// EnableColors globally enables colors.
func EnableColors() {
	atomic.AddInt32(&disableColors, -1)
}

// RegisterStyle registers a new custom style.
//
// It will look like this:
//
//   cfmt.RegisterStyle("code", func(s string) string {
//	   return cfmt.Sprintf("{{%s}}::red|underline", s)
//   })
//
// The first argument is the name by which this style will be used,
// and the second is the styling function itself.
func RegisterStyle(name string, fn func(string) string) {
	internal.CustomMap[name] = fn
}

// Sprint is the same as fmt.Sprint.
func Sprint(a ...interface{}) string {
	return internal.ParseAndApply(fmt.Sprint(a...), disableColors == 1)
}

// Fprint is the same as fmt.Fprint.
func Fprint(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprint(w, Sprint(a...))
}

// Print is the same as fmt.Print.
func Print(a ...interface{}) (n int, err error) {
	return Fprint(os.Stdout, a...)
}

// Sprintln is the same as fmt.Sprintln.
func Sprintln(a ...interface{}) string {
	return Sprint(a...) + "\n"
}

// Fprintln is the same as fmt.Fprintln.
func Fprintln(w io.Writer, a ...interface{}) (n int, err error) {
	return fmt.Fprint(w, Sprintln(a...))
}

// Println is the same as fmt.Println.
func Println(a ...interface{}) (n int, err error) {
	return Fprintln(os.Stdout, a...)
}

// Sprintf is the same as fmt.Sprintf.
func Sprintf(format string, a ...interface{}) string {
	return Sprint(fmt.Sprintf(format, a...))
}

// Fprintf is the same as fmt.Fprintf.
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	return fmt.Fprint(w, Sprintf(format, a...))
}

// Printf is the same as fmt.Printf.
func Printf(format string, a ...interface{}) (n int, err error) {
	return Fprintf(os.Stdout, format, a...)
}

// Errorf is the same as fmt.Errorf.
func Errorf(format string, a ...interface{}) error {
	return fmt.Errorf(Sprintf(format, a...))
}
