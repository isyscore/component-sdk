package isyscore

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type ComponentSDK struct {
	licInfo        *LicenseData
	compInfo       *ResultComponentRegister
	licHost        string
	licPort        int
	componentName  string
	componentKey   string
	IsValid        bool
	InvalidMessage string
	License        *ComponentLicense
	Producer       *ComponentProducer
}

func NewComponentSDKDef(compName string, compKey string) *ComponentSDK {
	sdk := new(ComponentSDK)
	sdk.licHost = "isc-license-service"
	sdk.licPort = 9013
	sdk.componentName = compName
	sdk.componentKey = compKey
	sdk.load()
	return sdk
}

func NewComponentSDK(compName string, compKey string, host string, port int) *ComponentSDK {
	sdk := new(ComponentSDK)
	sdk.licHost = host
	sdk.licPort = port
	sdk.componentName = compName
	sdk.componentKey = compKey
	sdk.load()
	return sdk
}

func (sdk *ComponentSDK) load() {
	paramComp := fmt.Sprintf("{\"compName\":\"%s\", \"compKey\":\"%s\"}", sdk.componentName, sdk.componentKey)
	urlLic := fmt.Sprintf("http://%s:%d/api/core/license/read", sdk.licHost, sdk.licPort)
	ret := httpGet(urlLic)
	sdk.licInfo = new(LicenseData)
	_ = json.Unmarshal([]byte(ret), sdk.licInfo)
	urlComp := fmt.Sprintf("http://license.isyscore.com:9990/api/license/cloud/component/one2?compName=%s&compKey=%s", sdk.componentName, sdk.componentKey)
	ret2 := httpGet(urlComp)
	sdk.compInfo = new(ResultComponentRegister)
	_ = json.Unmarshal([]byte(ret2), sdk.compInfo)
	if sdk.compInfo.Code == 0 {
		sdk.IsValid = false
		sdk.InvalidMessage = "无法顺利请求云端服务器"
	} else {
		if sdk.compInfo.Data == nil {
			sdk.IsValid = false
			if sdk.compInfo.Message == "" {
				sdk.InvalidMessage = "无法从云端获取数据"
			} else {
				sdk.InvalidMessage = sdk.compInfo.Message
			}
		} else {
			if sdk.compInfo.Code == 200 {
				if sdk.licInfo.LicenseCode != "" {
					if sdk.licInfo.Customer.EnterpriseName == sdk.compInfo.Data.ProducerCompany && sdk.licInfo.Customer.ContactName == sdk.compInfo.Data.ProducerContact {
						// 直接授权给自己
						sdk.IsValid = true
						sdk.InvalidMessage = ""
					} else {
						// 查授权
						urlCompValid := fmt.Sprintf("http://%s:%d/api/core/license/component/valid", sdk.licHost, sdk.licPort)
						ret3 := httpPost(urlCompValid, paramComp)
						r := new(ResultComponentLicensed)
						_ = json.Unmarshal([]byte(ret3), r)
						if r.Code == 200 {
							sdk.IsValid = true
							sdk.InvalidMessage = r.Message
						} else {
							sdk.IsValid = false
							if r.Message == "" {
								sdk.InvalidMessage = "未能获取到组件状态"
							} else {
								sdk.InvalidMessage = r.Message
							}
						}
					}
				} else {
					sdk.IsValid = false
					sdk.InvalidMessage = "OS 未授权"
				}
			} else {
				sdk.IsValid = false
				sdk.InvalidMessage = sdk.compInfo.Message
			}
		}
	}
	// load license
	urlCompLic := fmt.Sprintf("http://%s:%d/api/core/license/component/license", sdk.licHost, sdk.licPort)
	ret4 := httpPost(urlCompLic, paramComp)
	rLic := new(ResultComponentLicense)
	_ = json.Unmarshal([]byte(ret4), rLic)
	sdk.License = rLic.Data
	// load producer
	urlProducer := fmt.Sprintf("http://%s:%d/api/core/license/component/producer", sdk.licHost, sdk.licPort)
	ret5 := httpPost(urlProducer, paramComp)
	rProc := new(ResultComponentProducer)
	_ = json.Unmarshal([]byte(ret5), rProc)
	sdk.Producer = rProc.Data
}

func httpGet(url string) string {
	var defaultTransport http.RoundTripper = &http.Transport{
		Proxy: nil,
		DialContext: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          30,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   15 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		DisableCompression:    true,
	}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "curl/7.64.1")
	req.Header.Set("Accept", "*/*")

	client := &http.Client{Transport: defaultTransport}

	resp, err := client.Do(req) // http.Get(url)
	if err != nil {
		fmt.Printf("1 => %v, resp = %v \n", err, resp)
		return ""
	}
	if resp.StatusCode != 200 {
		fmt.Printf("2 => %d \n", resp.StatusCode)
		return ""
	}
	buf, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		fmt.Printf("3 => %v \n", err)
		return ""
	}
	return string(buf)
}

func httpPost(url string, json string) string {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(json)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	if resp.StatusCode != 200 {
		return ""
	}
	buf, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return ""
	}
	return string(buf)
}
