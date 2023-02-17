package app

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/izaakdale/lib/publisher"
	db "github.com/izaakdale/service-event/internal/datastore/sqlc"
	"github.com/izaakdale/service-event/pkg/proto/event"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	name             = "service-event"
	querier          db.Querier
	spec             specification
	grpcServerSocket string
	grpcMuxSocket    string
)

type (
	GServer struct {
		event.EventServiceServer
	}
	specification struct {
		Host             string `envconfig:"HOST"`
		Port             string `envconfig:"PORT"`
		GRPCHost         string `envconfig:"GRPC_HOST"`
		GRPCPort         string `envconfig:"GRPC_PORT"`
		DBDriver         string `envconfig:"DB_DRIVER"`
		DBDataSourceName string `envconfig:"DB_DATA_SOURCE_NAME"`
		AWSRegion        string `envconfig:"AWS_REGION" default:"eu-west-2"`
		TopicArn         string `envconfig:"TOPIC_ARN"`
		AWSEndpoint      string `envconfig:"AWS_ENDPOINT"`
	}
)

func Run() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := envconfig.Process("", &spec)
	if err != nil {
		panic(err)
	}

	grpcServerSocket = fmt.Sprintf("%s:%s", spec.GRPCHost, spec.GRPCPort)
	grpcMuxSocket = fmt.Sprintf("%s:%s", spec.Host, spec.Port)

	cfg, err := config.LoadDefaultConfig(context.Background(), func(o *config.LoadOptions) error {
		o.Region = spec.AWSRegion
		return nil
	})
	if err != nil {
		panic(err)
	}
	err = publisher.Initialise(cfg, spec.TopicArn, publisher.WithEndpoint(spec.AWSEndpoint))
	if err != nil {
		panic(err)
	}

	// datastore setup
	dbConn, err := sql.Open(spec.DBDriver, spec.DBDataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to db %s", err.Error())
	}
	querier = db.New(dbConn)

	// grpc server setup
	lis, err := net.Listen("tcp", grpcServerSocket)
	if err != nil {
		log.Fatalf("Failed to listen on %v\n", err)
	}
	gsrv := grpc.NewServer()
	ls := &GServer{}
	event.RegisterEventServiceServer(gsrv, ls)

	// grpc mux setup
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = event.RegisterEventServiceHandlerFromEndpoint(ctx, mux, grpcServerSocket, opts)
	if err != nil {
		panic(err)
	}

	// kick off server and mux
	errChan := make(chan error, 1)
	go func(errChan chan<- error) {
		log.Printf("listening for GRPC clients on %s\n", grpcServerSocket)
		errChan <- gsrv.Serve(lis)
	}(errChan)
	go func(errChan chan<- error) {
		log.Printf("serving http clients on %s\n", grpcMuxSocket)
		errChan <- http.ListenAndServe(grpcMuxSocket, mux)
	}(errChan)

	// subscribe to shutdown signals
	shutCh := make(chan os.Signal, 1)
	signal.Notify(shutCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// wait on error from server, mux or shutdown sig.
	select {
	case err = <-errChan:
		if err != nil {
			log.Fatal(err)
		}
	case signal := <-shutCh:
		log.Printf("got shutdown signal: %s, exiting\n", signal)
	}
}
