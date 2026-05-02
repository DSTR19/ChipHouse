package repository

import (
	"context"

	"ChipHouse/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TenantRepo struct {
	db *pgxpool.Pool
}

func NewTenantRepo(db *pgxpool.Pool) *TenantRepo {
	return &TenantRepo{db: db}
}

func (r *TenantRepo) GetActiveByStudioID(ctx context.Context, studioID int) (*domain.Tenant, error) {
	row := r.db.QueryRow(ctx,
		"SELECT id, full_name, studio_id, move_in, move_out, persons, is_active FROM tenants WHERE studio_id=$1 AND is_active=true",
		studioID,
	)

	var t domain.Tenant
	if err := row.Scan(&t.ID, &t.FullName, &t.StudioID, &t.MoveIn, &t.MoveOut, &t.Persons, &t.IsActive); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TenantRepo) Create(ctx context.Context, t *domain.Tenant) error {
	return r.db.QueryRow(ctx,
		`INSERT INTO tenants (full_name, studio_id, move_in, persons, is_active)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		t.FullName, t.StudioID, t.MoveIn, t.Persons, t.IsActive,
	).Scan(&t.ID)
}

func (r *TenantRepo) Deactivate(ctx context.Context, tenantID int, moveOut interface{}) error {
	_, err := r.db.Exec(ctx,
		"UPDATE tenants SET is_active=false, move_out=$2 WHERE id=$1",
		tenantID, moveOut,
	)
	return err
}
