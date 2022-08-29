package storage

import (
	"context"
	"fmt"

	"github.com/modular-project/address-service/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeliveryStorage struct {
	c *mongo.Collection
}

func NewDeliveryStorage(db *mongo.Database) DeliveryStorage {
	return DeliveryStorage{db.Collection("delivery")}
}

func (ds DeliveryStorage) Create(ctx context.Context, d *model.Delivery) (string, error) {
	r, err := ds.c.InsertOne(ctx, d)
	if err != nil {
		return "", fmt.Errorf("InsertOne: %w", err)
	}
	id, ok := r.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("InsertOneResult is not an ObjectID")
	}
	return id.Hex(), nil
}

func (ds DeliveryStorage) GetAll(ctx context.Context, uID uint64) ([]model.Address, error) {
	var as []model.Address
	r, err := ds.c.Find(ctx, bson.M{"user_id": uID})
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}
	if err := r.All(ctx, &as); err != nil {
		return nil, fmt.Errorf("decode all: %w", err)
	}
	return as, nil
}

func (ds DeliveryStorage) GetByID(ctx context.Context, uID uint64, aID string) (model.Address, error) {
	var a model.Address
	id, err := primitive.ObjectIDFromHex(aID)
	if err != nil {
		return model.Address{}, fmt.Errorf("ObjectIDFromHex: %w", err)
	}
	r := ds.c.FindOne(ctx, bson.M{"_id": id, "user_id": uID})
	if r.Err() != nil {
		return model.Address{}, fmt.Errorf("findOne: %w", r.Err())
	}
	if err := r.Decode(&a); err != nil {
		return model.Address{}, fmt.Errorf("decode: %w", err)
	}
	return a, nil
}

func (ds DeliveryStorage) DeleteByID(ctx context.Context, uID uint64, aID string) (int64, error) {
	id, err := primitive.ObjectIDFromHex(aID)
	if err != nil {
		return 0, fmt.Errorf("ObjectIDFromHex: %w", err)
	}
	r, err := ds.c.DeleteOne(ctx, bson.M{"user_id": uID, "_id": id})
	if err != nil {
		return 0, fmt.Errorf("DeleteOne: %w", err)
	}
	return r.DeletedCount, nil
}
