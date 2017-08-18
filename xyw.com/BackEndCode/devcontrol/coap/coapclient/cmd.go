package coapclient

import (
	"fmt"

	"canopus"
)

// CoAP命令
type CoapCommand struct {
	HostAddr    string            // IP:Port
	Method      canopus.CoapCode  // GET/POST...
	RequestURI  string            // 资源路径
	Payload     string            // 负载
	QueryParams map[string]string // 查询参数
}

// CoAP URI
func (p *CoapCommand) URIString() string {
	str := fmt.Sprintf("coap://%s%s", p.HostAddr, p.RequestURI)
	index := 0
	for k, v := range p.QueryParams {
		if index == 0 {
			str += fmt.Sprintf("?%s=%s", k, v)
			index++
			continue
		}

		str += fmt.Sprintf("&%s=%s", k, v)
	}

	return str
}
