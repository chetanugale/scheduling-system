package mocker

import (
	"context"

	"github.com/chetanugale/scheduling-system/models"
	"github.com/stretchr/testify/mock"
)

type MockEventService struct {
	mock.Mock
}

type MockAvailabilityService struct {
	mock.Mock
}

func (m *MockEventService) CreateEvent(ctx context.Context, event models.Event) (*models.Event, error) {
	args := m.Called(ctx, event)
	return args.Get(0).(*models.Event), args.Error(1)
}

func (m *MockEventService) GetEvent(ctx context.Context, id string) (*models.Event, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Event), args.Error(1)
}
func (m *MockEventService) DeleteEvent(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockEventService) GetAllEvents(ctx context.Context, filter string) ([]models.Event, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Event), args.Error(1)
}
func (m *MockEventService) UpdateEvent(ctx context.Context, id string, event models.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}
func (m *MockAvailabilityService) AddAvailability(ctx context.Context, a models.Availability) (*models.Availability, error) {
	args := m.Called(ctx, a)
	return args.Get(0).(*models.Availability), args.Error(1)
}

func (m *MockAvailabilityService) GetAvailabilitiesByEvent(ctx context.Context, eventID string) ([]models.Availability, error) {
	args := m.Called(ctx, eventID)
	return args.Get(0).([]models.Availability), args.Error(1)
}
func (m *MockAvailabilityService) DeleteAvailability(ctx context.Context, eventID string) error {
	args := m.Called(ctx, eventID)
	return args.Error(0)
}
func (m *MockAvailabilityService) UpdateAvailability(ctx context.Context, id string, avail models.Availability) error {
	args := m.Called(ctx, avail)
	return args.Error(0)
}

type MockRepo[T any] struct {
	mock.Mock
}

func (m *MockRepo[T]) Insert(ctx context.Context, obj T) (*T, error) {
	args := m.Called(ctx, obj)
	return args.Get(0).(*T), args.Error(1)
}

func (m *MockRepo[T]) GetByID(ctx context.Context, id string) (*T, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*T), args.Error(1)
}
