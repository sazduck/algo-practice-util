package main

import (
	"bufio"
	"io"
	"os"
	"strconv"
)

func main() {
	Run(os.Stdin, os.Stdout)
}

func Run(r io.Reader, w io.Writer) error {
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanWords)
	b := bufio.NewWriterSize(w, 1<<20)
	defer b.Flush()

	next := func() int {
		s.Scan()
		val, _ := strconv.Atoi(s.Text())
		return val
	}

	val := next()

	b.WriteString(strconv.Itoa(val))
	b.WriteByte('\n')
	return nil
}
