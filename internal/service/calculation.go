package service

import "time"

type CalculationService struct{}

func NewCalculationService() *CalculationService {
	return &CalculationService{}
}

func (c *CalculationService) DailyRate(total float64, days int) float64 {
	if days <= 0 {
		return 0
	}
	return total / float64(days)
}

// PeriodBounds returns the [start, end) of the billing period containing date d.
// Billing period: 25th of month → 25th of next month.
func (c *CalculationService) PeriodBounds(d time.Time) (time.Time, time.Time) {
	year, month, day := d.Date()
	loc := d.Location()
	var start time.Time
	if day >= 25 {
		start = time.Date(year, month, 25, 0, 0, 0, 0, loc)
	} else {
		start = time.Date(year, month-1, 25, 0, 0, 0, 0, loc)
	}
	end := start.AddDate(0, 1, 0)
	return start, end
}

// DaysInPeriod returns the number of days in the billing period containing d.
func (c *CalculationService) DaysInPeriod(d time.Time) int {
	start, end := c.PeriodBounds(d)
	return int(end.Sub(start).Hours() / 24)
}

// RealignmentDate is the 25th-of-month on which the tenant's payment schedule
// is aligned to the building-wide billing day. It is the first 25th strictly
// after the second advance ends (moveIn + 60 days).
func (c *CalculationService) RealignmentDate(moveIn time.Time) time.Time {
	advance2End := startOfDay(moveIn).AddDate(0, 0, 60)

	year, month, _ := advance2End.Date()
	loc := advance2End.Location()
	candidate := time.Date(year, month, 25, 0, 0, 0, 0, loc)
	if !candidate.After(advance2End) {
		candidate = candidate.AddDate(0, 1, 0)
	}
	return candidate
}

// RealignmentDays is the size, in days, of the gap between the second advance
// end and the realignment 25th. Always positive (or zero on the boundary).
func (c *CalculationService) RealignmentDays(moveIn time.Time) int {
	advance2End := startOfDay(moveIn).AddDate(0, 0, 60)
	realign := c.RealignmentDate(moveIn)
	return int(realign.Sub(advance2End).Hours() / 24)
}

func startOfDay(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}
