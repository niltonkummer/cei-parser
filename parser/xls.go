package parser

import (
	"github.com/extrame/xls"
	"time"

	"strings"
)

// TODO: Parse headers

type xlsDecoder struct {
	path string
	wb   *xls.WorkBook
}

func (x *xlsDecoder) Decode() (*FinancialAssets, error) {

	f, err := xls.Open(x.path, "utf-8")
	if err != nil {
		return nil, err
	}
	x.wb = f

	return x.parseXlsSheet(f.GetSheet(0)), nil

}

func (x *xlsDecoder) parseXlsSheet(sheet *xls.WorkSheet) (financial *FinancialAssets) {

	date, err := time.Parse("02/01/2006", strings.TrimSpace(
		strings.Replace(strings.ToLower(sheet.Row(CellRowDate).Col(CellColDate)), PrefixToGetDate, "", -1)))
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

	for line := uint16(0); line < sheet.MaxRow; line++ {
		row := sheet.Row(int(line))
		for col := 0; col < row.LastCol(); col++ {

			cell := row.Col(col)

			switch {
			case strings.Contains(normalizeToLower(cell), MagicKeyResumoSaldos):
				assets = x.parseCurrentAssets(sheet, line, col)

			case strings.Contains(normalizeToLower(cell), MagicKeyProventosEmDinheiroCreditados):
				proceeds = x.parseCreditedEarnings(sheet, line, col)

			case strings.Contains(normalizeToLower(cell), MagicKeyInformacoesDeNegociacaoDeAtivos):
				deals = x.parseAssetDeals(sheet, line, col)

			}

		}
	}

	financial.Stocks = assets
	financial.Proceeds = proceeds
	financial.Deals = deals

	return
}

func (x *xlsDecoder) parseAssetDeals(sheet *xls.WorkSheet, line uint16, col int) (deals []*Deal) {

	var header []string
	lineHeader := line + 2
	rowHeader := sheet.Row(int(lineHeader))
	for colH := uint16(col); colH < sheet.MaxRow; colH++ {
		c := rowHeader.Col(int(colH))
		if c != "" {
			header = append(header, cleanValue(c))
		}
	}

	lineDeals := lineHeader + 1

	fields := make(map[string]interface{}, 8)

stopDeals:
	for ; ; lineDeals++ {
		currentRow := sheet.Row(int(lineDeals))

		var currentPos int

		emptyLine := true

		for cellA := 0; cellA < currentRow.LastCol(); cellA++ {

			cell := cleanValue(currentRow.Col(cellA))
			if cell != "" && cell != "#" {

				emptyLine = false

				fields[header[currentPos]] = cell
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

func (x *xlsDecoder) parseCreditedEarnings(sheet *xls.WorkSheet, line uint16, col int) (list []*Proceed) {

	var header []string
	lineHeader := line + 1
	rowHeader := sheet.Row(int(lineHeader))
	for colH := uint16(col); colH < sheet.MaxRow; colH++ {
		c := rowHeader.Col(int(colH))
		if c != "" {
			header = append(header, c)
		}
	}
	lineProceeds := lineHeader + 1

	fields := make(map[string]interface{}, 8)

stopProceeds:
	for ; ; lineProceeds++ {
		currentRow := sheet.Row(int(lineProceeds))

		var currentPos int

		for cellA := 0; cellA < currentRow.LastCol(); cellA++ {
			cell := cleanValue(currentRow.Col(cellA))

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

func (x *xlsDecoder) parseCurrentAssets(sheet *xls.WorkSheet, line uint16, col int) (assets []*Stock) {

	var header []string
	lineHeader := line + 2

	rowHeader := sheet.Row(int(lineHeader))
	for colH := uint16(col); colH < sheet.MaxRow; colH++ {
		c := rowHeader.Col(int(colH))
		if c != "" {
			header = append(header, c)
		}
	}

	lineAssets := lineHeader + 1

	fields := make(map[string]interface{}, 8)

stopAssets:
	for ; ; lineAssets++ {
		currentRow := sheet.Row(int(lineAssets))

		var currentPos int

		for cellA := 0; cellA < currentRow.LastCol(); cellA++ {
			cell := cleanValue(currentRow.Col(cellA))

			// Point to stop the loop
			if normalizeToLower(cell) == MagicKeyValorizacaoEmReais {
				break stopAssets
			}

			if cell != "" && cell != "#" {
				fields[header[currentPos]] = cell
				currentPos++
			}
		}

		stock := &Stock{}
		fieldsToStruct(fields, stock)
		assets = append(assets, stock)
	}

	return
}
