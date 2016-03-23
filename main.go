package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Debug, erase when done
var _ = fmt.Printf

var ch byte
var ch1 byte
var buf *bufio.Reader
var outbuf *bufio.Writer

var override bool

func getc() {
	var err error
	switch ch1 {
	case 0:
		ch, err = buf.ReadByte()
		if err != nil {
			ch = 0
			return
		}
		break

	default:
		ch = ch1
		ch1 = 0
	}
}

func putc(b byte) {
	ch1 = b
}

func gen(inf, of *os.File) {
	buf = bufio.NewReader(inf)
	outbuf = bufio.NewWriter(of)

	fmt.Fprintf(outbuf, "<html xmlns=\"http://www.w3.org/1999/xhtml\"><head><meta http-equiv=\"Content-Type\" content=\"text/html; charset=\"><link rel=\"stylesheet\" href=\"style.css\" type=\"text/css\"></head><body>\n")

	incode := false

done:
	for {
		getc()
		switch ch {
		case 0:
			break done
		case '(', '/':
			chold := ch
			getc()
			if ch == '*' {
				if incode {
					fmt.Fprintf(outbuf, "</code>\n")
					incode = false
				}
				fmt.Fprintf(outbuf, "<table width=\"100%%\"><tr class=\"doc\"><td>\n")
				break
			}
			putc(ch)
			fmt.Fprintf(outbuf, "%c", chold)
		case '*':
			getc()
			if ch == ')' || ch == '/' {
				fmt.Fprintf(outbuf, "</td></tr></table>\n")
				fmt.Fprintf(outbuf, "<code>\n")
				incode = true
				break
			}
			putc(ch)
			fmt.Fprintf(outbuf, "*")
		case '\n':
			if !incode {
				fmt.Fprintf(outbuf, "</br>\n")
			} else {
				fmt.Fprintf(outbuf, "\n")
			}
		default:
			fmt.Fprintf(outbuf, "%c", ch)
		}
	}

	fmt.Fprintf(outbuf, "</body></html>\n")
	outbuf.Flush()
}

func main() {
	flag.BoolVar(&override, "override", false, "override output file")

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		flag.Usage()
		return
	}

	ifname := args[0]
	inf, err := os.Open(ifname)
	if err != nil {
		fmt.Printf("Error reading file %s!\n", args[0])
		return
	}

	ofname := fmt.Sprintf("%s.html", strings.Split(ifname, ".")[0])
	if !override {
		if _, err := os.Stat(ofname); err == nil {
			fmt.Printf("Output file %s already exists.\n", ofname)
			return
		}
	}

	of, err := os.Create(ofname)
	if err != nil {
		fmt.Printf("Error creating file %s!\n", ofname)
		return
	}

	gen(inf, of)
}
