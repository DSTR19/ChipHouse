package domain

import "time"

type Charge struct {
	ID          int
	StudioID    int
	TenantID    int
	PeriodStart time.Time
	PeriodEnd   time.Time
	Rent        float64
	Utilities   float64
	Adjustment  float64
	Total       float64
	Paid        float64
	IssuedAt    time.Time
}
