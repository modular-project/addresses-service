package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	once sync.Once
)

type DBConnection struct {
	User     string
	Password string
	NameDB   string
	Cluster  string
	Host     string
}

func NewDB(conn *DBConnection) (*mongo.Database, error) {
	var err error
	var client *mongo.Client
	once.Do(func() {
		dns := fmt.Sprintf("%s://%s:%s@%s?retryWrites=true&w=majority",
			conn.Host, conn.User, conn.Password, conn.Cluster)
		serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
		clientOptions := options.Client().
			ApplyURI(dns).
			SetServerAPIOptions(serverAPIOptions)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		client, err = mongo.Connect(ctx, clientOptions)
	})
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}
	if client == nil {
		return nil, fmt.Errorf("nil client")
	}
	return client.Database(conn.NameDB), nil
}
