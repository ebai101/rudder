package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

type IntervalType string

const (
	PAST_30_DAYS   IntervalType = "PAST_30_DAYS"
	THIS_WEEK      IntervalType = "THIS_WEEK"
	THIS_QUARTER   IntervalType = "THIS_QUARTER"
	THIS_MONTH     IntervalType = "THIS_MONTH"
	THIS_YEAR      IntervalType = "THIS_YEAR"
	LAST_12_MONTHS IntervalType = "LAST_12_MONTHS"
	LAST_WEEK      IntervalType = "LAST_WEEK"
	LAST_QUARTER   IntervalType = "LAST_QUARTER"
	LAST_MONTH     IntervalType = "LAST_MONTH"
	LAST_YEAR      IntervalType = "LAST_YEAR"
)

type IntervalPair struct {
	Start time.Time
	End   time.Time
}

func NewIntervalPair(i IntervalType) (IntervalPair, error) {
	today := time.Now().UTC()

	switch i {
	case PAST_30_DAYS:
		start := today.Add(-1 * 30 * 24 * time.Hour)
		pair := IntervalPair{
			Start: start,
			End:   today,
		}
		fmt.Println(pair)
		return pair, nil
	}

	return IntervalPair{}, errors.New("could not create IntervalPair")
}

func (i IntervalPair) CalcAvgDailyExpense(e decimal.Decimal) decimal.Decimal {
	dur := i.End.Sub(i.Start)
	days := int64(dur.Hours() / 24)
	return e.Div(decimal.NewFromInt(days))
}
