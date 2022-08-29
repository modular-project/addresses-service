package model

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Location struct {
	Long float64 // between -180 and 180
	Lat  float64 // between -90 and 90
}

type Address struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Street     string             `bson:"street,omitempty"`
	Suburb     string             `bson:"suburb,omitempty"`
	City       string             `bson:"city,omitempty"`
	PostalCode string             `bson:"pc,omitempty"`
	State      string             `bson:"state,omitempty"`
	Country    string             `bson:"country,omitempty"`
	Location   Location           `bson:",inline"`
}

type Delivery struct {
	Address `bson:",inline"`
	UserID  uint64 `bson:"user_id,omitempty"`
}

type Search struct {
	Limit  int64
	Offset int64
}

func (a Address) String() string {
	return fmt.Sprintf("%s, %s, %s, %s, %s, %s",
		a.Street, a.Suburb, a.PostalCode, a.City, a.State, a.Country)
}
