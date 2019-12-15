package main

import (
	"fmt"
	"github.com/animatinator/learning/endpoint"
	"github.com/animatinator/learning/server"
)

func main() {
	fmt.Println("[main] Creating and starting server...")
	s := server.New()
	
	// TODO: Wrap endpoints in CLients that perform some random behaviour
	e1 := endpoint.New(&s)
	e2 := endpoint.New(&s)
	
	done := make(chan struct{})
	go s.Run(done)
	
	// TODO: Remove this manual testing once CLients have been added
	e1.Write("hello", "world")
	e2.Write("whats", "up")
	fmt.Println("[main] Result after read: ", <- e1.Read("hello"))
	fmt.Println("[main] Result after read: ", <- e2.Read("whats"))
	
	e1.Disconnect()
	e2.Disconnect()
	
	<-done
}