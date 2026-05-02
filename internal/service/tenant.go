package service

import (
	"context"
	"errors"
	"time"

	"ChipHouse/internal/domain"
	"ChipHouse/internal/repository"
)

type TenantService struct {
	repo *repository.TenantRepo
}

func NewTenantService(repo *repository.TenantRepo) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) MoveIn(ctx context.Context, t *domain.Tenant) error {
	existing, _ := s.repo.GetActiveByStudioID(ctx, t.StudioID)
	if existing != nil {
		return errors.New("studio already has an active tenant")
	}
	if t.Persons < 1 {
		t.Persons = 1
	}
	t.IsActive = true
	if t.MoveIn.IsZero() {
		t.MoveIn = time.Now()
	}
	return s.repo.Create(ctx, t)
}

func (s *TenantService) MoveOut(ctx context.Context, tenantID int, moveOut time.Time) error {
	return s.repo.Deactivate(ctx, tenantID, moveOut)
}
