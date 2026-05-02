package repository

import (
	"context"

	"ChipHouse/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StudioRepo struct {
	db *pgxpool.Pool
}

func NewStudioRepo(db *pgxpool.Pool) *StudioRepo {
	return &StudioRepo{db: db}
}

func (r *StudioRepo) GetAll(ctx context.Context) ([]domain.Studio, error) {
	rows, err := r.db.Query(ctx, "SELECT id, number, area, rent FROM studios")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.Studio
	for rows.Next() {
		var s domain.Studio
		if err := rows.Scan(&s.ID, &s.Number, &s.Area, &s.Rent); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *StudioRepo) GetByID(ctx context.Context, id int) (*domain.Studio, error) {
	var s domain.Studio
	err := r.db.QueryRow(ctx,
		"SELECT id, number, area, rent FROM studios WHERE id=$1",
		id,
	).Scan(&s.ID, &s.Number, &s.Area, &s.Rent)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
