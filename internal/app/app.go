package app

import (
	"ChipHouse/config"
	"ChipHouse/internal/repository"
	"ChipHouse/internal/service"
)

type App struct {
	Cfg *config.Config
	DB  *repository.DB

	StudioRepo *repository.StudioRepo
	TenantRepo *repository.TenantRepo
	UserRepo   *repository.UserRepo
	ChargeRepo *repository.ChargeRepo

	CalculationService *service.CalculationService
	UtilityService     *service.UtilityService
	BillingService     *service.BillingService
	TenantService      *service.TenantService
	AuthService        *service.AuthService
}

func NewApp() (*App, error) {
	cfg := config.New()

	pool, err := repository.NewDB(cfg)
	if err != nil {
		return nil, err
	}
	db := &repository.DB{Pool: pool}

	studioRepo := repository.NewStudioRepo(pool)
	tenantRepo := repository.NewTenantRepo(pool)
	userRepo := repository.NewUserRepo(pool)
	chargeRepo := repository.NewChargeRepo(pool)

	calc := service.NewCalculationService()

	return &App{
		Cfg:                cfg,
		DB:                 db,
		StudioRepo:         studioRepo,
		TenantRepo:         tenantRepo,
		UserRepo:           userRepo,
		ChargeRepo:         chargeRepo,
		CalculationService: calc,
		UtilityService:     service.NewUtilityService(),
		BillingService:     service.NewBillingService(calc),
		TenantService:      service.NewTenantService(tenantRepo),
		AuthService:        service.NewAuthService(userRepo),
	}, nil
}
