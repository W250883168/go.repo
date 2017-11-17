package xnet

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"runtime"
	"time"

	"container/list"
	"encoding/binary"

	"go.repo/xutils/xdebug"
)

type tagICMP struct {
	Type        uint8
	Code        uint8
	Checksum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func ping(host string) {
	var icmp tagICMP
	var laddr = net.IPAddr{IP: net.ParseIP("0.0.0.0")}
	var raddr, _ = net.ResolveIPAddr("ip", host)

	conn, err := net.DialIP("ip:icmp", &laddr, raddr)
	xdebug.LogError(err)
	defer conn.Close()

	icmp.Type = 8
	icmp.Code = 0
	icmp.Checksum = 0
	icmp.Identifier = 0
	icmp.SequenceNum = 0

	var buffer bytes.Buffer
	binary.Write(&buffer, binary.BigEndian, icmp)
	// binary.Write(&buffer, binary.BigEndian, []byte("aaaaaaaa"))
	icmp.Checksum = checkSum(buffer.Bytes())
	buffer.Reset()
	binary.Write(&buffer, binary.BigEndian, icmp)
	// binary.Write(&buffer, binary.BigEndian, []byte("aaaaaaaa"))

	statistic := list.New()
	sended_packets := 0
	recv_buff := make([]byte, 1024)

	for i := 1000; i > 0; i-- {
		fmt.Printf("\n正在 Ping %s 具有 %d 字节的数据:\n", raddr.String(), buffer.Len())
		if _, err := conn.Write(buffer.Bytes()); err != nil {
			fmt.Println(err.Error())
			return
		}
		sended_packets++
		tStart := time.Now()

		conn.SetReadDeadline((time.Now().Add(time.Second * 5)))
		_, err := conn.Read(recv_buff)
		// log.Println(string(recv_buff))

		if err != nil {
			fmt.Println("请求超时")
			continue
		}

		tEnd := time.Now()
		dur := tEnd.Sub(tStart).Nanoseconds() / 1e6
		log.Printf("\t来自 %s 的回复: 时间 = %dms\n", raddr.String(), dur)
		statistic.PushBack(dur)

		time.Sleep(time.Millisecond * 400)
	}
}

func checkSum(data []byte) uint16 {
	var sum uint32
	var length int = len(data)
	var index int

	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[index])
	}

	sum += (sum >> 16)
	return uint16(^sum)
}

func init() {
	if ptr, _, line, ok := runtime.Caller(0); ok {
		fun := runtime.FuncForPC(ptr)
		str := fmt.Sprintf("初始化： %s /%d", fun.Name(), line)
		log.Println(str)
	}
}
