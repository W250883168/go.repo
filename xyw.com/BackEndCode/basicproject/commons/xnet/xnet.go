package xnet

import "net"

// 获取本机IP地址
func GetLocalIPAddr() string {
	var ip string = ""
	if addrs, err := net.InterfaceAddrs(); err == nil {
		// 检查ip地址判断是否回环地址
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					ip = ipnet.IP.String()
				}
			}
		}
	}

	return ip
}
