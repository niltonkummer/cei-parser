package main

import (
	"github.com/shopspring/decimal"
	"os"
	"path/filepath"
	"strings"
)

type Stock struct {
	Name   string
	Spec   string
	Code   string
	Total  int
	Price  decimal.Decimal
	Amount decimal.Decimal
}

type Parser interface {
	Decode() error
}

func cleanValue(v string) string {
	return strings.Trim(v, " \t\n")
}

func detectExtension(filename string) Parser {
	ext := filepath.Ext(filename)
	switch ext {
	case ".xls":
		return &xlsDecoder{path: filename}
	case ".xlsx":
		return &xlsxDecoder{path: filename}
	}
	return nil
}

func main() {

	parser := detectExtension(os.Args[1])
	parser.Decode()
}
