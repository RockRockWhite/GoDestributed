package grades

import (
	"log"
	"net/http"
)

func RegisterHandlers() {
	http.HandleFunc("/grades", func(writer http.ResponseWriter, request *http.Request) {
		log.Println("Request received")
		switch request.Method {
		case http.MethodGet:
			writer.Write([]byte("GET Grades Service"))
		default:
			writer.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}
