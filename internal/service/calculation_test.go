package service

import (
	"math"
	"testing"
	"time"
)

func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 1e-9
}

func TestCalculationService_DailyRate(t *testing.T) {
	c := NewCalculationService()

	tests := []struct {
		name  string
		total float64
		days  int
		want  float64
	}{
		{"30 days", 3000, 30, 100},
		{"31 days", 3100, 31, 100},
		{"zero days returns zero", 1000, 0, 0},
		{"negative days returns zero", 1000, -5, 0},
		{"fractional", 1000, 28, 1000.0 / 28.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.DailyRate(tt.total, tt.days)
			if !almostEqual(got, tt.want) {
				t.Errorf("DailyRate(%v, %d) = %v, want %v", tt.total, tt.days, got, tt.want)
			}
		})
	}
}

func TestCalculationService_PeriodBounds(t *testing.T) {
	c := NewCalculationService()

	tests := []struct {
		name      string
		input     time.Time
		wantStart time.Time
		wantEnd   time.Time
	}{
		{
			name:      "exactly the 25th — period starts that day",
			input:     date(2026, time.March, 25),
			wantStart: date(2026, time.March, 25),
			wantEnd:   date(2026, time.April, 25),
		},
		{
			name:      "26th — within the period that started on the 25th",
			input:     date(2026, time.March, 26),
			wantStart: date(2026, time.March, 25),
			wantEnd:   date(2026, time.April, 25),
		},
		{
			name:      "1st — period started on 25th of previous month",
			input:     date(2026, time.April, 1),
			wantStart: date(2026, time.March, 25),
			wantEnd:   date(2026, time.April, 25),
		},
		{
			name:      "24th — still in previous period",
			input:     date(2026, time.April, 24),
			wantStart: date(2026, time.March, 25),
			wantEnd:   date(2026, time.April, 25),
		},
		{
			name:      "January 5th — period crosses year boundary",
			input:     date(2026, time.January, 5),
			wantStart: date(2025, time.December, 25),
			wantEnd:   date(2026, time.January, 25),
		},
		{
			name:      "December 25th — period crosses into next year",
			input:     date(2025, time.December, 25),
			wantStart: date(2025, time.December, 25),
			wantEnd:   date(2026, time.January, 25),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStart, gotEnd := c.PeriodBounds(tt.input)
			if !gotStart.Equal(tt.wantStart) {
				t.Errorf("start = %v, want %v", gotStart, tt.wantStart)
			}
			if !gotEnd.Equal(tt.wantEnd) {
				t.Errorf("end = %v, want %v", gotEnd, tt.wantEnd)
			}
		})
	}
}

func TestCalculationService_DaysInPeriod(t *testing.T) {
	c := NewCalculationService()

	tests := []struct {
		name  string
		input time.Time
		want  int
	}{
		{"Mar 25 → Apr 25 = 31 days", date(2026, time.March, 25), 31},
		{"Apr 25 → May 25 = 30 days", date(2026, time.April, 25), 30},
		{"Feb 25 → Mar 25 = 28 in non-leap year", date(2026, time.February, 25), 28},
		{"Feb 25 → Mar 25 = 29 in leap year", date(2024, time.February, 25), 29},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.DaysInPeriod(tt.input)
			if got != tt.want {
				t.Errorf("DaysInPeriod(%v) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestCalculationService_RealignmentDate(t *testing.T) {
	c := NewCalculationService()

	tests := []struct {
		name   string
		moveIn time.Time
		want   time.Time
	}{
		{
			name:   "move-in Mar 20 → advance2 ends May 19 → realign May 25",
			moveIn: date(2026, time.March, 20),
			want:   date(2026, time.May, 25),
		},
		{
			name:   "move-in Apr 5 → advance2 ends Jun 4 → realign Jun 25",
			moveIn: date(2026, time.April, 5),
			want:   date(2026, time.June, 25),
		},
		{
			name:   "move-in Mar 25 → advance2 ends May 24 → realign May 25 (next day)",
			moveIn: date(2026, time.March, 25),
			want:   date(2026, time.May, 25),
		},
		{
			name:   "move-in Mar 26 → advance2 ends May 25 → realign Jun 25 (must be strictly after)",
			moveIn: date(2026, time.March, 26),
			want:   date(2026, time.June, 25),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.RealignmentDate(tt.moveIn)
			if !got.Equal(tt.want) {
				t.Errorf("RealignmentDate(%v) = %v, want %v",
					tt.moveIn.Format("2006-01-02"), got, tt.want)
			}
		})
	}
}

func TestCalculationService_RealignmentDays(t *testing.T) {
	c := NewCalculationService()

	tests := []struct {
		name   string
		moveIn time.Time
		want   int
	}{
		{"Mar 20: advance2 ends May 19 → 6 days to May 25", date(2026, time.March, 20), 6},
		{"Apr 5:  advance2 ends Jun 4 → 21 days to Jun 25", date(2026, time.April, 5), 21},
		{"Mar 25: advance2 ends May 24 → 1 day to May 25", date(2026, time.March, 25), 1},
		{"Mar 26: advance2 ends May 25 → 31 days to Jun 25 (skip same-day)", date(2026, time.March, 26), 31},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := c.RealignmentDays(tt.moveIn)
			if got != tt.want {
				t.Errorf("RealignmentDays(%v) = %d, want %d",
					tt.moveIn.Format("2006-01-02"), got, tt.want)
			}
		})
	}
}
