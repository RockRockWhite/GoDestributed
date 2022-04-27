package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

const ServerPort = ":3001"
const ServicesUrl = "http://localhost" + ServerPort + "/services"

type registry struct {
	registration []Registration
	mutex        *sync.RWMutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registration = append(r.registration, reg)
	r.mutex.Unlock()
	err := r.sendRequiredServices(reg)
	return err
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

func (r *registry) sendRequiredServices(registration Registration) error {
	// 读操作上锁
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var p patch
	for _, serviceReg := range r.registration {
		for _, requireService := range registration.RequiredServices {
			if requireService == serviceReg.ServiceName {
				p.Added = append(p.Added, patchEntry{
					ServiceName: serviceReg.ServiceName,
					Url:         serviceReg.ServiceUrl,
				})
			}
		}
	}
	// 把信息发给注册的服务
	err := r.sendPatch(registration.ServiceUpdateUrl, p)
	if err != nil {
		log.Println(err)
	}
	return nil
}

func (r *registry) sendPatch(url string, p patch) interface{} {
	d, err := json.Marshal(p)
	if err != nil {
		return err
	}
	res, err := http.Post(url, "application/json", bytes.NewBuffer(d))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("send patch with error code %s", res.Status)
	}

	return nil
}

var reg = &registry{
	registration: make([]Registration, 0),
	mutex:        &sync.RWMutex{},
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
