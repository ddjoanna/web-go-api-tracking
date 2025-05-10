package service

import (
	"context"
	"time"
	datastructure "tracking-service/internal/datastructures"
	errdefs "tracking-service/internal/errors"
	model "tracking-service/internal/models"
	repository "tracking-service/internal/repositories"

	"github.com/bwmarrin/snowflake"
)

type TenantService struct {
	snowflake *snowflake.Node
	repo      repository.TenantRepository
}

func NewTenantService(
	snowflake *snowflake.Node,
	repo repository.TenantRepository,
) *TenantService {
	return &TenantService{
		snowflake: snowflake,
		repo:      repo,
	}
}

func (s *TenantService) CreateTenant(ctx context.Context, in *datastructure.Tenant) (*model.Tenant, error) {
	tenant := &model.Tenant{
		ID:          s.snowflake.Generate().String(),
		Name:        in.Name,
		Description: in.Description,
		CreatedAt:   time.Now(),
	}

	if err := s.repo.CreateTenant(ctx, tenant); err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	return tenant, nil
}

func (s *TenantService) GetTenant(ctx context.Context, id string) (*model.Tenant, error) {
	tenant, err := s.repo.GetTenantByID(ctx, id)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}
	return tenant, nil
}

func (s *TenantService) UpdateTenant(ctx context.Context, id string, in *datastructure.Tenant) error {
	tenant, err := s.repo.GetTenantByID(ctx, id)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	tenant.Name = in.Name
	tenant.Description = in.Description
	tenant.UpdatedAt = time.Now()

	return s.repo.UpdateTenant(ctx, tenant)
}

func (s *TenantService) GetTenants(ctx context.Context) ([]*model.Tenant, error) {
	return s.repo.GetTenants(ctx)
}
