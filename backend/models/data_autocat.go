package models

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type AutoCatOverride struct {
	ColumnName  string `json:"column_name"`
	ColumnValue string `json:"override_value"`
	Order       int    `json:"override_order"`
}

type AutoCatCriterion struct {
	ColumnName         string           `json:"column_name"`
	Operator           string           `json:"operator"`
	FilterValueString  *string          `json:"filter_value_text,omitempty"`
	FilterValueDecimal *decimal.Decimal `json:"filter_value_numeric,omitempty"`
	FilterValueTime    *time.Time       `json:"filter_value_timestamptz,omitempty"`
	Order              int              `json:"criteria_order"`
}

type AutoCatRule struct {
	Criteria  []AutoCatCriterion
	Overrides []AutoCatOverride
}

func getValueName(str string) string {
	words := strings.Split(str, "_")
	caser := cases.Title(language.English)
	for i := range words {
		if words[i] == "id" {
			words[i] = "ID"
		} else {
			words[i] = caser.String(words[i])
		}
	}
	converted := strings.Join(words, "")
	return converted
}

func matchDecimal(c decimal.Decimal, t decimal.Decimal, operator string) bool {
	switch operator {
	case "equals":
		return t.Equal(c)
	case "min":
		return t.GreaterThanOrEqual(c)
	case "max":
		return t.LessThanOrEqual(c)
	default:
		return false
	}
}

func matchTime(c time.Time, t time.Time, operator string) bool {
	switch operator {
	case "min":
		return t.After(c) || t.Equal(c)
	case "max":
		return t.Before(c) || t.Equal(c)
	default:
		return false
	}
}

func matchString(c string, t string, operator string) bool {
	tLower := strings.ToLower(t)
	cLower := strings.ToLower(c)

	switch operator {
	case "equals":
		return tLower == cLower
	case "contains":
		return strings.Contains(tLower, cLower)
	case "starts_with":
		return strings.HasPrefix(tLower, cLower)
	case "ends_with":
		return strings.HasSuffix(tLower, cLower)
	case "regex":
		match, _ := regexp.MatchString(c, tLower)
		return match
	default:
		return false
	}
}

func (c AutoCatCriterion) matches(txn TransactionRow) bool {
	txnVal := reflect.ValueOf(txn)
	colVal := txnVal.FieldByName(getValueName(c.ColumnName))

	if c.FilterValueDecimal != nil {
		decimalVal, _ := colVal.Interface().(decimal.Decimal)
		return matchDecimal(*c.FilterValueDecimal, decimalVal, c.Operator)
	} else if c.FilterValueTime != nil {
		timeVal, _ := colVal.Interface().(time.Time)
		return matchTime(*c.FilterValueTime, timeVal, c.Operator)
	} else if c.FilterValueString != nil {
		stringVal := colVal.String()
		return matchString(*c.FilterValueString, stringVal, c.Operator)
	}
	return false
}

func (c AutoCatCriterion) String() string {
	if c.FilterValueString != nil {
		return fmt.Sprintf("%v %v %v", c.ColumnName, c.Operator, *c.FilterValueString)
	} else if c.FilterValueDecimal != nil {
		return fmt.Sprintf("%v %v %v", c.ColumnName, c.Operator, *c.FilterValueDecimal)
	} else if c.FilterValueTime != nil {
		return fmt.Sprintf("%v %v %v", c.ColumnName, c.Operator, *c.FilterValueTime)
	} else {
		return "<invalid criterion>"
	}
}

func (rule AutoCatRule) MatchesRow(txn TransactionRow) bool {
	sort.SliceStable(rule.Criteria, func(i, j int) bool {
		return rule.Criteria[i].Order < rule.Criteria[j].Order
	})

	for _, c := range rule.Criteria {
		match := c.matches(txn)
		if !match {
			return false
		}
	}
	return true
}

func (rule AutoCatRule) ApplyOverrides(txn TransactionRow) (TransactionRow, error) {
	sort.SliceStable(rule.Overrides, func(i, j int) bool {
		return rule.Overrides[i].Order < rule.Overrides[j].Order
	})

	newTxn := txn

	for _, o := range rule.Overrides {
		ovrVal := reflect.Indirect(reflect.ValueOf(&newTxn).Elem())
		ovrField := ovrVal.FieldByName(getValueName(o.ColumnName))
		if ovrField.IsValid() && ovrField.CanSet() {
			ovrVal := reflect.ValueOf(o.ColumnValue)
			ovrField.Set(ovrVal)
		} else {
			return TransactionRow{}, fmt.Errorf("field %v is invalid", ovrField)
		}
	}

	return newTxn, nil
}
