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
	"strconv"
	"strings"

	humanize "github.com/dustin/go-humanize"
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
	return runLines(lines, w)
}

func runLines(lines []string, w io.Writer) error {
	n, err := clearance(lines)
	if err != nil {
		return fmt.Errorf("couldn't calculate clearance for line: %q", err)
	}

	for i, line := range lines {
		// Skip the newlines at the end of the file.
		if i == len(lines)-1 {
			continue
		}
		if !strings.HasPrefix(line, " ") {
			fmt.Fprintf(w, "%v\n", line)
			continue
		}
		err := writePostingLine(w, line, n)
		if err != nil {
			return fmt.Errorf("error on line %v: %v", i, err)
		}
	}
	return nil
}

func clearance(lines []string) (int, error) {
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
		amount, err := addCommas(tokens[2])
		if err != nil {
			return 0, fmt.Errorf("couldn't add comma to line %q: %v", line, err)
		}
		if i := 2 + len(acct) + 2 + dotPosition(amount); i > max {
			max = i
		}
	}
	return max, nil
}

var startNumbers = regexp.MustCompile("^([-0-9]*)(.*)$")

func addCommas(amount string) (string, error) {
	// grab the digits off the front, including commas and minuses (not dots, everything after the decimal we leave)
	digitsParts := startNumbers.FindStringSubmatch(amount)
	numbersAtStart := digitsParts[1]
	endBit := digitsParts[2]
	if numbersAtStart == "" {
		return amount, nil
	}

	// parse it into an integer
	numbers, err := strconv.ParseInt(numbersAtStart, 10, 64)
	if err != nil {
		return "", fmt.Errorf("couldn't parse number %q from amount %q: %v", numbers, amount, err)
	}

	// format it with commas
	return humanize.Comma(numbers) + endBit, nil
}

func writePostingLine(w io.Writer, line string, n int) error {
	tokens := whitespace.Split(line, -1)

	// handle account
	acct := tokens[1]
	fmt.Fprintf(w, "  %v", acct)
	if len(tokens) == 2 {
		fmt.Fprint(w, "\n")
		return nil
	}

	// handle amount
	amount, err := addCommas(tokens[2])
	if err != nil {
		return fmt.Errorf("couldn't add comma to line %q: %v", line, err)
	}
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
	return nil
}

func dotPosition(s string) int {
	if idx := strings.Index(s, "."); idx != -1 {
		return idx
	}
	if idx := strings.LastIndexAny(s, "0123456789,"); idx != -1 {
		return idx + 1
	}
	return -1
}

func fPrintSpaces(w io.Writer, n int) {
	for i := 0; i < n; i++ {
		fmt.Fprint(w, " ")
	}
}
