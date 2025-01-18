package models

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

type IntervalType string

const (
	PAST_30_DAYS   IntervalType = "30d"
	THIS_WEEK      IntervalType = "week"
	THIS_MONTH     IntervalType = "quarter"
	THIS_QUARTER   IntervalType = "month"
	THIS_YEAR      IntervalType = "year"
	LAST_12_MONTHS IntervalType = "12m"
	LAST_WEEK      IntervalType = "last-week"
	LAST_MONTH     IntervalType = "last-month"
	LAST_QUARTER   IntervalType = "last-quarter"
	LAST_YEAR      IntervalType = "last-year"
)

type IntervalPair struct {
	Start time.Time
	End   time.Time
}

func NewIntervalPair(i IntervalType) (IntervalPair, error) {
	today := time.Now().UTC()

	switch i {
	case PAST_30_DAYS:
		start := today.Add(-30 * 24 * time.Hour)
		return IntervalPair{
			Start: start,
			End:   today,
		}, nil

	case THIS_WEEK:
		// Get the start of the current week (Monday)
		weekday := today.Weekday()
		if weekday == time.Sunday {
			weekday = 7
		}
		start := today.Add(-time.Duration(weekday-1) * 24 * time.Hour)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
		return IntervalPair{
			Start: start,
			End:   today,
		}, nil

	case THIS_MONTH:
		start := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, time.UTC)
		return IntervalPair{
			Start: start,
			End:   today,
		}, nil

	case THIS_QUARTER:
		quarter := (today.Month() - 1) / 3
		firstMonthOfQuarter := quarter*3 + 1
		start := time.Date(today.Year(), firstMonthOfQuarter, 1, 0, 0, 0, 0, time.UTC)
		return IntervalPair{
			Start: start,
			End:   today,
		}, nil

	case THIS_YEAR:
		start := time.Date(today.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		return IntervalPair{
			Start: start,
			End:   today,
		}, nil

	case LAST_12_MONTHS:
		start := today.AddDate(0, -12, 0)
		return IntervalPair{
			Start: start,
			End:   today,
		}, nil

	case LAST_WEEK:
		// Get the start of last week (Monday)
		weekday := today.Weekday()
		if weekday == time.Sunday {
			weekday = 7
		}
		start := today.Add(-time.Duration(weekday+6) * 24 * time.Hour)
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
		end := start.Add(7 * 24 * time.Hour)
		return IntervalPair{
			Start: start,
			End:   end,
		}, nil

	case LAST_MONTH:
		end := time.Date(today.Year(), today.Month(), 1, 0, 0, 0, 0, time.UTC)
		start := end.AddDate(0, -1, 0)
		return IntervalPair{
			Start: start,
			End:   end,
		}, nil

	case LAST_QUARTER:
		quarter := (today.Month() - 1) / 3
		firstMonthOfQuarter := quarter*3 + 1
		end := time.Date(today.Year(), firstMonthOfQuarter, 1, 0, 0, 0, 0, time.UTC)
		start := end.AddDate(0, -3, 0)
		return IntervalPair{
			Start: start,
			End:   end,
		}, nil

	case LAST_YEAR:
		end := time.Date(today.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
		start := end.AddDate(-1, 0, 0)
		return IntervalPair{
			Start: start,
			End:   end,
		}, nil
	}

	return IntervalPair{}, errors.New("could not create IntervalPair")
}

func (i IntervalPair) CalcAvgDailyExpense(e decimal.Decimal) decimal.Decimal {
	dur := i.End.Sub(i.Start)
	days := int64(dur.Hours() / 24)
	return e.Div(decimal.NewFromInt(days))
}
