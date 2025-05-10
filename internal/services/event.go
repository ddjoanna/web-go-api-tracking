package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	shared "tracking-service/internal"
	datastructure "tracking-service/internal/datastructures"
	errdefs "tracking-service/internal/errors"
	model "tracking-service/internal/models"
	repository "tracking-service/internal/repositories"

	util "tracking-service/internal/utils"

	"github.com/IBM/sarama"
	"github.com/bwmarrin/snowflake"
	log "github.com/sirupsen/logrus"
)

type EventService struct {
	snowflake     *snowflake.Node
	repo          repository.EventRepository
	app_repo      repository.ApplicationRepository
	platform_repo repository.PlatformRepository
	event_repo    repository.EventRepository
	producer      sarama.SyncProducer
}

func NewEventService(
	snowflake *snowflake.Node,
	repo repository.EventRepository,
	app_repo repository.ApplicationRepository,
	platform_repo repository.PlatformRepository,
	event_repo repository.EventRepository,
	producer sarama.SyncProducer,
) *EventService {
	return &EventService{
		snowflake:     snowflake,
		repo:          repo,
		app_repo:      app_repo,
		platform_repo: platform_repo,
		event_repo:    event_repo,
		producer:      producer,
	}
}

func (s *EventService) CreateEvent(ctx context.Context, in *datastructure.Event) (*model.Event, error) {
	application, err := s.app_repo.GetApplicationByID(ctx, in.ApplicationID)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	platform, err := s.platform_repo.GetPlatformByID(ctx, in.PlatformID)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	event := &model.Event{
		ID:            s.snowflake.Generate().String(),
		ApplicationID: application.ID,
		PlatformID:    platform.ID,
		Name:          in.Name,
		Description:   in.Description,
		IsActive:      in.IsActive,
		CreatedAt:     time.Now(),
	}

	if err := s.event_repo.CreateEvent(ctx, event); err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	return event, nil
}

func (s *EventService) GetEventByID(ctx context.Context, id string) (*model.Event, error) {
	event, err := s.repo.GetEventByID(ctx, id)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}
	return event, nil
}

func (s *EventService) UpdateEventByID(ctx context.Context, id string, in *datastructure.Event) error {
	event, err := s.repo.GetEventByID(ctx, id)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	event.Name = in.Name
	event.Description = in.Description
	event.IsActive = in.IsActive
	event.UpdatedAt = time.Now()

	return s.repo.UpdateEvent(ctx, event)
}

func (s *EventService) DeleteEventByID(ctx context.Context, id string) error {
	event, err := s.repo.GetEventByID(ctx, id)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	return s.repo.DeleteEvent(ctx, event)
}

func (s *EventService) GetEvents(ctx context.Context) ([]*model.Event, error) {
	events, err := s.repo.GetEvents(ctx)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}
	return events, nil
}

func (s *EventService) CreateEventField(ctx context.Context, in *datastructure.EventField) (*model.EventField, error) {
	event, err := s.repo.GetEventByID(ctx, in.EventID)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	eventField := &model.EventField{
		ID:          s.snowflake.Generate().String(),
		EventID:     event.ID,
		Name:        in.Name,
		DataType:    in.DataType,
		IsRequired:  in.IsRequired,
		Description: in.Description,
		CreatedAt:   time.Now(),
	}

	if err := s.event_repo.CreateEventField(ctx, eventField); err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	return eventField, nil
}

func (s *EventService) GetEventFieldByEventIDAndID(ctx context.Context, eventID string, fieldID string) (*model.EventField, error) {
	field, err := s.repo.GetEventFieldByEventIDAndID(ctx, eventID, fieldID)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}
	return field, nil
}

func (s *EventService) UpdateEventFieldByEventIDAndID(ctx context.Context, eventID string, fieldID string, in *datastructure.EventField) error {
	field, err := s.repo.GetEventFieldByEventIDAndID(ctx, eventID, fieldID)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	field.Name = in.Name
	field.DataType = in.DataType
	field.Description = in.Description
	field.IsRequired = in.IsRequired
	field.UpdatedAt = time.Now()

	return s.repo.UpdateEventField(ctx, field)
}

func (s *EventService) DeleteEventFieldByEventIDAndID(ctx context.Context, eventID string, fieldID string) error {
	field, err := s.repo.GetEventFieldByEventIDAndID(ctx, eventID, fieldID)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	return s.repo.DeleteEventField(ctx, field)
}

func (s *EventService) GetEventFields(ctx context.Context, eventID string) ([]*model.EventField, error) {
	eventFields, err := s.repo.GetEventFields(ctx, eventID)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}
	return eventFields, nil
}

