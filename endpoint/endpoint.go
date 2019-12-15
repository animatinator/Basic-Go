package endpoint

import (
	"github.com/animatinator/learning/ipc"
	"github.com/animatinator/learning/server"
)

type endpoint struct {
	// Send channel
	S chan<- ipc.Message
	// Results channel
	R <-chan ipc.Message
}

func New(svr *server.Server) endpoint {
	s := make(chan ipc.Message)
	r := svr.Connect(s)
	return endpoint {s, r}
}

func (e endpoint) Read(key string) <-chan string {
	m := ipc.Message{Type: ipc.READ, Key: key}
	e.S <- m
	
	res := make(chan string, 1)
	
	go func() {
		ans := <-e.R
		res <- ans.Payload
	}()
	
	return res
}

func (e endpoint) Write(key, value string) <-chan struct{} {
	m := ipc.Message{ipc.WRITE, key, value}
	e.S <- m
	
	res := make(chan struct{}, 1)
	
	go func() {
		<-e.R
		res <- struct{}{}
	}()
	
	return res
}

func (e endpoint) Disconnect() {
	close(e.S)
}