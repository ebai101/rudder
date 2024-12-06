package models

type AutoCatOverride struct {
	ColumnName  string `json:"column_name"`
	ColumnValue string `json:"override_value"`
	Order       int    `json:"override_order"`
}

type AutoCatCriterion struct {
	ColumnName  string `json:"column_name"`
	Operator    string `json:"operator"`
	FilterValue string `json:"filter_value"`
	Order       int    `json:"criteria_order"`
}

type AutoCatRule struct {
	Criteria  []AutoCatCriterion
	Overrides []AutoCatOverride
}
