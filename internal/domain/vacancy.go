package domain

import "time"

type VacancyPeriod struct {
	ID        int
	StudioID  int
	StartDate time.Time
	EndDate   *time.Time
	LostRent  float64
}
