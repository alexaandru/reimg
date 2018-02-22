package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var (
	testRoot = "test"
	testA    = filepath.Join(testRoot, "a")
	testB    = filepath.Join(testRoot, "b")
	testC    = filepath.Join(testRoot, "c")
)

func TestAutoDetectFile(t *testing.T) {
	testCases := []struct {
		path, expPath string
		expErr        error
	}{
		{testA, "good.txt", nil},
		{testB, "bad.txt", nil},
		{testC, "", errNoMatch},
	}

	for _, tc := range testCases {
		if p, err := runIn(t, tc.path, autoDetectFile); tc.expErr != nil && err == nil ||
			tc.expErr == nil && err != nil ||
			tc.expErr != nil && err != nil && tc.expErr != err {
			t.Fatal("Expected", tc.expErr, "got", err, "for", tc.path)
		} else if tc.expPath != p {
			t.Fatal("Expected", tc.expPath, "got", p)
		}
	}
}

func TestAutoDetectTitle(t *testing.T) {
	testCases := []struct {
		path, expTitle string
		expErr         error
	}{
		{"bo/gu/sss", "", fmt.Errorf("open bo/gu/sss: no such file or directory")},
		{filepath.Join(testA, "good.txt"), "Hello", nil},
		{filepath.Join(testB, "bad.txt"), "", &errNoTitle{filepath.Join(testB, "bad.txt")}},
		{filepath.Join(testC, "bad.json"), "", &errNoTitle{filepath.Join(testC, "bad.json")}},
	}

	for _, tc := range testCases {
		if title, err := autoDetectTitle(tc.path); tc.expErr != nil && err == nil ||
			tc.expErr == nil && err != nil ||
			tc.expErr != nil && err != nil && tc.expErr.Error() != err.Error() {
			t.Fatal("Expected", tc.expErr, "got", err, "for", tc.path)
		} else if tc.expTitle != title {
			t.Fatal("Expected", tc.expTitle, "got", title)
		}
	}
}

func TestNormalizeTitle(t *testing.T) {
	t.Skip("Not implemented")
}

func TestExtractNum(t *testing.T) {
	t.Skip("Not implemented")
}

func TestPatchAll(t *testing.T) {
	t.Skip("Not implemented")
}

func TestIntegrationMain(t *testing.T) {
	t.Skip("Not implemented")
}

// helpers

func runIn(t *testing.T, folder string, fn func() (string, error)) (out string, errOut error) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	if err = os.Chdir(folder); err != nil {
		t.Fatal(err)
	}

	out, errOut = fn()

	if err = os.Chdir(cwd); err != nil {
		t.Fatal(err)
	}

	return
}
