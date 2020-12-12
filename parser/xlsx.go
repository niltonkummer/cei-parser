package parser

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/tealeg/xlsx"
	"reflect"
	"strings"
)

type xlsxDecoder struct {
	path string
}

func (x *xlsxDecoder) Decode() ([]*Stock, error) {

	f, err := xlsx.OpenFile(x.path)
	if err != nil {
		return nil, err
	}

	for _, sheet := range f.Sheets {
		return parseSheet(sheet), nil
	}
	return nil, errors.New("empty spread sheet")
}

func parseSheet(sheet *xlsx.Sheet) (assets []*Stock) {

	for line, row := range sheet.Rows {
		for col, cell := range row.Cells {
			switch {
			case strings.Contains(strings.ToLower(cleanValue(cell.Value)), "resumo dos saldos"):
				fmt.Println("Ativos em custódia:")
				headers := []string{}
				lineHeaders := line + 2
				for _, c := range sheet.Rows[lineHeaders].Cells[col:] {
					if c.String() != "" {
						headers = append(headers, c.String())
					}
				}
				lineActives := lineHeaders + 1

			stopAssets:
				for _, row := range sheet.Rows[lineActives:] {
					var currentPos int

					stock := &Stock{}
					valueOf := reflect.ValueOf(stock).Elem()
					for _, cell := range row.Cells {
						// Point to stop loop
						if strings.ToLower(cleanValue(cell.String())) == "valorização em reais" {
							break stopAssets
						}

						valueCleaned := cleanValue(cell.String())
						if valueCleaned != "" && valueCleaned != "#" && currentPos < valueOf.NumField() {

							if valueOf.Field(currentPos).Type().Kind() == reflect.Int {
								v, err := cell.Int64()
								if err == nil {
									valueOf.Field(currentPos).SetInt(v)
								}
							} else if valueOf.Field(currentPos).Type().Name() == "Decimal" {
								v, _ := decimal.NewFromString(valueCleaned)
								valueOf.Field(currentPos).Set(reflect.ValueOf(v))
							} else {
								valueOf.Field(currentPos).SetString(valueCleaned)
							}

							currentPos++
						}
					}
					assets = append(assets, stock)
				}
			}
		}
	}
	return
}
