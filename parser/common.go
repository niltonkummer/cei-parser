package parser

import (
	"errors"
	"github.com/shopspring/decimal"
	"path/filepath"
	"strings"
	"time"
)

const (
	MagicKeyResumoSaldos                    = "resumo dos saldos"
	MagicKeyProventosEmDinheiroCreditados   = "proventos em dinheiro - creditados"
	MagicKeyInformacoesDeNegociacaoDeAtivos = "informações de negociação de ativos"
	MagicKeyTotalCreditado                  = "total creditado"
	MagicKeyValorizacaoEmReais              = "valorização em reais"

	PrefixToGetDate = "resumo dos saldos em"

	CellRowDate = 17
	CellColDate = 5
)

type FinancialAssets struct {
	MonthPeriod time.Time  `json:"month-period"`
	Stocks      []*Stock   `json:"stocks,omitempty"`
	Proceeds    []*Proceed `json:"proceeds,omitempty"`
	Deals       []*Deal    `json:"deals,omitempty"`
}

type Stock struct {
	Name    string          `json:"name,omitempty" column:"Ativo"`
	Spec    string          `json:"spec,omitempty" column:"Especif."`
	Code    string          `json:"code" column:"Cód. Neg."`
	Balance int             `json:"balance" column:"Saldo"`
	Price   decimal.Decimal `json:"price" column:"Cotação"`
	Amount  decimal.Decimal `json:"amount" column:"Valor"`
}

type Deal struct {
	Code         string          `json:"code" column:"Cód"`
	Date         time.Time       `json:"date" column:"Data Negócio"`
	AmountBuy    int             `json:"amount-buy" column:"Qtde.Compra"`
	AmountSell   int             `json:"amount-sell" column:"Qtd.Venda"`
	BuyAvgPrice  decimal.Decimal `json:"buy-avg-price" column:"Preço Médio Compra"`
	SellAvgPrice decimal.Decimal `json:"sell-avg-price" column:"Preço Médio Venda"`
	NetAmount    int             `json:"net-amount" column:"Qtde. Liquida"`
	Position     string          `json:"position" column:"Posição"`
}

type Proceed struct {
	Name     string          `json:"name" column:"Ativo"`
	Spec     string          `json:"spec" column:"Especif."`
	Code     string          `json:"code" column:"Cód. Neg."`
	Credited decimal.Decimal `json:"credited" column:"Creditado No Mês"`
}

type Parser interface {
	Decode() (*FinancialAssets, error)
}

func cleanValue(v string) string {
	return strings.Trim(v, " \t\n")
}

func normalizeToLower(v string) string {
	return strings.ToLower(cleanValue(v))
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
