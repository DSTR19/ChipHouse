package domain

import "time"

type ComplaintStatus string

const (
	ComplaintOpen     ComplaintStatus = "OPEN"
	ComplaintResolved ComplaintStatus = "RESOLVED"
)

type Complaint struct {
	ID         int
	TenantID   int
	StudioID   int
	Subject    string
	Body       string
	Status     ComplaintStatus
	CreatedAt  time.Time
	ResolvedAt *time.Time
}
