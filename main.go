package main

import (
	"context"
	"hello/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)


func main (){

	l := log.New(os.Stdout, "Product-API", log.LstdFlags)

	ph := handlers.NewProduct(l)

	//Create a new serve and register the handlers
	sm := http.NewServeMux() 
	sm.Handle("/",ph)
	
	
	//create a new server
	s:= &http.Server{
		Addr: ":9090",					// configure th bind address
		Handler: sm,					// set the default handler
		ErrorLog: l,					// set the logger for the server
		IdleTimeout: 120*time.Second,	// max time to read request from the client
		ReadTimeout: 5*time.Second,		// max time to write response to the client
		WriteTimeout: 10*time.Second,	// max time fir cinnectiins using tcp keep-alive
	}

	// start the server 
	go func(){
		l.Println("Starting server on port 9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	
	sig := <- c
	l.Println("Recieved terminate, graceful shutdown", sig)
	
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(tc)

}
