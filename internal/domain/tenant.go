package domain

import "time"

type Tenant struct {
	ID       int
	FullName string
	StudioID int
	MoveIn   time.Time
	MoveOut  *time.Time
	Persons  int
	IsActive bool
}
