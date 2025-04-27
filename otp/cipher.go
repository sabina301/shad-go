//go:build !solution

package otp

import (
	"io"
)

type MyReader struct {
	reader     io.Reader
	prngReader io.Reader
}

func NewReader(r io.Reader, prng io.Reader) io.Reader {
	return MyReader{r, prng}
}

func (mr MyReader) Read(p []byte) (int, error) {
	n, err := mr.reader.Read(p)
	if n == 0 {
		return n, err
	}
	prng := make([]byte, n)
	n, err = mr.prngReader.Read(prng)
	if n == 0 {
		return n, err
	}

	for i, _ := range p {
		p[i] = p[i] ^ prng[i%len(prng)]
	}

	return len(prng), nil
}

type MyWriter struct {
	writer io.Writer
	prng   io.Reader
}

func NewWriter(w io.Writer, prng io.Reader) io.Writer {
	return MyWriter{w, prng}
}

func (mw MyWriter) Write(p []byte) (int, error) {
	prng := make([]byte, len(p))
	n, err := mw.prng.Read(prng)
	if n == 0 {
		return n, err
	}
	p2 := make([]byte, len(p))
	copy(p2, p)
	for i, _ := range p2 {
		p2[i] = p2[i] ^ prng[i%len(prng)]
	}
	n, err = mw.writer.Write(p2)
	return n, err
}
