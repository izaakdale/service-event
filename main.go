package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	db "github.com/izaakdale/service-event/datastore/sqlc"
	"github.com/izaakdale/service-event/schema/event"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var (
	querier db.Querier
)

type (
	Service struct {
		Name       string
		GrpcServer *grpc.Server
	}

	GServer struct {
		event.EventServiceServer
	}
)

func (g *GServer) GetEvent(ctx context.Context, ev *event.EventRequest) (*event.EventResponse, error) {
	dbe, err := querier.GetEvent(ctx, ev.EventId)
	if err != nil {
		return nil, err
	}
	return &event.EventResponse{
		EventId:          dbe.EventID,
		EventName:        dbe.EventName,
		TicketsRemaining: dbe.TicketsRemaining,
		Timestamp:        dbe.EventTimestamp.GoString(),
	}, nil
}

func main() {

	dbConn, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/events?sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to connect to db %s", err.Error())
	}

	querier = db.New(dbConn)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", "localhost", "50001"))
	if err != nil {
		log.Fatalf("Failed to listen on %v\n", err)
	}

	gsrv := grpc.NewServer()
	ls := &GServer{}
	event.RegisterEventServiceServer(gsrv, ls)

	log.Fatal(gsrv.Serve(lis))
}
