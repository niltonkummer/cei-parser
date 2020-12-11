package main

import (
	"encoding/json"
	"fmt"
	"github.com/extrame/xls"
	"github.com/shopspring/decimal"

	"os"
	"reflect"
	"strconv"
	"strings"
)

type xlsDecoder struct {
	path string
}

func (x *xlsDecoder) Decode() error {

	f, err := xls.Open(x.path, "utf-8")
	if err != nil {
		return err
	}

	parseXlsSheet(f.GetSheet(0))

	return nil
}

func parseXlsSheet(sheet *xls.WorkSheet) {

	months := [][]*Stock{}

	for line := uint16(0); line < sheet.MaxRow; line++ {
		row := sheet.Row(int(line))
		for col := 0; col < row.LastCol(); col++ {
			cell := row.Col(col)
			switch {
			case strings.Contains(strings.ToLower(cleanValue(cell)), "resumo dos saldos"):
				fmt.Println("Ativos em custódia:")
				headers := []string{}
				lineHeaders := line + 2
				rowHeader := sheet.Row(int(lineHeaders))
				for colH := uint16(col); colH < sheet.MaxRow; colH++ {
					c := rowHeader.Col(int(colH))
					if c != "" {
						headers = append(headers, c)
					}
				}
				lineActives := lineHeaders + 1
				var assets []*Stock
			stopAssets:
				for ; ; lineActives++ {
					row := sheet.Row(int(lineActives))

					var currentPos int

					stock := &Stock{}
					valueOf := reflect.ValueOf(stock).Elem()
					for cellA := 0; cellA < row.LastCol(); cellA++ {
						cell = row.Col(cellA)
						// Point to stop loop
						if strings.ToLower(cleanValue(cell)) == "valorização em reais" {
							break stopAssets
						}

						valueCleaned := cleanValue(cell)
						if valueCleaned != "" && valueCleaned != "#" && currentPos < valueOf.NumField() {

							if valueOf.Field(currentPos).Type().Kind() == reflect.Int {
								v, err := strconv.ParseInt(cell, 10, 64)
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

				fmt.Println("|" + strings.Join(headers, "\t|") + "|")

				for _, activePart := range assets {
					v := reflect.ValueOf(activePart).Elem()
					item := []string{}
					for i := 0; i < v.NumField(); i++ {
						item = append(item, fmt.Sprintf("%v", v.Field(i).Interface()))
					}
					fmt.Println("|" + strings.Join(item, "\t|") + "|")
				}

				months = append(months, assets)

			default:

			}

			/*if strings.ToLower(strings.Trim(cell.Value, " \t\n")) == "total creditado" {
				fmt.Println("Total Creditado", line, col, cell)
				foundCredit = true
				continue
			}
			if foundCredit && cell.Value != "" {
				fmt.Println("Valor", line, col, cell)
			}
			*/
		}
	}
	json.NewEncoder(os.Stdout).Encode(months)
}
