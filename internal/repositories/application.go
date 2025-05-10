package repository

import (
	"context"
	model "tracking-service/internal/models"

	"gorm.io/gorm"
)

type ApplicationRepository interface {
	CreateApplication(ctx context.Context, application *model.Application) error
	GetApplicationByID(ctx context.Context, id string) (*model.Application, error)
	UpdateApplication(ctx context.Context, application *model.Application) error
	DeleteApplication(ctx context.Context, application *model.Application) error
	GetApplications(ctx context.Context) ([]*model.Application, error)
	GetApplicationByAPIKey(ctx context.Context, apiKey string) (*model.Application, error)
	CreateSession(ctx context.Context, session *model.Session) error
	GetSessionByApplicationIDAndID(ctx context.Context, applicationID string, id string) (*model.Session, error)
	UpdateSession(ctx context.Context, session *model.Session) error
	DeleteSession(ctx context.Context, session *model.Session) error
	GetApplicationByTenantIDAndID(ctx context.Context, tenantID string, id string) (*model.Application, error)
	CreateApplicationAPIKey(ctx context.Context, application *model.ApplicationApiKey) error
	GetApplicationKeyByID(ctx context.Context, apiKeyID string) (*model.ApplicationApiKey, error)
	DeleteApplicationAPIKey(ctx context.Context, apiKey *model.ApplicationApiKey) error
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{
		db: db,
	}
}

func (r *applicationRepository) CreateApplication(ctx context.Context, application *model.Application) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(application).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *applicationRepository) GetApplicationByID(ctx context.Context, id string) (*model.Application, error) {
	var application model.Application
	err := r.db.WithContext(ctx).
		Preload("ApiKeys").
		First(&application, "id = ?", id).Error
	return &application, err
}

func (r *applicationRepository) UpdateApplication(ctx context.Context, application *model.Application) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(application).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *applicationRepository) DeleteApplication(ctx context.Context, application *model.Application) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(application).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *applicationRepository) GetApplications(ctx context.Context) ([]*model.Application, error) {
	var applications []*model.Application
	err := r.db.WithContext(ctx).Find(&applications).Error
	return applications, err
}

func (r *applicationRepository) GetApplicationByAPIKey(ctx context.Context, apiKey string) (*model.Application, error) {
	var applicationApiKey model.ApplicationApiKey
	err := r.db.WithContext(ctx).
		First(&applicationApiKey, "api_key = ?", apiKey).Error

	if err != nil {
		return nil, err
	}

	var application model.Application
	err = r.db.WithContext(ctx).
		First(&application, "id = ?", applicationApiKey.ApplicationID).Error

	if err != nil {
		return nil, err
	}

	return &application, nil
}

func (r *applicationRepository) CreateSession(ctx context.Context, session *model.Session) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(session).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *applicationRepository) GetSessionByApplicationIDAndID(ctx context.Context, applicationID string, id string) (*model.Session, error) {
	var session model.Session
	err := r.db.WithContext(ctx).First(&session, "application_id = ? AND id = ?", applicationID, id).Error
	return &session, err
}

func (r *applicationRepository) UpdateSession(ctx context.Context, session *model.Session) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(session).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *applicationRepository) DeleteSession(ctx context.Context, session *model.Session) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(session).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *applicationRepository) GetApplicationByTenantIDAndID(ctx context.Context, tenantID string, id string) (*model.Application, error) {
	var application model.Application
	err := r.db.WithContext(ctx).First(&application, "tenant_id = ? AND id = ?", tenantID, id).Error
	return &application, err
}

func (r *applicationRepository) CreateApplicationAPIKey(ctx context.Context, application *model.ApplicationApiKey) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(application).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *applicationRepository) GetApplicationKeyByID(ctx context.Context, apiKeyID string) (*model.ApplicationApiKey, error) {
	var applicationApiKey model.ApplicationApiKey
	err := r.db.WithContext(ctx).First(&applicationApiKey, "id = ?", apiKeyID).Error
	return &applicationApiKey, err
}

func (r *applicationRepository) DeleteApplicationAPIKey(ctx context.Context, apiKey *model.ApplicationApiKey) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(apiKey).Error; err != nil {
			return err
		}
		return nil
	})
}
