package models

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type AutocatOverride struct {
	ColumnName  string `json:"column_name"`
	ColumnValue string `json:"override_value"`
	Order       int    `json:"override_order"`
}

type AutocatCriterion struct {
	ColumnName         string           `json:"column_name"`
	Operator           string           `json:"operator"`
	FilterValueString  *string          `json:"filter_value_text,omitempty"`
	FilterValueDecimal *decimal.Decimal `json:"filter_value_numeric,omitempty"`
	FilterValueTime    *time.Time       `json:"filter_value_timestamptz,omitempty"`
	Order              int              `json:"criteria_order"`
}

type AutocatRule struct {
	RuleID    int64
	Criteria  []AutocatCriterion
	Overrides []AutocatOverride
}

func (c AutocatCriterion) String() string {
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

func (c AutocatOverride) String() string {
	return fmt.Sprintf("%v = %v", c.ColumnName, c.ColumnValue)
}
