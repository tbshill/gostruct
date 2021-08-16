package main

import (
	"bufio"
	"github.com/tbshill/csv"
	"flag"
	"fmt"
	"os"
	"strings"
	"io"
)

var delimeter string
var recordName string
var pkgName string
var outfile string

func main() {
	flag.StringVar(&outfile, "out", "", "-out mycsvrecord.go")
	flag.StringVar(&pkgName, "pkg", ",", "-pkg mypackage")
	flag.StringVar(&delimeter, "dl", ",", "delimeter -dl '|'")
	flag.StringVar(&recordName, "name", "", "-name 'MyCSVRecortd'")
	flag.Parse()

	var out io.Writer
	if outfile != "" {
		var err error
		var f *os.File
		f, err = os.OpenFile(outfile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0655)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating a file:%v", err)
			os.Exit(1)
		}
		defer f.Close()
		out = f
	} else {
		out = os.Stdout
	}


	if recordName == "" {
		fmt.Fprintf(os.Stderr, "Please specify a --name")
		os.Exit(1)
	}

	if len(flag.Args()) == 0 {
		fmt.Fprintf(os.Stderr, "Please specify at least one file to process")
		os.Exit(1)
	}

	if pkgName != "" {
		fmt.Fprintf(out,"package %s\n\n", pkgName)
	}

	for _, filename := range flag.Args() {
		firstLine, err := GetFirstLineOfFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}

		genText := GenerateText(firstLine, delimeter, recordName)
		fmt.Fprintln(out,genText)
	}

}

func GetFirstLineOfFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", fmt.Errorf("Could not open file %s:%v\n", filename, err)
	}

	scanner := bufio.NewScanner(f)
	scanner.Scan()
	firstLine := scanner.Text()
	f.Close()
	return firstLine, nil
}

func GenerateText(header string, del string, name string) string {
	columns := csv.RowToCols(header, del)

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("type %s struct {\n", name))

	for _, column := range columns {
		if column == "" {
			column = "_"
		}
		sb.WriteString(fmt.Sprintf("\t%s string `csv:\"%s\"`\n", toPublicName(column), column))
	}

	sb.WriteString("}\n")

	return sb.String()
}

func toPublicName(s string) string {
	titled := strings.Title(s)
	var sb strings.Builder
	allowedCharacters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	digits := "0123456789"

	for _, c := range titled {
		if c == '(' {
			break
		}

		if strings.Contains(allowedCharacters, string(c)) {
			sb.WriteRune(c)
		}
	}

	filtered := sb.String()
	if strings.Contains(digits, string(filtered[0])) {
		filtered = filtered[1:] + filtered[:1]
	}

	return filtered
}
