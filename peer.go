package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	ring "github.com/Mlth/Assignment4/proto"
	"google.golang.org/grpc"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	ownPort := int32(arg1) + 5000

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:              ownPort,
		participant:     false,
		inCriticalState: false,
		nextPeer:        nil,
		ctx:             ctx,
	}

	var maxClient int32 = 5

	// Create listener tcp on port ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
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

	fmt.Printf("Trying to dial: %v\n", (((ownPort + 1) % maxClient) + 5000))
	conn, err := grpc.Dial(fmt.Sprintf(":%v", ((ownPort+1)%maxClient)+5000), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()
	p.nextPeer = ring.NewRingClient(conn)

	go takeInput(p, &wg)
	wg.Wait()
}

func takeInput(p *peer, wg *sync.WaitGroup) {
	defer wg.Done()
	reader := bufio.NewReader(os.Stdin)
	for {
		if !p.participant {
			log.Println("Do you want to enter the critical state? (type 'yes' or 'no')")
			inputMessage, _ := reader.ReadString('\n')
			inputMessage = strings.TrimSpace(inputMessage)
			if inputMessage == "yes" {
				p.participant = true
				p.nextPeer.Send(p.ctx, &ring.Request{Id: p.id, Type: "election"})
				p.participant = false
				for {
					log.Println("Do you want to exit the critical state? (type 'yes' or 'no')")
					inputMessage2, _ := reader.ReadString('\n')
					inputMessage2 = strings.TrimSpace(inputMessage2)
					if inputMessage2 == "yes" {
						p.inCriticalState = false
						break
					}
				}
			}
		}
	}
}

func (p *peer) Send(ctx context.Context, req *ring.Request) (*ring.EmptyMessage, error) {
	if req.Type == "election" {
		for p.inCriticalState {

		}
	}
	if req.Type == "election" {
		if p.id == req.Id {
			log.Println("Peer with id ", p.id, " is in critical state")
			p.inCriticalState = true
			p.participant = false
			p.nextPeer.Send(ctx, &ring.Request{Id: p.id, Type: "elected"})
		} else if req.Id > p.id || !p.participant {
			p.nextPeer.Send(ctx, &ring.Request{Id: req.Id, Type: "election"})
		} else if req.Id < p.id && p.participant {
			p.nextPeer.Send(ctx, &ring.Request{Id: p.id, Type: "election"})
		}
	} else if req.Type == "elected" {
		if p.id != req.Id {
			log.Println("Peer with id ", req.Id, " is in critical state")
			p.nextPeer.Send(ctx, &ring.Request{Id: req.Id, Type: "elected"})
		}
	}
	return &ring.EmptyMessage{}, nil
}

type peer struct {
	ring.UnimplementedRingServer
	id              int32
	participant     bool
	inCriticalState bool
	nextPeer        ring.RingClient
	ctx             context.Context
}