func (s *EventService) GetEventByTenant(ctx context.Context, applicationID string, eventID string) (*model.Event, error) {
	event, err := s.repo.GetEventByApplicationIDAndID(ctx, applicationID, eventID)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}
	return event, nil
}

func (s *EventService) UpdateEventByTenant(ctx context.Context, applicationID string, in *datastructure.Event) error {
	event, err := s.repo.GetEventByApplicationIDAndID(ctx, applicationID, in.ID)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	event.Name = in.Name
	event.Description = in.Description
	event.IsActive = in.IsActive
	event.UpdatedAt = time.Now()

	return s.repo.UpdateEvent(ctx, event)
}

func (s *EventService) DeleteEventByTenant(ctx context.Context, applicationID string, eventID string) error {
	event, err := s.repo.GetEventByApplicationIDAndID(ctx, applicationID, eventID)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	return s.repo.DeleteEvent(ctx, event)
}

func (s *EventService) GetEventFieldByTenant(
	ctx context.Context,
	applicationID string,
	eventID string,
	fieldID string,
) (*model.EventField, error) {
	event, err := s.repo.GetEventByApplicationIDAndID(ctx, applicationID, eventID)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	field, err := s.repo.GetEventFieldByEventIDAndID(ctx, event.ID, fieldID)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	return field, nil
}

func (s *EventService) UpdateEventFieldByTenant(
	ctx context.Context,
	applicationID string,
	eventID string,
	fieldID string,
	in *datastructure.EventField,
) error {
	event, err := s.repo.GetEventByApplicationIDAndID(ctx, applicationID, eventID)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	field, err := s.repo.GetEventFieldByEventIDAndID(ctx, event.ID, fieldID)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	field.Name = in.Name
	field.DataType = in.DataType
	field.Description = in.Description
	field.IsRequired = in.IsRequired
	field.UpdatedAt = time.Now()

	return s.repo.UpdateEventField(ctx, field)
}

func (s *EventService) DeleteEventFieldByTenant(
	ctx context.Context,
	applicationID string,
	eventID string,
	fieldID string,
) error {
	event, err := s.repo.GetEventByApplicationIDAndID(ctx, applicationID, eventID)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	field, err := s.repo.GetEventFieldByEventIDAndID(ctx, event.ID, fieldID)
	if err != nil {
		return errdefs.WrapGormError(err)
	}

	return s.repo.DeleteEventField(ctx, field)
}

func (s *EventService) GetEventsByApplicationID(ctx context.Context, applicationID string) ([]*model.Event, error) {
	events, err := s.repo.GetEventsByApplicationID(ctx, applicationID)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}
	return events, nil
}

func (s *EventService) CreateEventLog(ctx context.Context, in *datastructure.EventLog) (*model.EventLog, error) {
	event, err := s.repo.GetEventByID(ctx, in.EventID)
	if err != nil {
		return nil, errdefs.WrapGormError(err)
	}

	eventLog := &model.EventLog{
		ID:            s.snowflake.Generate().String(),
		ApplicationID: in.ApplicationID,
		SessionID:     in.SessionID,
		EventID:       event.ID,
		PlatformID:    in.PlatformID,
		Properties:    in.Properties,
		CreatedAt:     time.Now(),
	}

	// 傳送 kafka 資料，若失敗則降級直接儲存資料；異常資料同步由NiFi執行
	msg, err := s.createKafkaMessage(eventLog)
	if err != nil {
		if err := s.repo.CreateEventLog(ctx, eventLog); err != nil {
			return nil, errdefs.WrapGormError(err)
		}

		return eventLog, nil
	}

	err = util.WithRetry(ctx, 3, func() error {
		_, _, err := s.producer.SendMessage(msg)
		return err
	})

	if err != nil {
		log.WithContext(ctx).Errorf("Failed to send message to kafka: %v", err)
		if err := s.repo.CreateEventLog(ctx, eventLog); err != nil {
			log.WithContext(ctx).Errorf("Failed to create event log: %v", err)
			return nil, errdefs.WrapGormError(err)
		}

		return eventLog, nil
	}

	log.WithContext(ctx).Infof("Successfully sent message to kafka: %v", msg)
	return eventLog, nil
}

func (s *EventService) createKafkaMessage(
	queue *model.EventLog,
) (*sarama.ProducerMessage, error) {
	jsonData, err := json.Marshal(queue)
	if err != nil {
		return nil, fmt.Errorf("marshal queue failed: %w", err)
	}

	driverTraceId := s.snowflake.Generate().String()
	msg := &sarama.ProducerMessage{
		Topic: shared.KafkaTopic,
		Key:   sarama.StringEncoder(driverTraceId),
		Value: sarama.ByteEncoder(jsonData),
	}

	return msg, nil
}
