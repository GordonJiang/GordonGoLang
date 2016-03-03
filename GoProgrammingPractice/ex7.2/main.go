package main

import (
	"bytes"
	"fmt"
	"io"
)

type countWriter struct {
	count int64
	w     io.Writer
}

func main() {
	var buf bytes.Buffer
	w, pcount := CountingWriter(&buf)

	fmt.Fprintf(w, "This is a test.")
	fmt.Println(buf.String(), *pcount)
	fmt.Fprintf(w, "Add more contents")
	fmt.Println(buf.String(), *pcount)
}

func (c *countWriter) Write(p []byte) (n int, err error) {
	n, err = c.w.Write(p)
	c.count += int64(n)
	return
}

func CountingWriter(w io.Writer) (io.Writer, *int64) {
	c := &countWriter{
		w: w,
	}
	return c, &c.count
}
