package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Exit codes.
const (
	_ = iota
	Fail
)

var substitutes = map[string]string{
	"Inspectoratul pentru Situații de Urgență":           "ISU",
	"Inspectoratul pentru Situatii de Urgenta":           "ISU",
	"Inspectoratului pentru Situații de Urgență":         "ISU",
	"Inspectoratului pentru Situatii de Urgenta":         "ISU",
	"Inspectoratul General pentru Situatii de Urgenta":   "IGSU",
	"Inspectoratul General pentru Situații de Urgență":   "IGSU",
	"Inspectoratului General pentru Situații de Urgență": "IGSU",
	"Inspectoratului General pentru Situatii de Urgenta": "IGSU",
	"Departamentul pentru Situații de Urgență":           "DSU",
	"Departamentul pentru Situatii de Urgenta":           "DSU",
	"Departamentului pentru Situații de Urgență":         "DSU",
	"Departamentului pentru Situatii de Urgenta":         "DSU",
	"Ministerul Afacerilor Interne":                      "MAI",
	"Ministerului Afacerilor Interne":                    "MAI",
	"Centrul de Formare Iniţială şi Continuă":            "CFIC",
	"Centrul de Formare Initiala si Continua":            "CFIC",
	"Á":          "a",
	"Ă":          "A",
	"Â":          "A",
	"Î":          "I",
	"Ș":          "S",
	"Ş":          "S",
	"Ț":          "T",
	"Ţ":          "T",
	"á":          "a",
	"ă":          "a",
	"â":          "a",
	"î":          "i",
	"ș":          "s",
	"ş":          "s",
	"ț":          "t",
	"ţ":          "t",
	"S.M.U.R.D.": "SMURD",
	"I.G.S.U.":   "IGSU",
	"I.S.U.":     "ISU",
	"A.S.F.R.":   "ASFR",
	"M.A.I.":     "MAI",
	"C.B.R.N.":   "CBRN",
	"D.S.U.":     "DSU",
}

var reNonASCII, reTrim, reNum = regexp.MustCompile(`[^a-zA-Z0-9]`), regexp.MustCompile(`^_+|_+$`),
	regexp.MustCompile(`_([\d])+\.`)

func autoDetectFile() (string, error) {
	gx, err := filepath.Glob("*.txt")
	if err != nil {
		return "", err
	}

	if len(gx) != 1 {
		return "", errNoMatch
	}

	return gx[0], nil
}

func autoDetectTitle(fname string) (string, error) {
	f, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner, candidate := bufio.NewScanner(f), ""
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" || line == "" {
			continue
		} else if strings.HasPrefix(line, "-----") {
			return normalizeTitle(candidate), nil
		} else {
			candidate = line
		}
	}

	return "", &errNoTitle{fname}
}

func normalizeTitle(title string) string {
	for k, v := range substitutes {
		title = strings.Replace(title, k, v, -1)
	}

	title = reNonASCII.ReplaceAllString(title, "_")
	return reTrim.ReplaceAllString(title, "")
}

func extractNum(image string) (int, error) {
	mx := reNum.FindAllStringSubmatch(image, -1)
	if len(mx) == 0 || len(mx[0]) < 2 {
		return 0, fmt.Errorf("Unable to extract number from %s", image)
	}

	return strconv.Atoi(mx[0][1])
}

func patchAll(title string) error {
	images, err := filepath.Glob("*.jpg")
	if err != nil {
		return err
	}

	nextImg := 1
	rest := []string{}
	for _, image := range images {
		if strings.HasPrefix(image, title) {
			n, err := extractNum(image)
			if err != nil {
				return err
			}
			if n > nextImg {
				nextImg = n + 1
			}
		} else {
			rest = append(rest, image)
		}
	}

	sort.Strings(rest)
	for _, image := range rest {
		newImage := fmt.Sprintf("%s_%d.jpg", title, nextImg)
		nextImg++
		err := os.Rename(image, newImage)
		if err != nil {
			return err
		}
	}

	all, _ := filepath.Glob("*.txt")
	article := all[0]

	return os.Rename(article, title+".txt")
}

func main() {
	fname, err := autoDetectFile()
	if err != nil {
		fmt.Println(err)
		os.Exit(Fail)
	}

	title, err := autoDetectTitle(fname)
	if err != nil {
		fmt.Println(err)
		os.Exit(Fail)
	}

	if err := patchAll(title); err != nil {
		fmt.Println(err)
		os.Exit(Fail)
	}
}
