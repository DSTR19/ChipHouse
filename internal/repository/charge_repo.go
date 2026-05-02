package repository

import (
	"ChipHouse/internal/domain"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChargeRepo struct {
	db *pgxpool.Pool
}

func NewChargeRepo(db *pgxpool.Pool) *ChargeRepo {
	return &ChargeRepo{db: db}
}

func (r *ChargeRepo) Create(ctx context.Context, c *domain.Charge) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO charges (studio_id, tenant_id, period_start, period_end, rent, utilities, adjustment, total, paid, issued_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		c.StudioID, c.TenantID, c.PeriodStart, c.PeriodEnd,
		c.Rent, c.Utilities, c.Adjustment, c.Total, c.Paid, c.IssuedAt,
	).Scan(&c.ID)
}

func (r *ChargeRepo) GetByTenantID(ctx context.Context, tenantID int) ([]domain.Charge, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, studio_id, tenant_id, period_start, period_end, rent, utilities, adjustment, total, paid, issued_at
		 FROM charges WHERE tenant_id=$1 ORDER BY period_start DESC`,
		tenantID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Charge
	for rows.Next() {
		var c domain.Charge
		if err := rows.Scan(&c.ID, &c.StudioID, &c.TenantID, &c.PeriodStart, &c.PeriodEnd,
			&c.Rent, &c.Utilities, &c.Adjustment, &c.Total, &c.Paid, &c.IssuedAt); err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}
