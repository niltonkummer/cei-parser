package parser

import (
	"errors"
	"github.com/shopspring/decimal"
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
	Div    decimal.Decimal
	JCP    decimal.Decimal
}

type Parser interface {
	Decode() ([]*Stock, error)
}

func cleanValue(v string) string {
	return strings.Trim(v, " \t\n")
}

func DetectExtension(filename string) (Parser, error) {
	ext := filepath.Ext(filename)
	switch ext {
	case ".xls":
		return &xlsDecoder{path: filename}, nil
	case ".xlsx":
		return &xlsxDecoder{path: filename}, nil
	}
	return nil, errors.New("invalid format")
}
