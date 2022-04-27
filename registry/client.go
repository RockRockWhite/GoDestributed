package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
)

// RegisterService 注册服务
func RegisterService(r Registration) error {
	serviceUpdateUrl, err := url.Parse(r.ServiceUpdateUrl)
	if err != nil {
		return err
	}
	http.HandleFunc(serviceUpdateUrl.Path, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			dec := json.NewDecoder(r.Body)
			var p patch
			err := dec.Decode(&p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			prov.Update(p)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})

	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err = enc.Encode(r)
	if err != nil {
		return err
	}

	res, err := http.Post(ServicesUrl, "application/json", &buf)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Fatal to reigster service with status code: %d", res.StatusCode)
	}
	return nil
}

// UnregisterService 取消注册服务
func UnregisterService(url string) error {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(struct {
		ServiceUrl string
	}{url})
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, ServicesUrl, &buf)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	fmt.Printf("Unregisterring service from %s\n", ServicesUrl)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("Fatal to unregister service with status code: %d", res.StatusCode)
	}
	return nil
}

type providers struct {
	services map[ServiceName]string
	mutex    *sync.RWMutex
}

func (p *providers) Get(name ServiceName) (string, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	provider, ok := p.services[name]
	if !ok {
		return "", fmt.Errorf("No provider for service %s", name)
	}
	return provider, nil
}

func (p *providers) Update(pat patch) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, addEach := range pat.Added {
		p.services[addEach.ServiceName] = addEach.Url
	}

	for _, rmEach := range pat.Removed {
		for _, name := range p.services {
			if name == string(rmEach.ServiceName) {
				delete(p.services, rmEach.ServiceName)
			}
		}
		p.services[rmEach.ServiceName] = rmEach.Url
	}
}

var prov = providers{
	services: make(map[ServiceName]string),
	mutex:    &sync.RWMutex{},
}

func GetProvider(name ServiceName) (string, error) {
	return prov.Get(name)
}
