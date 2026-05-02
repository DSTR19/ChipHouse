package service

import "ChipHouse/internal/domain"

type UtilityService struct{}

func NewUtilityService() *UtilityService {
	return &UtilityService{}
}

// SplitByArea splits a utility cost across studios proportional to their area.
func (u *UtilityService) SplitByArea(totalCost float64, studios []domain.Studio, target domain.Studio) float64 {
	var totalArea float64
	for _, s := range studios {
		totalArea += s.Area
	}
	if totalArea == 0 {
		return 0
	}
	return totalCost * (target.Area / totalArea)
}

// SplitByPersons splits a utility cost proportional to number of persons in each studio.
func (u *UtilityService) SplitByPersons(totalCost float64, allTenants []domain.Tenant, target domain.Tenant) float64 {
	totalPersons := 0
	for _, t := range allTenants {
		if t.IsActive {
			totalPersons += t.Persons
		}
	}
	if totalPersons == 0 {
		return 0
	}
	return totalCost * (float64(target.Persons) / float64(totalPersons))
}
