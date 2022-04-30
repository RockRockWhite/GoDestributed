package main

import (
	"GoDestributed/registry"
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {
	registry.SetupRegistryService()
	http.Handle("/services", &registry.RegistyService{})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var srv http.Server
	srv.Addr = registry.ServerPort

	go func() {
		log.Println("Registry serveice start", srv.ListenAndServe())
		cancel()
	}()

	go func() {
		log.Println("press any key to exit")
		var s string
		fmt.Scanln(&s)
		cancel()
	}()

	<-ctx.Done()
	log.Println("Stopping registry server")
}
