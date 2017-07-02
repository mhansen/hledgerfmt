package main

import "testing"

func TestDotPosition(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"", -1},
		{".", 0},
		{"1.0", 1},
		{".0", 0},
		{"100.00", 3},
		{"1", 1},
		{"10", 2},
		{"100", 3},
		{"-1", 2},
		{"-10", 3},
		{"-100", 4},
		{"100AUD", 3},
		{"100.00USD", 3},
	}
	for _, c := range cases {
		if got := dotPosition(c.in); got != c.want {
			t.Errorf("dotPosition(%q)=%v, want=%v", c.in, got, c.want)
		}
	}
}
