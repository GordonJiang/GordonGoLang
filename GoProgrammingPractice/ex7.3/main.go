package main

import (
	"fmt"
	"io"
	"strings"
)

type limitReader struct {
	n      int64
	reader io.Reader
}

func main() {
	reader := strings.NewReader("This is a test.This is a test.This is a test.This is a test.This is a test.This is a test")
	r := NewLimitReader(reader, 20)
	data := make([]byte, 4)
	for {
		n, err := r.Read(data)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Printf("Read %d bytes from r:%s \n", n, string(data))
	}
}

func (l *limitReader) Read(p []byte) (n int, err error) {
	n, err = l.reader.Read(p)
	if err != nil {
		return 0, err
	}
	l.n -= int64(n)
	if l.n < 0 {
		return 0, io.EOF
	}
	return
}

func NewLimitReader(r io.Reader, n int64) io.Reader {
	return &limitReader{
		n:      n,
		reader: r,
	}
}
