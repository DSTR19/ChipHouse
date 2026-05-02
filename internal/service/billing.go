package service

import (
	"time"

	"ChipHouse/internal/domain"
)

// ChargeKind classifies the billing event so callers (and tests) can reason
// about the schedule explicitly.
type ChargeKind int

const (
	// ChargeMoveIn — first ever payment: 30-day rent advance + deposit (= rent).
	ChargeMoveIn ChargeKind = iota
	// ChargeAdvanceRenewal — second payment, on moveIn+30d: rent advance + utilities.
	ChargeAdvanceRenewal
	// ChargeRealignment — one-shot, on the first 25th after the second advance ends:
	// adjustment only (no rent, no utilities) to shift the tenant onto the
	// building-wide 25→25 schedule.
	ChargeRealignment
	// ChargeRegular — steady-state monthly rent + utilities, every 25th.
	ChargeRegular
)

type BillingService struct {
	calc *CalculationService
}

func NewBillingService(calc *CalculationService) *BillingService {
	return &BillingService{calc: calc}
}

// MoveInCharge generates the very first charge: 30-day rent advance + deposit.
// Deposit equals rent and is held until move-out.
func (b *BillingService) MoveInCharge(studio domain.Studio, tenant domain.Tenant) domain.Charge {
	moveIn := startOfDay(tenant.MoveIn)
	return domain.Charge{
		StudioID:    studio.ID,
		TenantID:    tenant.ID,
		PeriodStart: moveIn,
		PeriodEnd:   moveIn.AddDate(0, 0, 30),
		Rent:        studio.Rent,
		Utilities:   0,
		Adjustment:  studio.Rent, // deposit, equals rent
		Total:       2 * studio.Rent,
		IssuedAt:    time.Now(),
	}
}

// AdvanceRenewalCharge generates the second payment, on moveIn+30d:
// rent advance for the next 30 days + utilities for the elapsed period.
func (b *BillingService) AdvanceRenewalCharge(studio domain.Studio, tenant domain.Tenant, utilities float64) domain.Charge {
	start := startOfDay(tenant.MoveIn).AddDate(0, 0, 30)
	end := start.AddDate(0, 0, 30)
	return domain.Charge{
		StudioID:    studio.ID,
		TenantID:    tenant.ID,
		PeriodStart: start,
		PeriodEnd:   end,
		Rent:        studio.Rent,
		Utilities:   utilities,
		Adjustment:  0,
		Total:       studio.Rent + utilities,
		IssuedAt:    time.Now(),
	}
}

// RealignmentCharge generates the one-shot charge that shifts the tenant onto
// the building-wide 25→25 schedule. It contains *only* an adjustment for the
// gap between the end of the second advance and the realignment 25th.
//
// adjustment = daily_rate × gap_days, where daily_rate = rent / 30.
// The sign is always positive (additional charge) under this model — the
// 25th is always after the advance end, never before.
func (b *BillingService) RealignmentCharge(studio domain.Studio, tenant domain.Tenant) domain.Charge {
	advance2End := startOfDay(tenant.MoveIn).AddDate(0, 0, 60)
	realign := b.calc.RealignmentDate(tenant.MoveIn)
	gapDays := b.calc.RealignmentDays(tenant.MoveIn)
	dailyRate := studio.Rent / 30.0
	adjustment := dailyRate * float64(gapDays)

	return domain.Charge{
		StudioID:    studio.ID,
		TenantID:    tenant.ID,
		PeriodStart: advance2End,
		PeriodEnd:   realign,
		Rent:        0,
		Utilities:   0,
		Adjustment:  adjustment,
		Total:       adjustment,
		IssuedAt:    time.Now(),
	}
}

// RegularCharge generates a steady-state monthly charge for the 25→25 period
// containing periodAnchor. No adjustment — realignment has already happened.
func (b *BillingService) RegularCharge(studio domain.Studio, tenant domain.Tenant, utilities float64, periodAnchor time.Time) domain.Charge {
	start, end := b.calc.PeriodBounds(periodAnchor)
	return domain.Charge{
		StudioID:    studio.ID,
		TenantID:    tenant.ID,
		PeriodStart: start,
		PeriodEnd:   end,
		Rent:        studio.Rent,
		Utilities:   utilities,
		Adjustment:  0,
		Total:       studio.Rent + utilities,
		IssuedAt:    time.Now(),
	}
}
