package parser

import (
	"errors"
	"fmt"
	"github.com/tealeg/xlsx"
	"strings"
	"time"
)

// TODO: needs refatory it's out of date
type xlsxDecoder struct {
	path string
}

func (x *xlsxDecoder) Decode() (*FinancialAssets, error) {

	f, err := xlsx.OpenFile(x.path)
	if err != nil {
		return nil, err
	}

	for _, sheet := range f.Sheets {
		return x.parseSheet(sheet), nil
	}
	return nil, errors.New("empty spread sheet")
}

func (x *xlsxDecoder) parseSheet(sheet *xlsx.Sheet) (financial *FinancialAssets) {

	date, err := time.Parse("02/01/2006", strings.TrimSpace(strings.Replace(strings.ToLower(sheet.Row(CellRowDate).Cells[CellColDate].String()), PrefixToGetDate, "", -1)))
	if err != nil {
		return nil
	}

	var (
		assets   []*Stock
		proceeds []*Proceed
		deals    []*Deal
	)

	financial = &FinancialAssets{
		MonthPeriod: date,
	}

	for line, row := range sheet.Rows {
		for col, cell := range row.Cells {

			switch {
			case strings.Contains(normalizeToLower(cell.String()), MagicKeyResumoSaldos):
				assets = x.parseCurrentAssets(sheet, line, col)

			case strings.Contains(normalizeToLower(cell.String()), MagicKeyProventosEmDinheiroCreditados):
				proceeds = x.parseCreditedEarnings(sheet, line, col)

			case strings.Contains(normalizeToLower(cell.String()), MagicKeyInformacoesDeNegociacaoDeAtivos):
				deals = x.parseAssetDeals(sheet, line, col)
			}
		}
	}
	financial.Stocks = assets
	financial.Proceeds = proceeds
	financial.Deals = deals

	return
}

func (x *xlsxDecoder) parseCurrentAssets(sheet *xlsx.Sheet, line int, col int) (assets []*Stock) {

	lineHeader := line + 2
	header := x.getHeader(sheet, lineHeader, col)

	lineAssets := lineHeader + 1

	fields := make(map[string]interface{}, 8)

stopAssets:
	for _, row := range sheet.Rows[lineAssets:] {
		var currentPos int

		for _, cell := range row.Cells {
			cell.Value = cleanValue(cell.Value)

			// Point to stop loop
			if normalizeToLower(cell.String()) == MagicKeyValorizacaoEmReais {
				break stopAssets
			}

			value := cell.String()
			if value != "" && value != "#" {

				if f, err := cell.Int64(); err == nil {
					value = fmt.Sprintf("%d", f)
				} else if f, err := cell.Float(); err == nil {
					value = fmt.Sprintf("%.2f", f*100/100)
				}

				fields[header[currentPos]] = value
				currentPos++
			}
		}
		stock := &Stock{}
		fieldsToStruct(fields, stock)
		assets = append(assets, stock)
	}
	return
}

func (x *xlsxDecoder) parseAssetDeals(sheet *xlsx.Sheet, line int, col int) (deals []*Deal) {

	lineHeader := line + 2

	header := x.getHeader(sheet, lineHeader, col)

	lineDeals := lineHeader + 1

	fields := make(map[string]interface{}, 8)

stopDeals:
	for _, row := range sheet.Rows[lineDeals:] {

		var currentPos int

		emptyLine := true

		for _, cell := range row.Cells {

			cell.Value = cleanValue(cell.Value)

			var value interface{} = cell.String()
			if value != "" && value != "#" {
				if cell.IsTime() {
					t, err := cell.GetTime(false)
					if err == nil {
						value = t
					}
				}

				emptyLine = false
				fields[header[currentPos]] = value
				currentPos++
			}

		}

		if emptyLine {
			break stopDeals
		}

		deal := &Deal{}
		fieldsToStruct(fields, deal)
		deals = append(deals, deal)
	}
	return
}

func (x *xlsxDecoder) parseCreditedEarnings(sheet *xlsx.Sheet, line int, col int) (list []*Proceed) {

	lineHeader := line + 1

	header := x.getHeader(sheet, lineHeader, col)

	lineProceeds := lineHeader + 1

	fields := make(map[string]interface{}, 8)

stopProceeds:
	for _, row := range sheet.Rows[lineProceeds:] {

		var currentPos int

		for _, cellA := range row.Cells {
			cell := cleanValue(cellA.String())

			// Point to stop loop

			if normalizeToLower(cell) == MagicKeyTotalCreditado {
				break stopProceeds
			}

			if cell != "" && cell != "#" {

				fields[header[currentPos]] = cell

				currentPos++
			}
		}
		profit := &Proceed{}
		fieldsToStruct(fields, profit)
		list = append(list, profit)
	}
	return
}

func (x *xlsxDecoder) getHeader(sheet *xlsx.Sheet, lineHeaders int, col int) []string {

	var header []string
	for _, c := range sheet.Rows[lineHeaders].Cells[col:] {
		if c.String() != "" {
			header = append(header, c.String())
		}
	}
	return header
}
