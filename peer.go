package main

import (
	"bufio"
	"context"
	"fmt"
	ring "github.com/Mlth/Assignment4/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

var reader *bufio.Reader

func main() {

	//Creating .log-file for logging output from program, while still printing to the command line
	err := os.Remove("OutputLog.log")
	if err != nil {
	}
	f, err := os.OpenFile("OutputLog.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	mw := io.MultiWriter(os.Stdout, f)
	if err != nil {
		fmt.Println("log does not work")
	}
	defer f.Close()
	log.SetOutput(mw)

	//Making reader and waitGroup, and taking arguments for the id's of the peers.
	reader = bufio.NewReader(os.Stdin)
	var wg sync.WaitGroup
	wg.Add(1)
	arg1, _ := strconv.ParseInt(os.Args[1], 10, 32)
	arg2, _ := strconv.ParseInt(os.Args[2], 10, 32)
	ownPort := int32(arg1) + 5500
	totalPorts := int32(arg2)

	//Setting up context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Instantiating a peer
	p := &peer{
		id:         ownPort,
		wantsToken: false,
		nextPeer:   nil,
		ctx:        ctx,
	}

	// Creating listener on ownPort
	list, err := net.Listen("tcp", fmt.Sprintf(":%v", ownPort))
	if err != nil {
		log.Fatalf("Failed to listen on port: %v", err)
	}

	//Making grpc server
	grpcServer := grpc.NewServer()
	ring.RegisterRingServer(grpcServer, p)

	//Serving listener on ownPort
	go func() {
		if err := grpcServer.Serve(list); err != nil {
			log.Fatalf("failed to server %v", err)
		}

	}()

	//Making connection to next peer
	log.Println(p.id, ": Trying to dial:", (((ownPort + 1) % totalPorts) + 5500))
	conn, err := grpc.Dial(fmt.Sprintf(":%v", ((ownPort+1)%totalPorts)+5500), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect: %s", err)
	}
	defer conn.Close()
	p.nextPeer = ring.NewRingClient(conn)

	//Run function for printing and reading input from user
	go takeInput(p, &wg)

	//If the current peer has the lowest id/portNumber out of all peers, send a connectionVerification to the nextPeer
	if ownPort == 5500 {
		p.nextPeer.CheckConnection(p.ctx, &ring.ConnectionVerification{Id: p.id})
	}
	wg.Wait()
}

func takeInput(p *peer, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		//Waiting for the user to choose whether to enter the critical state
		log.Println(p.id, ": Do you want to enter the critical state? (type 'yes' or 'no')")
		inputMessage, _ := reader.ReadString('\n')
		inputMessage = strings.TrimSpace(inputMessage)
		log.Println("User (", p.id, "):", inputMessage)
		//If the peer wants to enter the critical state, the peer's 'wantsToken' is set to true, and the peer waits for it to be changed back.
		if inputMessage == "yes" {
			p.wantsToken = true
			log.Println(p.id, ": Okay! Wait for the token to be passed along to you.")
			for p.wantsToken {

			}
		}
	}
}

/*
This method sends a connectionVerification to the next peer, until the verification has reached all peers.
when this happens, the token is passed to the first peer
*/
func (p *peer) CheckConnection(ctx context.Context, msg *ring.ConnectionVerification) (*ring.EmptyMessage, error) {
	if msg.Id == p.id {
		p.nextPeer.PassToken(ctx, &ring.Token{})
	} else {
		for p.nextPeer == nil {

		}
		p.nextPeer.CheckConnection(p.ctx, &ring.ConnectionVerification{Id: msg.Id})
	}
	return &ring.EmptyMessage{}, nil
}

/*
If the peer does not want the token, it passes it on to the next peer.
If the peer does want the token, it prints that the peer is in the critical state, and waits for input from the user to leave the state.
*/
func (p *peer) PassToken(ctx context.Context, token *ring.Token) (*ring.EmptyMessage, error) {
	if !p.wantsToken {
		go p.nextPeer.PassToken(p.ctx, &ring.Token{})
	} else {
		log.Println(p.id, ": You are now in the critical state. Do you want to exit the state? (type 'yes' or 'no')")
		inputMessage, _ := reader.ReadString('\n')
		inputMessage = strings.TrimSpace(inputMessage)
		log.Println("User (", p.id, "):", inputMessage)
		//When the peer wants to leave the critical state, a message is printed to them, and the token is passed on to the next peer.
		if inputMessage == "yes" {
			log.Println(p.id, ": Okay! You have passed the token along")
			p.wantsToken = false
			go p.nextPeer.PassToken(p.ctx, &ring.Token{})
		}
	}
	return &ring.EmptyMessage{}, nil
}

type peer struct {
	ring.UnimplementedRingServer
	id         int32
	wantsToken bool
	nextPeer   ring.RingClient
	ctx        context.Context
}
