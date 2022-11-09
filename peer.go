package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	ring "github.com/Mlth/Assignment4/proto"
	"google.golang.org/grpc"
)

func main() {
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + 5000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:       ownPort,
		nextPeer: nil,
		ctx:      ctx,
	}

	var maxClient int32 = 5

	// Create listener tcp on port ownPort

	//fixme: noget galt med matematikken her!!
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ((ownPort-1)%maxClient)+5000))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	grpcServer := grpc.NewServer()
	ring.RegisterRingServer(grpcServer, p)

	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}

	}()

	fmt.Printf("Trying to dial: %v\n", ownPort+1)
	conn, err := grpc.Dial(fmt.Sprintf(":%v", ownPort+1), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()
	p.nextPeer = ring.NewRingClient(conn)

}

type peer struct {
	ring.UnimplementedRingServer
	id       int32
	nextPeer ring.RingClient
	ctx      context.Context
}
