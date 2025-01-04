package logging

import (
	"bytes"
	"testing"
)

func TestLinePrefixWriter(t *testing.T) {
	t.Run("Write", func(t *testing.T) {
		var out bytes.Buffer
		writer := LinePrefixWriter{Writer: &out, Prefix: "@cee: "}
		writer.Write([]byte("Hello\n World!\n"))

		expected := "@cee: Hello\n@cee:  World!\n"
		if out.String() != expected {
			t.Errorf("Expected %q, got %q", expected, out.String())
		}
	})

	t.Run("With no new line", func(t *testing.T) {
		var out bytes.Buffer
		writer := LinePrefixWriter{Writer: &out, Prefix: "@cee: "}
		writer.Write([]byte("Hello World!"))

		expected := "@cee: Hello World!"
		if out.String() != expected {
			t.Errorf("Expected %q, got %q", expected, out.String())
		}
	})
}
