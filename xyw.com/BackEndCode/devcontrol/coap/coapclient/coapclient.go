package coapclient

import (
	"errors"
	"fmt"
	"log"
	"net"
	"runtime"
	"sync/atomic"
	"time"

	"canopus"

	"xutils/xdebug"
	"xutils/xerr"

	"dev.project/BackEndCode/devcontrol/app"
)

// 发送CoAP命令
func Send(cmd CoapCommand) (response canopus.CoapResponse) {
	addr, _ := net.ResolveUDPAddr("udp4", cmd.HostAddr)
	log.Printf("%+v\n", addr.String())
	uri := cmd.RequestURI
	client := canopus.NewCoapClient()
	client.OnError(onError)
	client.OnMessage(onMessage)

	client.Dial(addr.String())
	client.OnStart(func(server canopus.CoapServer) {
		defer xerr.CatchPanic()
		defer func() {
			// 在新goroutine中关闭
			go client.Stop()
		}()

		request := canopus.NewRequest(canopus.MessageConfirmable, cmd.Method, canopus.GenerateMessageID())
		request.SetStringPayload(cmd.Payload)
		request.SetRequestURI(uri)
		for k, v := range cmd.QueryParams {
			request.SetURIQuery(k, v)
		}

		ack_chan := make(chan canopus.CoapResponse) // 应答通道
		defer close(ack_chan)
		var count = int32(0) // 重发计数
		proc := func() {
			atomic.AddInt32(&count, 1)
			txt := fmt.Sprintf("消息发送，第%d次: \t; MessageID=%d", count, request.GetMessage().MessageID)
			log.Println(txt)
			log.Printf("	URI: %s\n\t%+v\n", cmd.URIString(), cmd)
			resp, err := client.Send(request)
			xdebug.LogError(err)
			if err == nil {
				ack_chan <- resp
			}
		}

		go proc()
		var ackTimeout = time.Duration(app.GetConfig().CoapConfig.AckTimeout)
		ticker := time.NewTicker(time.Second * ackTimeout)
		defer ticker.Stop()
		for done := false; !done; {
			select {
			case <-ticker.C:
				if done = (int(count) >= app.GetConfig().CoapConfig.MaxSendCount); !done {
					log.Println("超时重发，，，，，，")
					go proc()
				}
			case response = <-ack_chan:
				// log.Printf("收到消息！ MessageID=%d\n", response.GetMessage().MessageID)
				done = true
			}
		}
	})

	client.Start()
	return response
}

// 发送CoAP命令
func Send2(cmd CoapCommand) (response canopus.CoapResponse, resend int, err error) {
	addr, _ := net.ResolveUDPAddr("udp4", cmd.HostAddr)
	log.Printf("%+v\n", addr.String())
	uri := cmd.RequestURI
	client := canopus.NewCoapClient()
	client.OnError(onError)
	client.OnMessage(onMessage)

	client.Dial(addr.String())
	client.OnStart(func(server canopus.CoapServer) {
		defer xerr.CatchPanic()
		defer func() {
			// 在新goroutine中关闭
			go client.Stop()
		}()

		request := canopus.NewRequest(canopus.MessageConfirmable, cmd.Method, canopus.GenerateMessageID())
		request.SetStringPayload(cmd.Payload)
		request.SetRequestURI(uri)
		for k, v := range cmd.QueryParams {
			request.SetURIQuery(k, v)
		}

		ack_chan := make(chan canopus.CoapResponse) // 应答通道
		defer close(ack_chan)
		var count = int32(0) // 重发计数
		proc := func() {
			atomic.AddInt32(&count, 1)
			resend = int(count)
			txt := fmt.Sprintf("消息发送，第%d次: \t; MessageID=%d", count, request.GetMessage().MessageID)
			log.Println(txt)
			log.Printf("	URI: %s\n\t%+v\n", cmd.URIString(), cmd)
			resp, err := client.Send(request)
			xdebug.LogError(err)
			if err == nil {
				ack_chan <- resp
			}
		}

		go proc()
		var ackTimeout = time.Duration(app.GetConfig().CoapConfig.AckTimeout)
		ticker := time.NewTicker(time.Second * ackTimeout)
		defer ticker.Stop()
		for done := false; !done; {
			select {
			case <-ticker.C:
				if done = (int(count) >= app.GetConfig().CoapConfig.MaxSendCount); !done {
					log.Println("超时重发，，，，，，")
					go proc()
				} else {
					err = errors.New(fmt.Sprintf("第%d次重发，响应超时", count))
				}
			case response = <-ack_chan:
				// log.Printf("收到消息！ MessageID=%d\n", response.GetMessage().MessageID)
				done = true
				err = nil
			}
		}
	})

	client.Start()
	return response, resend, err
}

func onMessage(msg *canopus.Message, inbound bool) {
	if inbound {
		log.Printf("收到消息！ MessageID=%d\n", msg.MessageID)
	}
}

func onError(err error) {
	xdebug.LogError(err)
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}
}
