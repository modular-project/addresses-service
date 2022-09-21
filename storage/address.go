package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/modular-project/address-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AddressStorage struct {
	c      *mongo.Collection
	maxDis int
}

func NewAddressStorage(db *mongo.Database, max int) AddressStorage {
	return AddressStorage{c: db.Collection("establishment"), maxDis: max}
}

func (as AddressStorage) GetByID(ctx context.Context, aID string) (model.Address, error) {
	var a model.Address
	id, err := primitive.ObjectIDFromHex(aID)
	if err != nil {
		return model.Address{}, fmt.Errorf("ObjectIDFromHex: %w", err)
	}
	r := as.c.FindOne(ctx, bson.M{"_id": id})
	if r.Err() != nil {
		return model.Address{}, fmt.Errorf("findOne: %w", r.Err())
	}
	if err := r.Decode(&a); err != nil {
		return model.Address{}, fmt.Errorf("decode: %w", err)
	}
	return a, nil
}

func (as AddressStorage) Create(ctx context.Context, add *model.Address) (string, error) {
	r, err := as.c.InsertOne(ctx, add)
	if err != nil {
		return "", fmt.Errorf("InsertOne: %w", err)
	}
	id, ok := r.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("InsertOneResult is not an ObjectID")
	}
	return id.Hex(), nil
}

func (as AddressStorage) DeleteByID(ctx context.Context, aID string) (int64, error) {
	id, err := primitive.ObjectIDFromHex(aID)
	if err != nil {
		return 0, fmt.Errorf("ObjectIDFromHex: %w", err)
	}
	r, err := as.c.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return 0, fmt.Errorf("DeleteOne: %w", err)
	}
	return r.DeletedCount, nil
}

func (as AddressStorage) Search(ctx context.Context, s *model.Search) ([]model.Address, error) {
	var ads []model.Address
	opt := options.FindOptions{
		Limit: &s.Limit,
		Skip:  &s.Offset,
		Sort:  s.OrderBy,
	}
	log.Println(s.Querys, s.OrderBy)
	r, err := as.c.Find(ctx, s.Querys, &opt)
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}
	if err := r.All(ctx, &ads); err != nil {
		return nil, fmt.Errorf("decode all: %w", err)
	}
	return ads, nil
}

func (as AddressStorage) Nearest(ctx context.Context, loc []float64) (string, error) {
	var near model.Address
	opts := options.FindOne().SetProjection(bson.M{"_id": 1})
	r := as.c.FindOne(ctx, bson.M{
		"location": bson.M{
			"$near": bson.M{
				"$geometry": bson.M{
					"type":        "Point",
					"coordinates": loc,
				},
				"$maxDistance": as.maxDis,
			},
		},
	}, opts)
	if r.Err() != nil {
		return "", fmt.Errorf("findOne: %w", r.Err())
	}
	if err := r.Decode(&near); err != nil {
		return "", fmt.Errorf("decode: %w", err)
	}
	return near.ID.Hex(), nil
}
