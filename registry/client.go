package registry

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RegisterService 注册服务
func RegisterService(r Registration) error {
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(r)
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
