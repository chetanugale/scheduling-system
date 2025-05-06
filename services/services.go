package services

import (
	"context"

	"github.com/chetanugale/scheduling-system/models"
	"github.com/chetanugale/scheduling-system/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventService interface {
	CreateEvent(ctx context.Context, e models.Event) (*models.Event, error)
	GetEvent(ctx context.Context, id string) (*models.Event, error)
	UpdateEvent(ctx context.Context, id string, update models.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetAllEvents(ctx context.Context, title string) ([]models.Event, error)
}

type MongoEventService struct {
	Repo repository.MongoRepository[models.Event]
}

func (s *MongoEventService) CreateEvent(ctx context.Context, e models.Event) (*models.Event, error) {
	return s.Repo.Insert(ctx, e)
}

func (s *MongoEventService) GetEvent(ctx context.Context, id string) (*models.Event, error) {
	return s.Repo.GetByID(ctx, id)
}

func (s *MongoEventService) UpdateEvent(ctx context.Context, id string, update models.Event) error {
	return s.Repo.UpdateByID(ctx, id, update)
}

func (s *MongoEventService) DeleteEvent(ctx context.Context, id string) error {
	return s.Repo.DeleteByID(ctx, id)
}

func (s *MongoEventService) GetAllEvents(ctx context.Context, title string) ([]models.Event, error) {
	filter := bson.M{}
	if title != "" {
		filter = bson.M{"title": bson.M{"$eq": title}}
	}
	return s.Repo.FindAll(ctx, filter)
}

// ---------------- Availability ----------------

type AvailabilityService interface {
	AddAvailability(ctx context.Context, a models.Availability) (*models.Availability, error)
	GetAvailabilitiesByEvent(ctx context.Context, eventID string) ([]models.Availability, error)
	DeleteAvailability(ctx context.Context, id string) error
	UpdateAvailability(ctx context.Context, id string, a models.Availability) error
}

type MongoAvailabilityService struct {
	Repo repository.MongoRepository[models.Availability]
}

func (s *MongoAvailabilityService) AddAvailability(ctx context.Context, a models.Availability) (*models.Availability, error) {
	return s.Repo.Insert(ctx, a)
}

func (s *MongoAvailabilityService) GetAvailabilitiesByEvent(ctx context.Context, eventID string) ([]models.Availability, error) {
	eid, err := primitive.ObjectIDFromHex(eventID)
	if err != nil {
		return nil, err
	}
	return s.Repo.FindAll(ctx, bson.M{"eventid": eid})
}

func (s *MongoAvailabilityService) DeleteAvailability(ctx context.Context, id string) error {
	return s.Repo.DeleteByID(ctx, id)
}

func (s *MongoAvailabilityService) UpdateAvailability(ctx context.Context, id string, a models.Availability) error {
	return s.Repo.UpdateByID(ctx, id, a)
}
