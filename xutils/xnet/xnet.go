package xnet

import (
	"net"
	"time"
)

// 获取本机IP地址
func GetLocalIPAddr() (ip string) {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		// 检查ip地址判断是否回环地址
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip = ipnet.IP.String()
					break
				}
			}
		}
	}

	return ip
}

func GetOneLocalIPAddr() string {
	const remote_addr = "baidu.com"      //
	const ping_timeout = 5 * time.Second //

	ip_chan := make(chan string, 1)
	count_chan := make(chan int, 1)
	defer close(ip_chan)
	defer close(count_chan)

	count := 0 // 统计计数
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, address := range addrs {
			// 检查ip地址判断是否回环地址
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip := ipnet.IP.String()
					count++
					go func() {
						laddr := net.IPAddr{IP: net.ParseIP(ip)}
						raddr, _ := net.ResolveIPAddr("ip4", remote_addr)
						conn, err := net.DialIP("ip4:icmp", &laddr, raddr)
						if err == nil {
							defer conn.Close()
							ip_chan <- ip
						}

						count_chan <- 1
					}()
				}
			}
		}
	}

	var list = []string{}
	acc := 0
	for done := false; !done; {
		select {
		case str, ok := <-ip_chan:
			if ok {
				list = append(list, str)
			}
		case i := <-count_chan:
			if acc += i; acc >= count {
				println("Complete,,,,,,", acc, count)
				done = true
			}
		case <-time.After(ping_timeout):
			println("Timeout,,,,,,,")
			done = true
		}

	}

	var ip string
	if len(list) > 0 {
		ip = list[0]
	}

	return ip
}

func GetLocalIPAddrList() []string {
	const remote_addr = "baidu.com"      //
	const ping_timeout = 5 * time.Second //

	ip_chan := make(chan string, 1)
	count_chan := make(chan int, 1)
	defer close(ip_chan)
	defer close(count_chan)

	count := 0 // 统计计数
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, address := range addrs {
			// 检查ip地址判断是否回环地址
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip := ipnet.IP.String()
					count++
					go func() {
						laddr := net.IPAddr{IP: net.ParseIP(ip)}
						raddr, _ := net.ResolveIPAddr("ip4", remote_addr)
						conn, err := net.DialIP("ip4:icmp", &laddr, raddr)
						if err == nil {
							defer conn.Close()
							ip_chan <- ip
						}

						count_chan <- 1
					}()
				}
			}
		}
	}

	var iplist = []string{}
	acc := 0
	for done := false; !done; {
		// println("begin,,,,,,")
		select {
		case ip, ok := <-ip_chan:
			if ok {
				// println(ip)
				iplist = append(iplist, ip)
			}
		case i := <-count_chan:
			if acc += i; acc >= count {
				println("Complete,,,,,,", acc, count)
				done = true
			}
		case <-time.After(ping_timeout):
			println("Timeout,,,,,,,")
			done = true
		}

	}

	return iplist
}
