package main

import (
	"GoDestributed/log"
	"GoDestributed/registry"
	"GoDestributed/services"
	"context"
	"fmt"
	stdlog "log"
)

func main() {
	log.Run("./distributed.log")
	host, port := "localhost", "3002"

	ctx, err := services.Start(context.Background(), host, port, log.RegisterHandlers, registry.Registration{
		ServiceName: "logService",
		ServiceUrl:  fmt.Sprintf("http://%s:%s", host, port),
	})

	if err != nil {
		stdlog.Fatalln(err)
	}

	// 等待管道信号
	<-ctx.Done()
	fmt.Println("shutting down")

	err = registry.UnregisterService(fmt.Sprintf("http://%s:%s", host, port))
	if err != nil {
		stdlog.Println(err)
	}
}
