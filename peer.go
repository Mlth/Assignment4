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

var reader *bufio.Reader

func main() {
	reader = bufio.NewReader(os.Stdin)
	var wg sync.WaitGroup
	wg.Add(1)
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	arg2, _ := strconv.ParseInt(os.Args[2], 10, 32)
	ownPort := int32(arg1) + 5000
	totalPorts := int32(arg2)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	p := &peer{
		id:                           ownPort,
		wantsToken:                   false,
		previousPeerHasConnectedToMe: false,
		nextPeer:                     nil,
		ctx:                          ctx,
	}

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

	fmt.Printf("Trying to dial: %v\n", (((ownPort + 1) % totalPorts) + 5000))
	conn, err := grpc.Dial(fmt.Sprintf(":%v", ((ownPort+1)%totalPorts)+5000), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()
	p.nextPeer = ring.NewRingClient(conn)

	go takeInput(p, &wg)

	if ownPort == 5000 {
		p.nextPeer.CheckConnection(p.ctx, &ring.ConnectionVerification{Id: p.id})
		if p.previousPeerHasConnectedToMe {
			p.nextPeer.PassToken(ctx, &ring.Token{})
		}
	}

	wg.Wait()
}

func takeInput(p *peer, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		log.Println("Do you want to enter the critical state? (type 'yes' or 'no')")
		inputMessage, _ := reader.ReadString('\n')
		inputMessage = strings.TrimSpace(inputMessage)
		if inputMessage == "yes" {
			p.wantsToken = true
			log.Println("Okay! Wait for the token to be passed along to you.")
			for p.wantsToken {

			}
		}
	}
}

func (p *peer) CheckConnection(ctx context.Context, msg *ring.ConnectionVerification) (*ring.EmptyMessage, error) {
	if msg.Id == p.id {
		p.previousPeerHasConnectedToMe = true
	} else {
		for p.nextPeer == nil {

		}
		p.nextPeer.CheckConnection(p.ctx, &ring.ConnectionVerification{Id: msg.Id})
	}
	return &ring.EmptyMessage{}, nil
}

func (p *peer) PassToken(ctx context.Context, token *ring.Token) (*ring.EmptyMessage, error) {
	if !p.wantsToken {
		go p.nextPeer.PassToken(p.ctx, &ring.Token{})
	} else {
		log.Println("You are now in the critical state. Do you want to exit the state? (type 'yes' or 'no')")
		inputMessage, _ := reader.ReadString('\n')
		inputMessage = strings.TrimSpace(inputMessage)
		if inputMessage == "yes" {
			log.Println("Okay! You have passed the token along")
			p.wantsToken = false
			go p.nextPeer.PassToken(p.ctx, &ring.Token{})
		}
	}
	return &ring.EmptyMessage{}, nil
}

type peer struct {
	ring.UnimplementedRingServer
	id                           int32
	wantsToken                   bool
	previousPeerHasConnectedToMe bool
	nextPeer                     ring.RingClient
	ctx                          context.Context
}
