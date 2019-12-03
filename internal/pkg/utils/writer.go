package utils

import "io"

// LineWriter is an io.Writer that appends a newline character to the argument
// to each call to Write.
type LineWriter struct {
	io.Writer
}

// Write writes the given bytes to the underlying writer and appends a newline
// character to the end.
func (sc *LineWriter) Write(p []byte) (int, error) {
	_, err := sc.writeLine(p)
	return len(p), err
}

func (sc *LineWriter) writeLine(p []byte) (int, error) {
	return sc.Writer.Write(append(p, []byte("\n")...))
}
