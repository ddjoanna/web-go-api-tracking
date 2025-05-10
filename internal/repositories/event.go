package repository

import (
	"context"
	model "tracking-service/internal/models"

	"gorm.io/gorm"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, event *model.Event) error
	GetEventByID(ctx context.Context, id string) (*model.Event, error)
	UpdateEvent(ctx context.Context, event *model.Event) error
	GetEvents(ctx context.Context) ([]*model.Event, error)
	DeleteEvent(ctx context.Context, event *model.Event) error
	CreateEventField(ctx context.Context, eventField *model.EventField) error
	GetEventFieldByEventIDAndID(ctx context.Context, eventID string, fieldID string) (*model.EventField, error)
	UpdateEventField(ctx context.Context, eventField *model.EventField) error
	DeleteEventField(ctx context.Context, eventField *model.EventField) error
	GetEventFields(ctx context.Context, eventID string) ([]*model.EventField, error)
	GetEventsByApplicationID(ctx context.Context, applicationID string) ([]*model.Event, error)
	GetEventByApplicationIDAndID(ctx context.Context, applicationID string, id string) (*model.Event, error)
	CreateEventLog(ctx context.Context, eventLog *model.EventLog) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) CreateEvent(ctx context.Context, event *model.Event) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(event).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *eventRepository) GetEventByID(ctx context.Context, id string) (*model.Event, error) {
	var event model.Event
	err := r.db.WithContext(ctx).
		Preload("Fields").
		First(&event, "id = ?", id).Error
	return &event, err
}

func (r *eventRepository) UpdateEvent(ctx context.Context, event *model.Event) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(event).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *eventRepository) DeleteEvent(ctx context.Context, event *model.Event) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(event).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *eventRepository) GetEvents(ctx context.Context) ([]*model.Event, error) {
	var events []*model.Event
	err := r.db.WithContext(ctx).Find(&events).Error
	return events, err
}

func (r *eventRepository) CreateEventField(ctx context.Context, eventField *model.EventField) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(eventField).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *eventRepository) UpdateEventField(ctx context.Context, eventField *model.EventField) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(eventField).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *eventRepository) DeleteEventField(ctx context.Context, eventField *model.EventField) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(eventField).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *eventRepository) GetEventFields(ctx context.Context, eventID string) ([]*model.EventField, error) {
	var eventFields []*model.EventField
	err := r.db.WithContext(ctx).Where("event_id = ?", eventID).Find(&eventFields).Error
	return eventFields, err
}

func (r *eventRepository) GetEventByApplicationID(ctx context.Context, id string) (*model.Event, error) {
	var event model.Event
	err := r.db.WithContext(ctx).First(&event, "application_id = ?", id).Error
	return &event, err
}

func (r *eventRepository) GetEventsByApplicationID(ctx context.Context, applicationID string) ([]*model.Event, error) {
	var events []*model.Event
	err := r.db.WithContext(ctx).
		Preload("Fields").
		Where("application_id = ?", applicationID).Find(&events).Error
	return events, err
}

func (r *eventRepository) GetEventByApplicationIDAndID(ctx context.Context, applicationID string, id string) (*model.Event, error) {
	var event model.Event
	err := r.db.WithContext(ctx).
		Preload("Fields").
		First(&event, "application_id = ? AND id = ?", applicationID, id).Error
	return &event, err
}

func (r *eventRepository) CreateEventLog(ctx context.Context, eventLog *model.EventLog) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(eventLog).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *eventRepository) GetEventFieldByEventIDAndID(ctx context.Context, eventID string, fieldID string) (*model.EventField, error) {
	var eventField model.EventField
	err := r.db.WithContext(ctx).First(&eventField, "event_id = ? AND id = ?", eventID, fieldID).Error
	return &eventField, err
}
