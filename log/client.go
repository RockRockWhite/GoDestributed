package log

import (
	"GoDestributed/registry"
	"bytes"
	"fmt"
	stdlog "log"
	"net/http"
)

func SetClientLogger(serviceUrl string, clientService registry.ServiceName) {
	stdlog.SetPrefix(fmt.Sprintf("[%s] -", clientService))
	stdlog.SetFlags(0)
	stdlog.SetOutput(&clientLogger{serviceUrl})
}

type clientLogger struct {
	url string
}

func (c clientLogger) Write(p []byte) (n int, err error) {
	b := bytes.NewBuffer(p)
	res, err := http.Post(c.url+"/log", "text/plain", b)
	if err != nil {
		return 0, err
	}
	if res.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("logger: %s", res.Status)
	}

	return len(p), nil
}
