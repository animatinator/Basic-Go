package server

import (
	"github.com/animatinator/learning/data"
	"github.com/animatinator/learning/ipc"
	"fmt"
	"sync"
)

type Server struct {
	Inputs [](<-chan ipc.Message)
	Outputs [](chan<- ipc.Message)
	Db data.Database
	// TODO: Logging
}

func New() Server {
	return Server {Db:data.NewDb()}
}

func (s *Server) Connect(c <-chan ipc.Message) <-chan ipc.Message {
	s.Inputs = append(s.Inputs, c)
	o := make(chan ipc.Message, 1)
	s.Outputs = append(s.Outputs, o)
	return o
}

func (s *Server) handleReadResponse(rsp <-chan string, id int) {
	go func() {
		v := <-rsp
		msg := ipc.Message{Type: ipc.RESULT, Payload: v}
		s.Outputs[id] <- msg
	}()
}

func (s *Server) handleWriteResponse(rsp <-chan struct{}, id int) {
	go func() {
		<-rsp
		msg := ipc.Message{Type: ipc.RESULT}
		s.Outputs[id] <- msg
	}()
}

func (s *Server) Run(done chan<- struct{}) {
	fmt.Println("[server] Starting up")
	fmt.Println("[server]", len(s.Inputs), "connected channels")
	
	var wg sync.WaitGroup
	wg.Add(len(s.Inputs))
	
	for i, _ := range s.Inputs {
		fmt.Println(i)
		go func(i int) {
			fmt.Println("[server] Listening on channel", i)
			// TODO: Add a logger that we send a stream of events to.
			for msg := range s.Inputs[i] {
				fmt.Println("Got a message of type:", msg.Type)
				switch msg.Type {
					case ipc.READ:
						res := s.Db.Read(msg.Key)
						s.handleReadResponse(res, i)
					case ipc.WRITE:
						res := s.Db.Write(msg.Key, msg.Payload)
						s.handleWriteResponse(res, i)
					default:
						// TODO
				}
			}

			wg.Done()
		}(i)
	}
	
	wg.Wait()
	
	done <- struct{}{}
}