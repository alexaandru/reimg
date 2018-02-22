package main

import "fmt"

type errNoTitle struct {
	path string
}

var errNoMatch = fmt.Errorf("Cannot find a unique match for *.txt")

func (e *errNoTitle) Error() string {
	return "Unable to detect title in " + e.path
}
