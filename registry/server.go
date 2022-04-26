package registry

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

const ServerPort = ":3001"
const ServicesUrl = "http://localhost" + ServerPort + "/services"

type registry struct {
	registration []Registration
	mutex        *sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registration = append(r.registration, reg)
	r.mutex.Unlock()
	return nil
}

func (r *registry) remove(url string) error {
	r.mutex.Lock()
	// 移除元素
	for i, each := range r.registration {
		if each.ServiceUrl == url {
			r.registration = append(r.registration[:i], r.registration[i+1:]...)
			break
		}
	}

	r.mutex.Unlock()
	return nil
}

var reg = &registry{
	registration: make([]Registration, 0),
	mutex:        &sync.Mutex{},
}

type RegistyService struct {
}

func (r RegistyService) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	log.Println("Request received")
	switch request.Method {
	case http.MethodPost:
		dec := json.NewDecoder(request.Body)
		var registration Registration
		err := dec.Decode(&registration)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Adding server: %s with url: %s \n", registration.ServiceName, registration.ServiceUrl)
		err = reg.add(registration)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	case http.MethodDelete:
		dec := json.NewDecoder(request.Body)
		var registration struct {
			ServiceUrl string
		}
		err := dec.Decode(&registration)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("Removing server: with url: %s \n", registration.ServiceUrl)
		err = reg.remove(registration.ServiceUrl)
		if err != nil {
			log.Println(err)
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
