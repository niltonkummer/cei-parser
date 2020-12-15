package parser

import (
	"fmt"
	"github.com/shopspring/decimal"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func fieldsToStruct(fields map[string]interface{}, str interface{}) {

	valueOf := reflect.ValueOf(str).Elem()

	for field := 0; field < valueOf.NumField(); field++ {

		key, ok := valueOf.Type().Field(field).Tag.Lookup("column")

		if !ok {
			continue
		}

		value := fields[key]

		currentPos := field
		if valueOf.Field(currentPos).Type().Kind() == reflect.Int {
			v, err := strconv.ParseInt(fmt.Sprintf("%v", value), 10, 64)
			if err == nil {
				valueOf.Field(currentPos).SetInt(v)
			}
		} else if valueOf.Field(currentPos).Type().Name() == "Time" {
			t, ok := value.(time.Time)
			if !ok {
				a, _ := strconv.Atoi(fmt.Sprintf("%v", value))
				v := time.Date(1900, 1, a-1, 0, 0, 0, 0, time.UTC)
				valueOf.Field(currentPos).Set(reflect.ValueOf(v))
			}
			valueOf.Field(currentPos).Set(reflect.ValueOf(t))
		} else if valueOf.Field(currentPos).Type().Name() == "Decimal" {

			v, err := decimal.NewFromString(fmt.Sprintf("%v", value))
			if err != nil {
				value = strings.Replace(strings.Replace(fmt.Sprintf("%v", value), ".", "", -1), ",", ".", 1)
				v, _ = decimal.NewFromString(fmt.Sprintf("%v", value))

			}
			vInt := v.Mul(decimal.NewFromFloat(100)).IntPart()
			v = decimal.NewFromInt(vInt).Div(decimal.NewFromFloat(100))
			valueOf.Field(currentPos).Set(reflect.ValueOf(v))
		} else {
			valueOf.Field(currentPos).SetString(fmt.Sprintf("%v", value))
		}

	}
}
