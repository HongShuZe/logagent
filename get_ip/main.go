package main

import (
	"net"
	"strings"
	"fmt"
)

//获取本地对外ip
func GetOutboundIP() (ip string, err error) {
	//Dial()连接到指定网络上的地址。
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return
	}
	defer conn.Close()
	//LocalAddr()返回本地网络地址。
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	ip = strings.Split(localAddr.IP.String(), ":")[0]
	return
}

func main()  {
	ip, err := GetOutboundIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ip)
}
