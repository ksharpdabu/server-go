package main

import (
	"fmt"
	"github.com/OpenIoTHub/server-go/config"
	"github.com/OpenIoTHub/server-go/session"
	"github.com/OpenIoTHub/server-go/udpapi"
	"time"
)

func main() {
	go session.RunTLS(config.ConfigMode.Common.TlsPort)
	go session.RunTCP(config.ConfigMode.Common.TcpPort)
	go session.RunKCP(config.ConfigMode.Common.KcpPort)
	go udpapi.RunApiServer(config.ConfigMode.Common.UdpApiPort)
	fmt.Println("服务器正在运行，内网端配置请根据本服务器配置填写！")
	for {
		time.Sleep(time.Hour * 99999)
	}
}
