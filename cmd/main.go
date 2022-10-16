package main

import (
	"fmt"
	"log"
	"net"
	"os"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	gmaps "github.com/modular-project/address-service/adapter/gmap"
	"github.com/modular-project/address-service/controller"
	"github.com/modular-project/address-service/http/handler"
	"github.com/modular-project/address-service/storage"
	pf "github.com/modular-project/protobuffers/address/address"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func newDBConnection() storage.DBConnection {
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
	return storage.DBConnection{User: user, Host: host, Password: pwd, Cluster: cluster, NameDB: "modular"}
}

func startGRPC() *grpc.Server {
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(Recovery),
	}
	server := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(opts...),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_middleware.ChainStreamServer(),
		)),
	)
	return server
}

func Recovery(i interface{}) error {
	return status.Errorf(codes.Unknown, "panic triggered: %v", i)
}

func main() {
	conn := newDBConnection()
	db, err := storage.NewDB(&conn)
	if err != nil {
		log.Fatalf("NewDB: %s", err)
	}
	coll, _ := os.LookupEnv("ADDR_COLLECTION")
	ast := storage.NewAddressStorage(db, 25000, coll)
	coll, _ = os.LookupEnv("DEL_COLLECTION")
	dst := storage.NewDeliveryStorage(db, coll)
	key, ok := os.LookupEnv("GMAP_APIKEY")
	if !ok {
		log.Fatal("enviroment variable GMAP_APIKEY not found")
	}
	gms, err := gmaps.NewGMapService(key)
	if err != nil {
		log.Fatalf("NewGMapService: %s", err)
	}
	env := "ADDR_PORT"
	port, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ads := controller.NewAddressService(ast, dst, gms)
	auc := handler.NewAddressUC(ads)
	srv := startGRPC()
	pf.RegisterAddressServiceServer(srv, auc)
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(pf.AddressService_ServiceDesc.ServiceName, healthpb.HealthCheckResponse_SERVING)
	log.Printf("Server started at :%s", port)
	healthpb.RegisterHealthServer(srv, healthServer)
	err = srv.Serve(lis)
	if err != nil {
		log.Fatalf("failed to server at :%s, got error: %s", port, err)
	}
}
