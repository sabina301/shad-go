//go:build !solution

package externalsort

import (
	"container/heap"
	"io"
	"os"
	"strings"
)

type MyReader struct {
	reader io.Reader
}

func (mR MyReader) ReadLine() (string, error) {
	buff := make([]byte, 1)
	sb := strings.Builder{}
	for {
		n, err := mR.reader.Read(buff)
		if err != nil {
			if err == io.EOF {
				if n == 1 {
					sb.WriteByte(buff[0])
				}
				if sb.Len() > 0 {
					return sb.String(), nil
				}
				return "", io.EOF
			}
			return "", err
		}
		if n > 0 {
			if buff[0] == '\n' {
				break
			}
			sb.WriteByte(buff[0])
		}
	}
	return sb.String(), nil
}

type MyWriter struct {
	writer io.Writer
}

func (mW MyWriter) Write(s string) error {
	buff := []byte(s)
	for i := 0; i < len(buff); i++ {
		_, err := mW.writer.Write([]byte{buff[i]})
		if err != nil {
			return err
		}
	}
	if len(buff) == 0 || buff[len(buff)-1] != '\n' {
		_, _ = mW.writer.Write([]byte{'\n'})
	}
	return nil
}

func NewReader(r io.Reader) LineReader {
	return MyReader{reader: r}
}

func NewWriter(w io.Writer) LineWriter {
	return MyWriter{writer: w}
}

func Merge(w LineWriter, readers ...LineReader) error {
	var myHeap MyHeapWittReaders
	heap.Init(&myHeap)
	for _, r := range readers {
		firstLine, err := r.ReadLine()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF && firstLine == "" {
			continue
		}
		firstElement := ValAndFileReader{Val: firstLine, Reader: r}
		heap.Push(&myHeap, firstElement)
	}

	for myHeap.Len() != 0 {
		upElement := heap.Pop(&myHeap).(ValAndFileReader)
		if err := w.Write(upElement.Val); err != nil {
			return err
		}

		newVal, err := upElement.Reader.ReadLine()
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF && newVal == "" {
			continue
		}
		newElement := ValAndFileReader{Val: newVal, Reader: upElement.Reader}
		heap.Push(&myHeap, newElement)
	}
	return nil
}

func Sort(w io.Writer, in ...string) error {
	lineReaders := make([]LineReader, 0)
	for _, fileName := range in {
		file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		reader := NewReader(file)
		h := MyHeap{}
		heap.Init(&h)
		for {
			line, err := reader.ReadLine()
			if err != nil && err != io.EOF {
				file.Close()
				return err
			}
			if err == io.EOF {
				if line != "" {
					heap.Push(&h, line)
				}
				break
			}
			heap.Push(&h, line)
		}
		file.Close()

		file2, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		writer := NewWriter(file2)
		delStr := make([]string, 0)
		for h.Len() != 0 {
			delStr = append(delStr, heap.Pop(&h).(string))
		}
		for _, str := range delStr {
			if err := writer.Write(str); err != nil {
				file2.Close()
				return err
			}
		}
		file2.Close()

		file3, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
		if err != nil {
			return err
		}
		lineReaders = append(lineReaders, NewReader(file3))
	}
	defer func() {
		for _, r := range lineReaders {
			if closer, ok := r.(io.Closer); ok {
				_ = closer.Close()
			}
		}
	}()
	return Merge(NewWriter(w), lineReaders...)
}
