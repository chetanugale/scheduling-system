package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/chetanugale/scheduling-system/constants"
	"github.com/chetanugale/scheduling-system/handlers"
	"github.com/chetanugale/scheduling-system/models"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/chetanugale/scheduling-system/repository"
	"github.com/chetanugale/scheduling-system/services"
)

func main() {
	router := gin.Default()
	ctx := context.Background()

	eventService, availService := dbInitializer(ctx)

	router = registerApi(router, eventService, availService)

	router.Run(constants.PORT)

}

func registerApi(router *gin.Engine, eventService *services.MongoEventService, availService *services.MongoAvailabilityService) *gin.Engine {
	// ----- Event management

	router.POST("/events", handlers.CreateEventHandler(eventService))       //create event   // TODO : add validators for duplicate data
	router.GET("/events", handlers.GetAllEventsHandler(eventService))       // get all events
	router.GET("/events/:id", handlers.GetEventHandler(eventService))       // get event with ID
	router.PUT("/events/:id", handlers.UpdateEventHandler(eventService))    // update event with ID
	router.DELETE("/events/:id", handlers.DeleteEventHandler(eventService)) // delete event with ID

	// ----- Availability management

	router.POST("/availability", handlers.AddAvailabilityHandler(availService))                 // every user will post availability    // TODO : add validators for duplicate data
	router.GET("/event/:id/availability", handlers.GetAvailabilityByEventHandler(availService)) // TODO : optimize this API
	router.PUT("/availability/:id", handlers.UpdateAvailabilityHandler(availService))
	router.DELETE("/availability/:id", handlers.DeleteAvailabilityHandler(availService))

	// ------ Recommendation

	router.GET("/events/:id/recommend", handlers.RecommendHandler(eventService, availService)) // recommend availability based on id for max users

	return router
}

func dbInitializer(ctx context.Context) (*services.MongoEventService, *services.MongoAvailabilityService) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(constants.MONGO_URI))
	if err != nil {
		log.Fatal(err)
	}

	db := client.Database(constants.DB_NAME)
	eventRepo := repository.NewMongoRepository[models.Event](db.Collection(constants.COLL_EVENTS))
	availRepo := repository.NewMongoRepository[models.Availability](db.Collection(constants.COLL_AVAIL))

	eventService := &services.MongoEventService{Repo: eventRepo}
	availService := &services.MongoAvailabilityService{Repo: availRepo}

	return eventService, availService
}
