package videosrv

import (
	"fmt"
	"log"
	"strconv"

	"dev.project/BackEndCode/devserver/model/core"
)

const (
	CmdAction_BeginVideo int = 1
	CmdAction_StopVideo  int = 2
	CmdAction_PauseVideo int = 3
)

var (
	ffServer_IPAddr  string = "192.168.0.201"
	ffServer_UDPPort int    = 1514

	ffServer_RtmpPort = 1935
	ffServer_RtmpPath = "/vod2/"

	ffServer_HttpPath = "/vod/"
	ffServer_HttpPort = 80
)

// HTTP点播地址
func VOD_HttpPath() string {
	ip, port, path := ffServer_IPAddr, ffServer_HttpPort, ffServer_HttpPath
	str := fmt.Sprintf("http://%s:%d%s", ip, port, path)
	return str
}

// RTMP点播地址
func VOD_RtmpPath() string {
	ip, port, path := ffServer_IPAddr, ffServer_RtmpPort, ffServer_RtmpPath
	str := fmt.Sprintf("rtmp://%s:%d%s", ip, port, path)
	return str
}

// HTTP直播地址
func Live_HttpPath() string {
	return "http://localhost/hls/mystream"
}

// RTMP直播地址
func Live_RtmpPath() string {
	return "rtmp://localhost:1935/hls/mystream"
}

// 服务器UDP地址
func FFmpegServer_UDPAddr() (addr string) {
	ip, port := ffServer_IPAddr, ffServer_UDPPort
	addr = fmt.Sprintf("%s:%d", ip, port)
	return addr
}

func init() {
	var config core.Config
	config.InitConfig("./config.ini")
	node := "ffserver"
	ffServer_IPAddr = config.Read(node, "ipaddr")
	ffServer_UDPPort, _ = strconv.Atoi(config.Read(node, "udpport"))
	ffServer_HttpPath = config.Read(node, "httppath")
	ffServer_HttpPort, _ = strconv.Atoi(config.Read(node, "httpport"))

	//	fmt.Println(ffServer_IPAddr)
	//	fmt.Println(ffServer_UDPPort)
	//	fmt.Println(ffServer_HttpPath)
	//	fmt.Println(ffServer_HttpPort)
	log.Println("<<<<<<<<<<<\t	VOD_HttpPath: " + VOD_HttpPath())
	log.Println("<<<<<<<<<<<\t	FFmpegServer_UDPAddr: " + FFmpegServer_UDPAddr())
}
