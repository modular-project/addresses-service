package model

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Sort int

type By int

const (
	ASC Sort = iota
	DES
)

type Location struct {
	// long between -180 and 180
	// lat between -90 and 90
	Type        string    `json:"-"`
	Coordinates []float64 `json:"-"` // long, Lat
}

type Address struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Street     string             `bson:"street,omitempty"`
	Suburb     string             `bson:"suburb,omitempty"`
	City       string             `bson:"city,omitempty"`
	PostalCode string             `bson:"pc,omitempty"`
	State      string             `bson:"state,omitempty"`
	Country    string             `bson:"country,omitempty"`
	Location   Location           `bson:"location"`
	IsDeleted  bool               `bson:is_deleted,omitempty`
}

type Delivery struct {
	Address `bson:",inline"`
	UserID  uint64 `bson:"user_id,omitempty"`
}

type OrderBy struct {
	Sort int32  `json:"by,omitempty"`
	By   string `json:"sort,omitempty"`
}

type Search struct {
	Limit   int64
	Offset  int64
	OrderBy primitive.D
	Querys  primitive.D
}

func (a Address) String() string {
	return fmt.Sprintf("%s, %s, %s, %s, %s, %s",
		a.Street, a.Suburb, a.PostalCode, a.City, a.State, a.Country)
}

// func (s Search) Query() primitive.D {
// 	if s.Querys == nil {
// 		return primitive.D{}
// 	}
// 	q := make([]primitive.E, len(s.OrderBy))
// 	for i := range s.Querys {
// 		q[i].Key = s.Querys[i].key
// 		q[i].Value = s.Querys[i].val
// 	}
// 	return q
// }

// func (s Search) Order() primitive.D {
// 	if s.OrderBy == nil {
// 		return primitive.D{}
// 	}
// 	o := make([]primitive.E, len(s.OrderBy))
// 	for i := range s.OrderBy {
// 		o[i].Key = s.OrderBy[i].By
// 		o[i].Value = s.OrderBy[i].Sort
// 	}
// 	return o
// }
