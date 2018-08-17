package camunda

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	_ "github.com/davecgh/go-spew/spew"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/proxy"
)

type CamundaConnection struct {
	Protocol     string
	Server       string
	Port         int
	BaseEndpoint string
	User         string
	Password     string
	ValidateSSL  bool
	Proxy        string
	ProxyIsSocks bool
	BaseURL      string
	Client       *http.Client
}

func NewCamunda(UseSSL bool, Server string, Port int, BaseEndpoint string, User string, Password string, ValidateSSL bool, Proxy string, ProxyIsSocks bool) *CamundaConnection {
	var camunda *CamundaConnection
	var tr *http.Transport

	camunda = new(CamundaConnection)
	tr = &http.Transport{
		DisableKeepAlives:   false,
		IdleConnTimeout:     0,
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 100,
	}
	if !ValidateSSL {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	if Proxy != "" {
		if ProxyIsSocks {
			dialer, err := proxy.SOCKS5("tcp", Proxy, nil, proxy.Direct)
			if err != nil {
				log.WithField("url", Proxy).Error("Can't connect to Socks5 proxy: " + err.Error())
			}
			tr.Dial = dialer.Dial
		} else {
			proxyUrl, err := url.Parse(Proxy)
			if err != nil {
				log.WithField("url", Proxy).Error(err)
			}
			tr.Proxy = http.ProxyURL(proxyUrl)
		}
	}
	camunda.Client = &http.Client{
		Transport: tr,
		Timeout:   time.Second * 10}
	if UseSSL {
		camunda.Protocol = "https"
	} else {
		camunda.Protocol = "http"
	}
	if BaseEndpoint == "" {
		BaseEndpoint = "/engine-rest"
	}
	camunda.BaseURL = camunda.Protocol + "://"
	camunda.Server = Server
	camunda.Port = Port
	camunda.BaseEndpoint = BaseEndpoint
	camunda.User = User
	camunda.Password = Password
	camunda.ValidateSSL = ValidateSSL
	camunda.Proxy = Proxy
	camunda.ProxyIsSocks = ProxyIsSocks
	if User != "" {
		camunda.BaseURL = camunda.BaseURL + User + ":" + Password + "@"
	}
	camunda.BaseURL = camunda.BaseURL + Server + ":" + strconv.Itoa(Port) + BaseEndpoint
	log.WithFields(log.Fields{
		"Protocol":     camunda.Protocol,
		"Server":       camunda.Server,
		"Port":         camunda.Port,
		"User":         camunda.User,
		"Password":     camunda.Password,
		"ValidateSSL":  camunda.ValidateSSL,
		"Proxy":        camunda.Proxy,
		"ProxyIsSocks": camunda.ProxyIsSocks,
		"BaseURL":      camunda.BaseURL,
	}).Debug("Camunda connection initialized")
	return camunda
}

func (camunda *CamundaConnection) Get(endpoint string) (map[string]interface{}, error) {
	var x []byte

	response, err := camunda.request("GET", endpoint, x)
	if err != nil {
		return nil, err
	}
	return toJSON(response)
}

func (camunda *CamundaConnection) GetRaw(endpoint string) ([]byte, error) {
	var x []byte
	return camunda.request("GET", endpoint, x)
}

func (camunda *CamundaConnection) Delete(endpoint string) (map[string]interface{}, error) {
	var x []byte

	response, err := camunda.request("DELETE", endpoint, x)
	if err != nil {
		return nil, err
	}
	return toJSON(response)
}

func (camunda *CamundaConnection) DeleteRaw(endpoint string) ([]byte, error) {
	var x []byte
	return camunda.request("DELETE", endpoint, x)
}

func (camunda *CamundaConnection) Post(endpoint string, jsonData []byte) (map[string]interface{}, error) {
	response, err := camunda.request("POST", endpoint, jsonData)
	if err != nil {
		return nil, err
	}
	return toJSON(response)
}

func (camunda *CamundaConnection) PostRaw(endpoint string, jsonData []byte) ([]byte, error) {
	return camunda.request("POST", endpoint, jsonData)
}

func (camunda *CamundaConnection) Put(endpoint string, jsonData []byte) (map[string]interface{}, error) {
	response, err := camunda.request("PUT", endpoint, jsonData)
	if err != nil {
		return nil, err
	}
	return toJSON(response)
}

func (camunda *CamundaConnection) PutRaw(endpoint string, jsonData []byte) ([]byte, error) {
	return camunda.request("PUT", endpoint, jsonData)
}

func (camunda *CamundaConnection) request(method string, endpoint string, jsonData []byte) ([]byte, error) {
	var req *http.Request
	var err error

	target := camunda.BaseURL + endpoint
	switch method {
	case "GET", "DELETE":
		req, err = http.NewRequest(method, target, nil)
	default:
		req, err = http.NewRequest(method, target, bytes.NewBuffer(jsonData))
	}

	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	r, err := camunda.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	response, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func toJSON(response []byte) (map[string]interface{}, error) {
	var data interface{}
	err := json.Unmarshal(response, &data)
	if err != nil {
		return nil, err
	}
	m := data.(map[string]interface{})
	return m, nil
}
