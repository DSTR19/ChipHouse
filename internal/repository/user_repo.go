package repository

import (
	"context"

	"ChipHouse/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u domain.User
	err := r.db.QueryRow(ctx,
		"SELECT id, email, password_hash, role, tenant_id FROM users WHERE email=$1",
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role, &u.TenantID)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
