package logging

import (
	"bytes"
	"io"
)

type LinePrefixWriter struct {
	io.Writer
	Prefix string
}

func (pw LinePrefixWriter) Write(p []byte) (n int, err error) {
	var buf bytes.Buffer
	buf.WriteString(pw.Prefix)

	n = len(p)

	for {
		i := bytes.IndexRune(p, '\n')
		if i == -1 {
			break
		}
		buf.Write(p[:i+1])

		if i+1 >= len(p) {
			p = nil
			break
		}
		buf.WriteString(pw.Prefix)

		p = p[i+1:]
	}

	buf.Write(p)
	_, err = buf.WriteTo(pw.Writer)

	return
}
