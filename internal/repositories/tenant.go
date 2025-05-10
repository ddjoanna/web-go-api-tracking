package repository

import (
	"context"
	model "tracking-service/internal/models"

	"gorm.io/gorm"
)

type TenantRepository interface {
	CreateTenant(ctx context.Context, tenant *model.Tenant) error
	GetTenantByID(ctx context.Context, id string) (*model.Tenant, error)
	UpdateTenant(ctx context.Context, tenant *model.Tenant) error
	GetTenants(ctx context.Context) ([]*model.Tenant, error)
}

type tenantRepository struct {
	db *gorm.DB
}

func NewTenantRepository(db *gorm.DB) TenantRepository {
	return &tenantRepository{
		db: db,
	}
}

func (r *tenantRepository) CreateTenant(ctx context.Context, tenant *model.Tenant) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(tenant).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *tenantRepository) GetTenantByID(ctx context.Context, id string) (*model.Tenant, error) {
	var tenant model.Tenant
	err := r.db.WithContext(ctx).First(&tenant, "id = ?", id).Error
	return &tenant, err
}

func (r *tenantRepository) UpdateTenant(ctx context.Context, tenant *model.Tenant) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(tenant).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *tenantRepository) GetTenants(ctx context.Context) ([]*model.Tenant, error) {
	var tenants []*model.Tenant
	err := r.db.WithContext(ctx).Find(&tenants).Error
	return tenants, err
}
