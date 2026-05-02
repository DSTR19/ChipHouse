package service

import (
	"testing"
	"time"

	"ChipHouse/internal/domain"
)

func TestBillingService_MoveInCharge_RentPlusDeposit(t *testing.T) {
	calc := NewCalculationService()
	b := NewBillingService(calc)

	studio := domain.Studio{ID: 1, Rent: 1000}
	tenant := domain.Tenant{ID: 42, StudioID: 1, MoveIn: date(2026, time.March, 20)}

	charge := b.MoveInCharge(studio, tenant)

	if charge.Rent != 1000 {
		t.Errorf("Rent = %v, want 1000", charge.Rent)
	}
	// deposit (= rent) is carried in Adjustment so the move-in payment doubles up.
	if charge.Adjustment != 1000 {
		t.Errorf("Deposit (Adjustment) = %v, want 1000", charge.Adjustment)
	}
	if charge.Total != 2000 {
		t.Errorf("Total = %v, want 2000 (rent + deposit)", charge.Total)
	}
	if !charge.PeriodStart.Equal(date(2026, time.March, 20)) {
		t.Errorf("PeriodStart = %v, want 2026-03-20", charge.PeriodStart)
	}
	if !charge.PeriodEnd.Equal(date(2026, time.April, 19)) {
		t.Errorf("PeriodEnd = %v, want 2026-04-19 (moveIn+30d)", charge.PeriodEnd)
	}
}

func TestBillingService_AdvanceRenewalCharge_RentPlusUtilities(t *testing.T) {
	calc := NewCalculationService()
	b := NewBillingService(calc)

	studio := domain.Studio{ID: 1, Rent: 1000}
	tenant := domain.Tenant{ID: 1, StudioID: 1, MoveIn: date(2026, time.March, 20)}

	charge := b.AdvanceRenewalCharge(studio, tenant, 250)

	if !charge.PeriodStart.Equal(date(2026, time.April, 19)) {
		t.Errorf("PeriodStart = %v, want 2026-04-19 (moveIn+30d)", charge.PeriodStart)
	}
	if !charge.PeriodEnd.Equal(date(2026, time.May, 19)) {
		t.Errorf("PeriodEnd = %v, want 2026-05-19 (moveIn+60d)", charge.PeriodEnd)
	}
	if charge.Rent != 1000 || charge.Utilities != 250 {
		t.Errorf("rent=%v utilities=%v, want 1000/250", charge.Rent, charge.Utilities)
	}
	if charge.Adjustment != 0 {
		t.Errorf("Adjustment = %v, want 0", charge.Adjustment)
	}
	if charge.Total != 1250 {
		t.Errorf("Total = %v, want 1250", charge.Total)
	}
}

func TestBillingService_RealignmentCharge_AdjustmentOnly(t *testing.T) {
	calc := NewCalculationService()
	b := NewBillingService(calc)

	tests := []struct {
		name           string
		moveIn         time.Time
		rent           float64
		wantStart      time.Time
		wantEnd        time.Time
		wantAdjustment float64
	}{
		{
			name:           "move-in Mar 20: advance2 ends May 19 → realign May 25 → +6 days × 1000/30",
			moveIn:         date(2026, time.March, 20),
			rent:           3000,
			wantStart:      date(2026, time.May, 19),
			wantEnd:        date(2026, time.May, 25),
			wantAdjustment: 6 * (3000.0 / 30.0),
		},
		{
			name:           "move-in Apr 5: advance2 ends Jun 4 → realign Jun 25 → +21 days × 3000/30",
			moveIn:         date(2026, time.April, 5),
			rent:           3000,
			wantStart:      date(2026, time.June, 4),
			wantEnd:        date(2026, time.June, 25),
			wantAdjustment: 21 * (3000.0 / 30.0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			studio := domain.Studio{ID: 1, Rent: tt.rent}
			tenant := domain.Tenant{ID: 1, StudioID: 1, MoveIn: tt.moveIn}

			charge := b.RealignmentCharge(studio, tenant)

			if !charge.PeriodStart.Equal(tt.wantStart) {
				t.Errorf("PeriodStart = %v, want %v", charge.PeriodStart, tt.wantStart)
			}
			if !charge.PeriodEnd.Equal(tt.wantEnd) {
				t.Errorf("PeriodEnd = %v, want %v", charge.PeriodEnd, tt.wantEnd)
			}
			if charge.Rent != 0 || charge.Utilities != 0 {
				t.Errorf("realignment must carry no rent/utilities; got rent=%v util=%v",
					charge.Rent, charge.Utilities)
			}
			if !almostEqual(charge.Adjustment, tt.wantAdjustment) {
				t.Errorf("Adjustment = %v, want %v", charge.Adjustment, tt.wantAdjustment)
			}
			if !almostEqual(charge.Total, tt.wantAdjustment) {
				t.Errorf("Total = %v, want %v (= adjustment)", charge.Total, tt.wantAdjustment)
			}
		})
	}
}

