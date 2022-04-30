package main

import (
	"GoDestributed/grades"
	"GoDestributed/log"
	"GoDestributed/registry"
	"GoDestributed/services"
	"context"
	"fmt"
	stdlog "log"
)

func main() {
	host, port := "localhost", "3003"

	ctx, err := services.Start(context.Background(), host, port, grades.RegisterHandlers, registry.Registration{
		ServiceName: "GradesService",
		ServiceUrl:  fmt.Sprintf("http://%s:%s", host, port),
		RequiredServices: []registry.ServiceName{
			"LogService",
		},
		ServiceUpdateUrl: fmt.Sprintf("http://%s:%s", host, port) + "/services",
		HeartbeatUrl:     fmt.Sprintf("http://%s:%s", host, port) + "/heartbeat",
	})

	if err != nil {
		stdlog.Fatalln(err)
	}

	if url, err := registry.GetProvider(registry.LogSerice); err == nil {
		fmt.Printf("Logging service found at %s\n", url)
		log.SetClientLogger(url, "LogService")

		stdlog.Println("Logging service found at!!!", url)
	}

	// 等待管道信号
	<-ctx.Done()
	fmt.Println("shutting down")

	err = registry.UnregisterService(fmt.Sprintf("http://%s:%s", host, port))
	if err != nil {
		stdlog.Println(err)
	}
}
