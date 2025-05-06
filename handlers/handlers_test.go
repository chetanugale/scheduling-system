package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/chetanugale/scheduling-system/mocker"
	"github.com/chetanugale/scheduling-system/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// --------- GET /events/:id/recommend -----------

func TestRecommendHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	slotID := primitive.NewObjectID()
	mockEventSvc := new(mocker.MockEventService)
	mockAvailSvc := new(mocker.MockAvailabilityService)

	event := &models.Event{
		ID:    primitive.NewObjectID(),
		Title: "Mock Event",
		Slots: []models.TimeSlot{{ID: slotID}},
	}

	mockEventSvc.On("GetEvent", mock.Anything, "abc123").Return(event, nil)
	mockAvailSvc.On("GetAvailabilitiesByEvent", mock.Anything, "abc123").
		Return([]models.Availability{{UserID: "u1", SlotID: slotID}}, nil)

	router := gin.New()
	router.GET("/events/:id/recommend", RecommendHandler(mockEventSvc, mockAvailSvc))

	req, _ := http.NewRequest(http.MethodGet, "/events/abc123/recommend", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockEventSvc.AssertExpectations(t)
	mockAvailSvc.AssertExpectations(t)
}

// --------- POST /events -----------
func TestCreateEventHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mocker.MockEventService)

	event := models.Event{Title: "Planning"}
	mockSvc.On("CreateEvent", mock.Anything, mock.AnythingOfType("Event")).
		Return(&event, nil)

	router := gin.New()
	router.POST("/events", CreateEventHandler(mockSvc))

	body, _ := json.Marshal(event)
	req, _ := http.NewRequest(http.MethodPost, "/events", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockSvc.AssertExpectations(t)
}

// --------- GET /events/:id -----------

func TestGetEventHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mocker.MockEventService)

	id := primitive.NewObjectID().Hex()
	event := &models.Event{Title: "Demo", ID: primitive.NewObjectID()}

	mockSvc.On("GetEvent", mock.Anything, id).Return(event, nil)

	router := gin.New()
	router.GET("/events/:id", GetEventHandler(mockSvc))

	req, _ := http.NewRequest(http.MethodGet, "/events/"+id, nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// --------- DELETE /events/:id -----------

func TestDeleteEventHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mocker.MockEventService)

	id := primitive.NewObjectID().Hex()

	mockSvc.On("DeleteEvent", mock.Anything, id).Return(nil)

	router := gin.New()
	router.DELETE("/events/:id", DeleteEventHandler(mockSvc))

	req, _ := http.NewRequest(http.MethodDelete, "/events/"+id, nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
}

// --------- UPDATE /events/:id -----------

func TestUpdateEventHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mocker.MockEventService)

	id := primitive.NewObjectID().Hex()
	event := &models.Event{Title: "Demo", ID: primitive.NewObjectID()}

	mockSvc.On("GetEvent", mock.Anything, id).Return(event, nil)
	mockSvc.On("UpdateEvent", mock.Anything, id).Return(event, nil)

	router := gin.New()
	router.PUT("/events/:id", UpdateEventHandler(mockSvc))

	body, _ := json.Marshal(event)
	req, _ := http.NewRequest(http.MethodPut, "/events"+id, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

// --------- GET /events -----------

func TestGetAllEventHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mocker.MockEventService)

	event := &models.Event{Title: "Demo", ID: primitive.NewObjectID()}
	eventList := &[]models.Event{*event}

	mockSvc.On("GetAllEvents", mock.Anything, mock.Anything).Return(*eventList, nil)

	router := gin.New()
	router.GET("/events", GetAllEventsHandler(mockSvc))

	req, _ := http.NewRequest(http.MethodGet, "/events", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// --------- POST /availability -----------
func TestCreateAvailability(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mocker.MockAvailabilityService)

	event := models.Availability{EventID: primitive.NewObjectID()}
	mockSvc.On("AddAvailability", mock.Anything, mock.AnythingOfType("Availability")).
		Return(&event, nil)

	router := gin.New()
	router.POST("/availability", AddAvailabilityHandler(mockSvc))

	body, _ := json.Marshal(event)
	req, _ := http.NewRequest(http.MethodPost, "/availability", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockSvc.AssertExpectations(t)
}

// --------- DELETE /availability/:id -----------

func TestDeleteAvailabilityHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mocker.MockAvailabilityService)

	id := primitive.NewObjectID().Hex()

	mockSvc.On("DeleteAvailability", mock.Anything, id).Return(nil)

	router := gin.New()
	router.DELETE("/availability/:id", DeleteAvailabilityHandler(mockSvc))

	req, _ := http.NewRequest(http.MethodDelete, "/availability/"+id, nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNoContent, resp.Code)
}

// --------- UPDATE /availability/:id -----------

func TestUpdateAvailabilityHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mocker.MockAvailabilityService)

	id := primitive.NewObjectID().String()
	avail := &models.Availability{EventID: primitive.NewObjectID(), ID: primitive.NewObjectID()}

	mockSvc.On("UpdateAvailability", mock.Anything, id).Return(nil)

	router := gin.New()
	router.PUT("/availability/:id", UpdateAvailabilityHandler(mockSvc))

	body, _ := json.Marshal(avail)
	req, _ := http.NewRequest(http.MethodPut, "/availability"+id, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

// --------- GET /event/:id/availability -----------

func TestGetAvailabilityByEventHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockSvc := new(mocker.MockAvailabilityService)

	id := primitive.NewObjectID().String()
	avail := &models.Availability{EventID: primitive.NewObjectID(), ID: primitive.NewObjectID()}

	mockSvc.On("GetAvailabilitiesByEvent", mock.Anything, id).Return(avail, nil)

	router := gin.New()
	router.GET("/event/:id/availability", GetAvailabilityByEventHandler(mockSvc))

	req, _ := http.NewRequest(http.MethodGet, "/events/"+id+"/availability", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

// func TestCreateEventHandler_Success(t *testing.T) {
// 	gin.SetMode(gin.TestMode)

// 	mockService := new(mocker.MockEventService)
// 	handler := CreateEventHandler(mockService)

// 	// Sample input/output
// 	input := models.Event{
// 		Title: "Team Sync",
// 	}
// 	expected := &models.Event{
// 		ID:    primitive.NewObjectID(),
// 		Title: "Team Sync",
// 	}

// 	// Marshal input to JSON
// 	body, _ := json.Marshal(input)

// 	// Create test request/response recorder
// 	req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(body))
// 	req.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()

// 	// Create Gin context manually
// 	c, _ := gin.CreateTestContext(w)
// 	c.Request = req
// 	// Set mock expectation
// 	mockService.
// 		On("CreateEvent", c, mock.AnythingOfType("models.Event")).
// 		Return(expected, nil)

// 	// Run handler
// 	handler(c)

// 	// Assert
// 	assert.Equal(t, http.StatusOK, w.Code)

// 	var response models.Event
// 	err := json.Unmarshal(w.Body.Bytes(), &response)
// 	assert.NoError(t, err)
// 	assert.Equal(t, expected.ID, response.ID)
// 	assert.Equal(t, expected.Title, response.Title)

// 	mockService.AssertExpectations(t)
// }