func TestBillingService_RegularCharge_NoAdjustment(t *testing.T) {
	calc := NewCalculationService()
	b := NewBillingService(calc)

	studio := domain.Studio{ID: 1, Rent: 1000}
	tenant := domain.Tenant{ID: 1, StudioID: 1, MoveIn: date(2026, time.March, 20)}

	// First steady-state period after the May 25 realignment.
	charge := b.RegularCharge(studio, tenant, 200, date(2026, time.June, 1))

	if !charge.PeriodStart.Equal(date(2026, time.May, 25)) {
		t.Errorf("PeriodStart = %v, want 2026-05-25", charge.PeriodStart)
	}
	if !charge.PeriodEnd.Equal(date(2026, time.June, 25)) {
		t.Errorf("PeriodEnd = %v, want 2026-06-25", charge.PeriodEnd)
	}
	if charge.Adjustment != 0 {
		t.Errorf("Adjustment = %v, want 0 (steady state)", charge.Adjustment)
	}
	if charge.Total != 1200 {
		t.Errorf("Total = %v, want 1200", charge.Total)
	}
}

// End-to-end: walk through the full sequence for the 20-March example and
// verify every payment lines up with what the spec describes.
func TestBillingService_FullSequence_MoveInMar20(t *testing.T) {
	calc := NewCalculationService()
	b := NewBillingService(calc)

	studio := domain.Studio{ID: 1, Rent: 3000}
	tenant := domain.Tenant{ID: 1, StudioID: 1, MoveIn: date(2026, time.March, 20)}

	moveIn := b.MoveInCharge(studio, tenant)
	if moveIn.Total != 6000 {
		t.Errorf("move-in Total = %v, want 6000 (rent + deposit)", moveIn.Total)
	}

	renewal := b.AdvanceRenewalCharge(studio, tenant, 400)
	if !renewal.PeriodStart.Equal(date(2026, time.April, 19)) {
		t.Errorf("renewal start = %v, want 2026-04-19", renewal.PeriodStart)
	}
	if renewal.Total != 3400 {
		t.Errorf("renewal Total = %v, want 3400", renewal.Total)
	}

	realign := b.RealignmentCharge(studio, tenant)
	if !realign.PeriodEnd.Equal(date(2026, time.May, 25)) {
		t.Errorf("realign end = %v, want 2026-05-25", realign.PeriodEnd)
	}
	want := 6 * (3000.0 / 30.0)
	if !almostEqual(realign.Total, want) {
		t.Errorf("realign Total = %v, want %v (6 days × 100)", realign.Total, want)
	}

	regular := b.RegularCharge(studio, tenant, 400, date(2026, time.June, 25))
	if !regular.PeriodStart.Equal(date(2026, time.June, 25)) {
		t.Errorf("regular start = %v, want 2026-06-25", regular.PeriodStart)
	}
	if regular.Total != 3400 {
		t.Errorf("regular Total = %v, want 3400", regular.Total)
	}
}

func TestBillingService_IssuedAt_RoughlyNow(t *testing.T) {
	calc := NewCalculationService()
	b := NewBillingService(calc)

	studio := domain.Studio{ID: 1, Rent: 100}
	tenant := domain.Tenant{ID: 1, StudioID: 1, MoveIn: date(2026, time.March, 25)}

	before := time.Now()
	charge := b.MoveInCharge(studio, tenant)
	after := time.Now()

	if charge.IssuedAt.Before(before) || charge.IssuedAt.After(after) {
		t.Errorf("IssuedAt = %v not in [%v, %v]", charge.IssuedAt, before, after)
	}
}
