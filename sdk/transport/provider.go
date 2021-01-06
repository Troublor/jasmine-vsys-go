package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	sdkErr "github.com/Troublor/jasmine-vsys-go/sdk/error"
	"github.com/virtualeconomy/go-v-sdk/vsys"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"reflect"
	"strings"
)

type NetType = vsys.NetType

const Testnet = vsys.Testnet
const Mainnet = vsys.Mainnet

var Endpoint = map[NetType]string{
	Testnet: "http://klymena.vos.systems:9924",
	Mainnet: "http://vnode.vos.systems:9922",
}

type Provider struct {
	NetType  NetType
	Endpoint string
	Client   *http.Client
}

func NewProvider(endpoint string, netType NetType) (*Provider, error) {
	endpoint = strings.TrimSpace(endpoint)
	// remove the last '/' in the endpoint url
	for endpoint[len(endpoint)-1] == '/' {
		endpoint = endpoint[:len(endpoint)-1]
	}

	jar, err := cookiejar.New(nil)
	if err != nil { // error handling
		return nil, err
	}
	client := &http.Client{
		Jar: jar,
	}
	p := &Provider{Client: client, Endpoint: endpoint, NetType: netType}

	// call /node/status to check if the endpoint is working
	var nodeStatus NodeStatus
	err = p.Get("/node/status", nil, &nodeStatus)
	if err != nil {
		return nil, sdkErr.NewError(fmt.Sprintf("unable to get node status: %s", err.Error()))
	}
	if nodeStatus.BlockchainHeight <= 0 {
		return nil, sdkErr.NewError("node is not active")
	}

	vsys.InitApi(endpoint, netType)
	return p, nil
}

func (p *Provider) genBody(v reflect.Value) (io.Reader, error) {
	switch v.Kind() {
	case reflect.String:
		return bytes.NewBuffer([]byte(v.String())), nil
	case reflect.Ptr:
		return p.genBody(v.Elem())
	default:
		d, err := json.Marshal(v.Interface())
		if err != nil {
			return nil, err
		}
		reader := bytes.NewBuffer(d)
		return reader, nil
	}
}

func (p *Provider) prepareRequest(method string, path string, queries map[string]string, reqData interface{}) (*http.Request, error) {
	if len(path) > 0 {
		if path[0] != '/' {
			path = "/" + path
		}
	}

	var body io.Reader
	if reqData != nil {
		var err error
		body, err = p.genBody(reflect.ValueOf(reqData))
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, p.Endpoint+path, body)
	if err != nil {
		return nil, err
	}

	if queries != nil {
		q := req.URL.Query()
		for k, v := range queries {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (p *Provider) SendRequest(method, path string, queries map[string]string, reqData, respData interface{}) error {
	req, err := p.prepareRequest(method, path, queries, reqData)
	if err != nil {
		return err
	}
	resp, err := p.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode/100 == 2 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return json.Unmarshal(body, respData)
	} else if resp.StatusCode/100 == 4 {
		return NewVsysError(resp)
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return sdkErr.NewError(fmt.Sprintf("http error: %s, %s", resp.Status, string(body)))
	}

}

func (p *Provider) Get(path string, queries map[string]string, respData interface{}) error {
	return p.SendRequest("GET", path, queries, nil, respData)
}

func (p *Provider) Post(path string, queries map[string]string, reqData, respData interface{}) error {
	return p.SendRequest("POST", path, queries, reqData, respData)
}
