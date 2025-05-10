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

type PlatformService struct {
	snowflake *snowflake.Node
	repo      repository.PlatformRepository
}

func NewPlatformService(
	snowflake *snowflake.Node,
	repo repository.PlatformRepository,
) *PlatformService {
	return &PlatformService{
		snowflake: snowflake,
		repo:      repo,
	}
}

func (s *PlatformService) CreatePlatform(ctx context.Context, in *datastructure.Platform) (*model.Platform, error) {
	platform := &model.Platform{
		Name:      in.Name,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreatePlatform(ctx, platform); err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	return platform, nil
}

func (s *PlatformService) GetPlatformByID(ctx context.Context, id int) (*model.Platform, error) {
	platform, err := s.repo.GetPlatformByID(ctx, id)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}
	return platform, nil
}

func (s *PlatformService) GetPlatforms(ctx context.Context) ([]*model.Platform, error) {
	platforms, err := s.repo.GetPlatforms(ctx)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}
	return platforms, nil
}
