package service

import (
	"context"
	"time"
	datastructure "tracking-service/internal/datastructures"
	errdefs "tracking-service/internal/errors"
	model "tracking-service/internal/models"
	repository "tracking-service/internal/repositories"
	util "tracking-service/internal/utils"

	"github.com/bwmarrin/snowflake"
)

type ApplicationService struct {
	snowflake     *snowflake.Node
	repo          repository.ApplicationRepository
	tenant_repo   repository.TenantRepository
	platform_repo repository.PlatformRepository
}

func NewApplicationService(
	snowflake *snowflake.Node,
	repo repository.ApplicationRepository,
	tantent_repo repository.TenantRepository,
	platform_repo repository.PlatformRepository,
) *ApplicationService {
	return &ApplicationService{
		snowflake:     snowflake,
		repo:          repo,
		tenant_repo:   tantent_repo,
		platform_repo: platform_repo,
	}
}

func (s *ApplicationService) CreateApplication(ctx context.Context, in *datastructure.Application) (*model.Application, error) {
	tenant, err := s.tenant_repo.GetTenantByID(ctx, in.TenantID)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	applicationID := s.snowflake.Generate().String()
	now := time.Now()
	apiKeyID := s.snowflake.Generate().String()
	apiKey, err := util.GenerateAPIKey(32)
	if err != nil {
		return nil, errdefs.ErrorInternalError
	}

	application := &model.Application{
		ID:          applicationID,
		TenantID:    tenant.ID,
		Name:        in.Name,
		Description: in.Description,
		CreatedAt:   now,
		ApiKeys: []model.ApplicationApiKey{
			{
				ID:            apiKeyID,
				ApplicationID: applicationID,
				APIKey:        apiKey,
				CreatedAt:     now,
			},
		},
	}

	if err := s.repo.CreateApplication(ctx, application); err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	return application, nil
}

func (s *ApplicationService) GetApplicationByID(ctx context.Context, id string) (*model.Application, error) {
	application, err := s.repo.GetApplicationByID(ctx, id)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}
	return application, nil
}

func (s *ApplicationService) UpdateApplicationByID(ctx context.Context, id string, in *datastructure.Application) error {
	application, err := s.repo.GetApplicationByID(ctx, id)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	application.Name = in.Name
	application.Description = in.Description
	application.UpdatedAt = time.Now()

	return s.repo.UpdateApplication(ctx, application)
}

func (s *ApplicationService) DeleteApplicationByID(ctx context.Context, id string) error {
	application, err := s.repo.GetApplicationByID(ctx, id)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	return s.repo.DeleteApplication(ctx, application)
}

func (s *ApplicationService) GetApplications(ctx context.Context) ([]*model.Application, error) {
	apps, err := s.repo.GetApplications(ctx)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}
	return apps, nil
}

func (s *ApplicationService) ValidateAPIKey(ctx context.Context, apiKey string) (*model.Application, error) {
	application, err := s.repo.GetApplicationByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	return application, nil
}

func (s *ApplicationService) CreateSession(ctx context.Context, in *datastructure.Session) (*model.Session, error) {
	application, err := s.repo.GetApplicationByID(ctx, in.ApplicationID)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	platform, err := s.platform_repo.GetPlatformByID(ctx, in.PlatformID)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	startedAt, err := util.ParseTimeDefaultFormat(in.StartedAt)
	if err != nil {
		return nil, errdefs.ErrorInvalidRequest
	}

	var endedAt *time.Time
	if in.EndedAt != nil {
		t, err := util.ParseTimeDefaultFormat(*in.EndedAt)
		if err != nil {
			return nil, errdefs.ErrorInvalidRequest
		}
		endedAt = &t
	}
	session := &model.Session{
		ID:            s.snowflake.Generate().String(),
		ApplicationID: application.ID,
		PlatformID:    platform.ID,
		SessionKey:    in.SessionKey,
		UserID:        in.UserID,
		UserAgent:     in.UserAgent,
		IPAddress:     in.IPAddress,
		StartedAt:     startedAt,
		EndedAt:       endedAt,
		CreatedAt:     time.Now(),
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	return session, nil
}

func (s *ApplicationService) GetSessionByApplicationIDAndID(ctx context.Context, applicationID string, id string) (*model.Session, error) {
	session, err := s.repo.GetSessionByApplicationIDAndID(ctx, applicationID, id)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}
	return session, nil
}

func (s *ApplicationService) UpdateSessionByApplicationIDAndID(ctx context.Context, applicationID string, id string, in *datastructure.UpdateSessionRequest) error {
	session, err := s.repo.GetSessionByApplicationIDAndID(ctx, applicationID, id)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	if in.UserID != "" {
		session.UserID = &in.UserID
	}

	if in.EndedAt != "" {
		t, err := util.ParseTimeDefaultFormat(in.EndedAt)
		if err != nil {
			return err
		}
		session.EndedAt = &t
	}

	session.UpdatedAt = time.Now()

	return s.repo.UpdateSession(ctx, session)
}

func (s *ApplicationService) GetApplicationByTenantIDAndID(ctx context.Context, tenantID string, id string) (*model.Application, error) {
	application, err := s.repo.GetApplicationByTenantIDAndID(ctx, tenantID, id)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}
	return application, nil
}

func (s *ApplicationService) CreateApplicationAPIKey(ctx context.Context, applicationID string) (*model.ApplicationApiKey, error) {
	application, err := s.repo.GetApplicationByID(ctx, applicationID)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	key, err := util.GenerateAPIKey(32)
	if err != nil {
		return nil, errdefs.ErrorInternalError
	}
	apiKey := &model.ApplicationApiKey{
		ID:            s.snowflake.Generate().String(),
		ApplicationID: application.ID,
		APIKey:        key,
		CreatedAt:     time.Now(),
	}

	err = s.repo.CreateApplicationAPIKey(ctx, apiKey)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}
	return apiKey, nil
}

func (s *ApplicationService) DeleteApplicationAPIKey(ctx context.Context, applicationID string, apiKeyID string) error {
	apiKey, err := s.repo.GetApplicationKeyByID(ctx, apiKeyID)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	if apiKey.ApplicationID != applicationID {
		return errdefs.ErrorInvalidRequest
	}

	return s.repo.DeleteApplicationAPIKey(ctx, apiKey)
}

func (s *ApplicationService) DeleteSessionByApplicationIDAndID(ctx context.Context, applicationID string, id string) error {
	session, err := s.repo.GetSessionByApplicationIDAndID(ctx, applicationID, id)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	return s.repo.DeleteSession(ctx, session)
}
