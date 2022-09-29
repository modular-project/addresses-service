package storage

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/modular-project/address-service/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	address = []model.Address{
		{
			ID:         primitive.NewObjectID(),
			Street:     "test 1",
			City:       "city test",
			PostalCode: "11111",
			State:      "test-s",
			Country:    "test-c",
		},
	}
)

func initDB(t *testing.T, db *mongo.Database) {
	col := db.Collection("delivery")
	des := []model.Delivery{
		{
			UserID:  1,
			Address: address[0],
		},
	}
	ides := make([]interface{}, len(des))
	for i := range des {
		ides[i] = des[i]
		t.Logf("ID: %s", des[i].ID.Hex())
	}
	r, err := col.InsertMany(context.Background(), ides)
	if err != nil {
		t.Fatalf("InserMany: %s", err)
	}
	if len(r.InsertedIDs) != len(des) {
		t.Fatalf("len inserted(%d) is not equal to len des(%d)", len(r.InsertedIDs), len(des))
	}
}

func dropTest(t *testing.T, db *mongo.Database) {
	if db.Name() != "test" {
		t.Fatalf("db is not tes")
		return
	}
	col := db.Collection("delivery")
	if err := col.Drop(context.Background()); err != nil {
		t.Fatalf("drop: %s", err)
	}
}

func newTestConnection() DBConnection {
	env := "ADDR_DB_HOST"
	host, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "ADDR_DB_USER"
	user, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "ADDR_DB_PWD"
	pwd, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "ADDR_DB_NAME"
	cluster, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	return DBConnection{User: user, Host: host, Password: pwd, Cluster: cluster, NameDB: "test"}
}

func TestDeliveryStorage_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		d   model.Delivery
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "insert OK",
			args: args{
				ctx: context.Background(),
				d: model.Delivery{
					UserID: 1,
					Address: model.Address{
						Street:     "test 1",
						City:       "city test",
						PostalCode: "11111",
						State:      "test-s",
						Country:    "test-c",
					},
				},
			},
		},
	}
	conn := newTestConnection()
	db, err := NewDB(&conn)
	if err != nil {
		t.Fatalf("failed to NewDB: %s", err)
	}
	ds := NewDeliveryStorage(db, "")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ds.Create(tt.args.ctx, &tt.args.d)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeliveryStorage.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != 12*2 {
				t.Errorf("DeliveryStorage.Create() len is not 24, got: %d, %s", len(got), got)
			}
		})
	}
}

func TestNewDeliveryStorage(t *testing.T) {
	type args struct {
		db *mongo.Database
	}
	tests := []struct {
		name string
		args args
		want DeliveryStorage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDeliveryStorage(tt.args.db, ""); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDeliveryStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeliveryStorage_GetAll(t *testing.T) {
	type args struct {
		ctx context.Context
		uID uint64
	}
	tests := []struct {
		name    string
		ds      DeliveryStorage
		args    args
		want    []model.Address
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ds.GetAll(tt.args.ctx, tt.args.uID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeliveryStorage.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeliveryStorage.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeliveryStorage_GetByID(t *testing.T) {
	type args struct {
		ctx context.Context
		uID uint64
		aID string
	}
	tests := []struct {
		name    string
		args    args
		want    model.Address
		wantErr bool
	}{
		{
			name: "ok user 1",
			want: address[0],
			args: args{ctx: context.Background(), uID: 1, aID: address[0].ID.Hex()},
		}, {
			name:    "forbidden, user 2 get user 1",
			args:    args{ctx: context.Background(), uID: 2, aID: address[0].ID.Hex()},
			wantErr: true,
		}, {
			name:    "aID not found",
			args:    args{ctx: context.Background(), uID: 1, aID: primitive.NewObjectID().Hex()},
			wantErr: true,
		},
	}
	conn := newTestConnection()
	db, err := NewDB(&conn)
	if err != nil {
		t.Errorf("newDB: %s", err)
	}
	initDB(t, db)
	t.Cleanup(func() { dropTest(t, db) })
	ds := NewDeliveryStorage(db, "")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ds.GetByID(tt.args.ctx, tt.args.uID, tt.args.aID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeliveryStorage.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeliveryStorage.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDeliveryStorage_DeleteByID(t *testing.T) {
	type args struct {
		ctx context.Context
		in1 uint64
		in2 string
	}
	tests := []struct {
		name    string
		ds      DeliveryStorage
		args    args
		wantErr bool
		deleted int64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ds.DeleteByID(tt.args.ctx, tt.args.in1, tt.args.in2)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeliveryStorage.DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.deleted != got {
				t.Errorf("DeliveryStorage.DeleteByID() gotDeleted = %d, wantDeleted %d", got, tt.deleted)
			}
		})
	}
}
