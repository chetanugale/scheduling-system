package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/chetanugale/scheduling-system/models"
	"github.com/chetanugale/scheduling-system/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateEventHandler(svc services.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var e models.Event
		if err := c.ShouldBindJSON(&e); err != nil {
			c.JSON(http.StatusBadRequest, bson.M{"error": err.Error()})
			return
		}
		// TODO : optimize and cleanup later
		// e.ID = uuid.New().String()
		for data := range e.Slots {
			e.Slots[data].ID = primitive.NewObjectID()
		}

		created, err := svc.CreateEvent(c, e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, bson.M{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, created)
	}
}

func GetEventHandler(svc services.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		event, err := svc.GetEvent(c, id)
		if err != nil {
			c.JSON(http.StatusNotFound, bson.M{"error": "Event not found"})
			return
		}
		c.JSON(http.StatusOK, event)
	}
}

func GetAllEventsHandler(svc services.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Query("title")
		event, err := svc.GetAllEvents(c, title)
		if err != nil {
			c.JSON(http.StatusNotFound, bson.M{"error": "Events not found.Empty Dataset"})
			return
		}
		c.JSON(http.StatusOK, event)
	}
}

func UpdateEventHandler(svc services.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var event models.Event
		if err := c.ShouldBindJSON(&event); err != nil {
			c.JSON(http.StatusBadRequest, bson.M{"error": err.Error()})
			return
		}
		err := svc.UpdateEvent(c, id, event)
		if err != nil {
			log.Fatalf("%+v", err)
			c.JSON(http.StatusInternalServerError, bson.M{"error": "Error while updating event."})
			return
		}
		c.JSON(http.StatusOK, event)
	}
}

func DeleteEventHandler(svc services.EventService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := svc.DeleteEvent(c, id); err != nil {
			c.JSON(http.StatusInternalServerError, bson.M{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func AddAvailabilityHandler(svc services.AvailabilityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var a models.Availability
		if err := c.ShouldBindJSON(&a); err != nil {
			c.JSON(http.StatusBadRequest, bson.M{"error": err.Error()})
			return
		}
		a.ID = primitive.NewObjectID()            // TODO : optimize and cleanup later
		created, err := svc.AddAvailability(c, a) // TODO : validate if eventId and SlotId present
		if err != nil {
			c.JSON(http.StatusInternalServerError, bson.M{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, created)
	}
}

func GetAvailabilityByEventHandler(svc services.AvailabilityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventId := c.Param("id")
		availList, err := svc.GetAvailabilitiesByEvent(c, eventId)
		if err != nil {
			c.JSON(http.StatusNotFound, fmt.Sprintf("%+v", err.Error()))
		}
		c.JSON(http.StatusOK, availList)
	}
}

func DeleteAvailabilityHandler(svc services.AvailabilityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		if err := svc.DeleteAvailability(c, id); err != nil {
			c.JSON(http.StatusInternalServerError, bson.M{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}

func UpdateAvailabilityHandler(svc services.AvailabilityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var avail models.Availability
		if err := c.ShouldBindJSON(&avail); err != nil {
			c.JSON(http.StatusBadRequest, bson.M{"error": err.Error()})
			return
		}
		err := svc.UpdateAvailability(c, id, avail)
		if err != nil {
			log.Fatalf("%+v", err)
			c.JSON(http.StatusInternalServerError, bson.M{"error": "Error while updating availability."})
			return
		}
		c.JSON(http.StatusOK, avail)
	}
}

func processRecommendations(event models.Event, available []models.Availability) ([]models.TimeSlot, map[string][]string, error) {
	// event
	// {
	// 	"title":"test6",
	// 	"estimatedMins":30,
	// 	"slots":[
	// 	{
	// 		"_id":"abcd",
	// 		"startTime":"2025-05-04T09:00:00Z",
	// 		"endTime":"2025-05-04T09:30:00Z"
	// 	},
	// 	{
	//		"_id":"pqrs",
	// 		"startTime":"2025-05-04T14:00:00Z",
	// 		"endTime":"2025-05-04T14:30:00Z"
	// 	}
	// 	]
	// }

	//availability
	// {
	//     "eventId":"6817a28c6d1b32a2fd46ec16",
	//     "slotId":"6817a28c6d1b32a2fd46ec14",
	//     "userId":"abcd"
	// }
	//target : timeslot which has max no of users available for event, list of users not available for max timeslot
	//1. fetch users for matching slotid and eventid and reverse the query for not matching timeslots

	slotList := map[string][]string{}
	for _, avail := range available {
		slotList[avail.SlotID.String()] = append(slotList[avail.SlotID.String()], avail.UserID)
	}

	var preciseUsers int
	idealSlots := []models.TimeSlot{}
	userSlotsList := map[string][]string{}

	for _, slot := range event.Slots {
		users := slotList[slot.ID.String()]
		if len(users) > preciseUsers {
			preciseUsers = len(users)
			idealSlots = []models.TimeSlot{slot}
		} else if len(users) == preciseUsers {
			idealSlots = append(idealSlots, slot)
		}
	}

	usersSet := map[string]struct{}{}
	for _, a := range available {
		usersSet[a.UserID] = struct{}{}
	}
	for _, slot := range idealSlots {
		presentUser := map[string]bool{}
		for _, userid := range slotList[slot.ID.String()] {
			presentUser[userid] = true
		}
		for userid := range usersSet {
			if !presentUser[userid] {
				userSlotsList[slot.ID.Hex()] = append(userSlotsList[slot.ID.Hex()], userid)
			}
		}
	}
	return idealSlots, userSlotsList, nil
}

func RecommendHandler(svcEvent services.EventService, svcAvail services.AvailabilityService) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventId := c.Param("id")
		event, err := svcEvent.GetEvent(c, eventId)
		if err != nil {
			c.JSON(http.StatusNotFound, fmt.Sprintf("%+v", err.Error()))
		}
		availList, err := svcAvail.GetAvailabilitiesByEvent(c, eventId)
		if err != nil {
			c.JSON(http.StatusNotFound, fmt.Sprintf("%+v", err.Error()))
		}
		idealSlots, notfeasible, err := processRecommendations(*event, availList) // TODO : Optimize this function
		if err != nil {
			c.JSON(http.StatusPreconditionFailed, fmt.Sprintf("%+v", err.Error()))
		}
		c.JSON(http.StatusOK, gin.H{"IdealSlots": idealSlots, "NotFeasibleforUsers": notfeasible})
	}
}
