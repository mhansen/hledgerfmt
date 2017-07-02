package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

var whitespace = regexp.MustCompile("\\s+")

func runFile(inFile string, outFile string) {
	r, err := os.Open(inFile)
	if err != nil {
		log.Fatalf("Couldn't open input: %v", err)
	}

	var b bytes.Buffer
	if err = run(r, &b); err != nil {
		log.Fatalf("Error: %v", err)
	}

	w, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("Couldn't open for write: %v", err)
	}

	fmt.Fprint(w, b.String())
}

func main() {
	flag.Parse()
	args := flag.Args()
	for _, file := range args {
		runFile(file, file)
	}
	if len(args) == 0 {
		runFile("/dev/stdin", "/dev/stdout")
	}
}

func run(r io.Reader, w io.Writer) error {
	f, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	lines := strings.Split(string(f), "\n")
	runLines(lines, w)
	return nil
}

func runLines(lines []string, w io.Writer) {
	n := clearance(lines)

	for i, line := range lines {
		// Skip the newlines at the end of the file.
		if i == len(lines)-1 {
			continue
		}
		if !strings.HasPrefix(line, " ") {
			fmt.Fprintf(w, "%v\n", line)
			continue
		}
		writePostingLine(w, line, n)
	}
}

func clearance(lines []string) int {
	max := 0
	for _, line := range lines {
		if !strings.HasPrefix(line, " ") {
			continue
		}
		tokens := whitespace.Split(line, -1)
		if len(tokens) == 2 {
			continue
		}
		acct := tokens[1]
		amount := tokens[2]
		if i := 2 + len(acct) + 2 + dotPosition(amount); i > max {
			max = i
		}
	}
	return max
}

func writePostingLine(w io.Writer, line string, n int) {
	tokens := whitespace.Split(line, -1)

	// handle account
	acct := tokens[1]
	fmt.Fprintf(w, "  %v", acct)
	if len(tokens) == 2 {
		fmt.Fprint(w, "\n")
		return
	}

	// handle amount
	amount := tokens[2]
	fmt.Fprint(w, "  ")
	if idx := dotPosition(amount); idx != -1 {
		charsSoFar := 2 + len(acct) + 2 + idx
		fPrintSpaces(w, n-charsSoFar)
	}
	fmt.Fprint(w, amount)

	// handle everything else
	for i := 3; i < len(tokens); i++ {
		fmt.Fprintf(w, " %v", tokens[i])
	}
	fmt.Fprint(w, "\n")
}

func dotPosition(s string) int {
	if idx := strings.Index(s, "."); idx != -1 {
		return idx
	}
	if idx := strings.LastIndexAny(s, "0123456789"); idx != -1 {
		return idx + 1
	}
	return -1
}

func fPrintSpaces(w io.Writer, n int) {
	for i := 0; i < n; i++ {
		fmt.Fprint(w, " ")
	}
}
