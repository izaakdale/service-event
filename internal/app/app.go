package app

import (
	"database/sql"
	"fmt"
	"log"
	"net"

	db "github.com/izaakdale/service-event/internal/datastore/sqlc"
	"github.com/izaakdale/service-event/pkg/schema/event"
	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
)

var (
	querier db.Querier
	name    = "service-event"
	spec    specification
)

type (
	Service struct {
		Name       string
		GrpcServer *grpc.Server
	}
	GServer struct {
		event.EventServiceServer
	}
	specification struct {
		GRPCHost         string
		GRPCPort         string
		DBDriver         string
		DBDataSourceName string
	}
)

func New() *Service {
	err := envconfig.Process("", &spec)
	if err != nil {
		panic(err)
	}

	dbConn, err := sql.Open(spec.DBDriver, spec.DBDataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to db %s", err.Error())
	}
	querier = db.New(dbConn)

	gsrv := grpc.NewServer()
	ls := &GServer{}
	event.RegisterEventServiceServer(gsrv, ls)

	return &Service{name, gsrv}
}

func (s *Service) Run() {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", spec.GRPCHost, spec.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen on %v\n", err)
	}
	log.Printf("listening for GRPC clients on %s:%s\n", spec.GRPCHost, spec.GRPCPort)
	log.Fatal(s.GrpcServer.Serve(lis))
}
