package service

import (
	"crypto/tls"
	"encoding/json"
	"hidevops.io/hiboot/pkg/app"
	"hidevops.io/mio/pkg/utils"
	"io"
	"net/http"
	"time"
)

type ClientService interface {
	Get(method, baseUrl string, v interface{}) (*http.Response, error)
}

func init() {
	app.Register(newClientSet)
}

type ClientServiceImpl struct {
	ClientService
	client *http.Client
}

func newClientSet() ClientService {

	return &ClientServiceImpl{
		client: NewClientSet(),
	}
}

func NewClientSet() *http.Client {
	transport := &utils.Transport{
		ConnectTimeout: 1 * time.Second,
		RequestTimeout: 2 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	defer transport.Close()
	return &http.Client{Transport: transport}
}

func (c *ClientServiceImpl) Get(method, baseUrl string, v interface{}) (*http.Response, error) {
	req2, _ := http.NewRequest(method, baseUrl, nil)
	resp, err := c.client.Do(req2)
	if err != nil {
		return nil, err
	}
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}
	resp.Body.Close()
	return resp, err
}


func Client(method, baseUrl string, v interface{}) (*http.Response, error) {
	transport := &utils.Transport{
		ConnectTimeout: 1 * time.Second,
		RequestTimeout: 2 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	defer transport.Close()
	client := &http.Client{Transport: transport}
	req2, _ := http.NewRequest(method, baseUrl, nil)
	resp, err := client.Do(req2)
	if err != nil {
		return nil, err
	}
	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}
	resp.Body.Close()
	return resp, err
}
