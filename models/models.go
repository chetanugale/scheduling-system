package models

import (
    "go.mongodb.org/mongo-driver/bson/primitive"
    "time"
)

type TimeSlot struct {
    ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    StartTime time.Time          `json:"startTime"`
    EndTime   time.Time          `json:"endTime"`
}

type Event struct {
    ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Title         string             `json:"title"`
    EstimatedMins int                `json:"estimatedMins"`
    Slots         []TimeSlot         `json:"slots"`
}

type Availability struct {
    ID      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    EventID primitive.ObjectID `json:"eventId"`
    UserID  string             `json:"userId"`
    SlotID  primitive.ObjectID `json:"slotId"`
}

type User struct{
	UserID string `json:"userId"`
	UserName string `json:"userName"`
}

