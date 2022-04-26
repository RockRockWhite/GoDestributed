package services

import (
	"GoDestributed/registry"
	"context"
	"fmt"
	"log"
	"net/http"
)

func Start(ctx context.Context, host, port string, registerHandlersFunc func(), reg registry.Registration) (context.Context, error) {
	registerHandlersFunc()

	ctx = startServices(ctx, host, port, reg)
	err := registry.RegisterService(reg)
	if err != nil {
		return ctx, err
	}
	return ctx, nil
}

func startServices(ctx context.Context, host, port string, reg registry.Registration) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	var srv http.Server
	srv.Addr = host + ":" + port

	go func() {
		log.Println(srv.ListenAndServe())
		cancel()
	}()

	go func() {
		fmt.Printf("%v started. Press any key to stop.\n", reg.ServiceName)
		var str string
		fmt.Scanln(&str)
		srv.Shutdown(ctx)
		cancel()
	}()

	return ctx
}
