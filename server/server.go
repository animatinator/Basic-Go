package server

import (
	"github.com/animatinator/learning/data"
	"github.com/animatinator/learning/ipc"
	"fmt"
	"sync"
)

type Connection struct {
	In <-chan ipc.Message
	Out chan<- ipc.Message
}

type Server struct {
	Conns []Connection
	Db data.Database
	// TODO: Logging
}

func New() Server {
	return Server {Db:data.NewDb()}
}

func (s *Server) Connect(c <-chan ipc.Message) <-chan ipc.Message {
	o := make(chan ipc.Message, 1)
	s.Conns = append(s.Conns, Connection{c, o})
	return o
}

func handleReadResponse(rsp <-chan string, c Connection) {
	go func() {
		v := <-rsp
		msg := ipc.Message{Type: ipc.RESULT, Payload: v}
		c.Out <- msg
	}()
}

func handleWriteResponse(rsp <-chan struct{}, c Connection) {
	go func() {
		<-rsp
		msg := ipc.Message{Type: ipc.RESULT}
		c.Out <- msg
	}()
}

func (s *Server) serveConnection(c Connection, wg *sync.WaitGroup) {
	// TODO: Add a logger that we send a stream of events to.
	for msg := range c.In {
		fmt.Println("[server] Got a message of type:", msg.Type)
		switch msg.Type {
			case ipc.READ:
				res := s.Db.Read(msg.Key)
				handleReadResponse(res, c)
			case ipc.WRITE:
				res := s.Db.Write(msg.Key, msg.Payload)
				handleWriteResponse(res, c)
			default:
				// TODO
		}
	}

	wg.Done()
}

func (s *Server) Run(done chan<- struct{}) {
	fmt.Println("[server] Starting up")
	fmt.Println("[server]", len(s.Conns), "connected channels")
	
	var wg sync.WaitGroup
	wg.Add(len(s.Conns))
	
	for _, c := range s.Conns {
		go s.serveConnection(c, &wg)
	}
	
	wg.Wait()
	
	done <- struct{}{}
}