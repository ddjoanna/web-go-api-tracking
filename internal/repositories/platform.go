package repository

import (
	"context"
	model "tracking-service/internal/models"

	"gorm.io/gorm"
)

type PlatformRepository interface {
	CreatePlatform(ctx context.Context, platform *model.Platform) error
	GetPlatformByID(ctx context.Context, id int) (*model.Platform, error)
	GetPlatforms(ctx context.Context) ([]*model.Platform, error)
}

type platformRepository struct {
	db *gorm.DB
}

func NewPlatformRepository(db *gorm.DB) PlatformRepository {
	return &platformRepository{
		db: db,
	}
}

func (r *platformRepository) CreatePlatform(ctx context.Context, platform *model.Platform) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(platform).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *platformRepository) GetPlatformByID(ctx context.Context, id int) (*model.Platform, error) {
	var platform model.Platform
	err := r.db.WithContext(ctx).First(&platform, "id = ?", id).Error
	return &platform, err
}

func (r *platformRepository) GetPlatforms(ctx context.Context) ([]*model.Platform, error) {
	var platforms []*model.Platform
	err := r.db.WithContext(ctx).Find(&platforms).Error
	return platforms, err
}
