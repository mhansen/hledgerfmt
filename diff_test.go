package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"github.com/pmezard/go-difflib/difflib"
)

func TestDiffs(t *testing.T) {
	for _, c := range []string{"Example"} {
		inFileName := "testdata/" + c + ".in.hledger"
		inFile, err := os.Open(inFileName)
		if err != nil {
			t.Fatalf("Could not open test data %v: %v", inFileName, err)
		}
		wantFileName := "testdata/" + c + ".out.hledger"
		wantBytes, err := ioutil.ReadFile(wantFileName)
		if err != nil {
			t.Fatalf("Could not open test data %v: %v", wantFileName, err)
		}
		want := string(wantBytes)
		var b bytes.Buffer
		if err := run(inFile, &b); err != nil {
			t.Fatalf("%q: unexpected error in run: %v", c, err)
		}
		got := b.String()
		if got != want {
			diff := difflib.UnifiedDiff{
				A:        difflib.SplitLines(want),
				B:        difflib.SplitLines(got),
				FromFile: "Want",
				ToFile:   "Got",
				Context:  5,
			}
			text, err := difflib.GetUnifiedDiffString(diff)
			if err != nil {
				t.Fatalf("%q: unexpected error in GetUnifiedDiffString: %v", c, err)
			}
			t.Errorf("%q: Diffs found. len(got)=%v, len(%q)=%v:\n%v", c, len(got), wantFileName, len(want), text)
		}
	}
}
