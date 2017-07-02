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

var f = flag.String("f", "", "file to autoformat")

func main() {
	flag.Parse()
	var err error
	var r io.Reader = os.Stdin
	if *f != "" {
		r, err = os.Open(*f)
		if err != nil {
			log.Fatalf("Couldn't open %q: %v", f, err)
		}
	}
	out, err := run(r)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	var w io.Writer = os.Stdout
	if *f != "" {
		w, err = os.Create(*f)
		if err != nil {
			log.Fatalf("Couldn't open %q for write: %v", f, err)
		}
	}
	fmt.Fprint(w, out)
}

func run(r io.Reader) (string, error) {
	f, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(f), "\n")
	return runLines(lines), nil
}

func runLines(lines []string) string {
	n := clearance(lines)

	var b bytes.Buffer
	for i, line := range lines {
		// Skip the newlines at the end of the file.
		if i == len(lines)-1 {
			continue
		}
		if !strings.HasPrefix(line, " ") {
			fmt.Fprintf(&b, "%v\n", line)
			continue
		}
		writePostingLine(&b, line, n)
	}
	return b.String()
}

func clearance(lines []string) int {
	constraintLine := ""
	max := 0
	for _, line := range lines {
		if !strings.HasPrefix(line, " ") {
			continue
		}
		tokens := re.Split(line, -1)
		if len(tokens) == 2 {
			continue
		}
		acct := tokens[1]
		amount := tokens[2]
		if i := 2 + len(acct) + 2 + dotPosition(amount); i > max {
			max = i
			constraintLine = line
		}
	}
	log.Printf("Constraint Line: %q\n", constraintLine)
	return max
}

func writePostingLine(w io.Writer, line string, n int) {
	tokens := re.Split(line, -1)

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

var re = regexp.MustCompile("\\s+")
