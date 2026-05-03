package main

import (
	"flag"
	"fmt"
	"pkg/shutdown"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/gateway"
)

var configFile = flag.String("f", "etc/gateway.yaml", "gateway config")

func main() {
	flag.Parse()

	var c gateway.GatewayConf

	conf.MustLoad(*configFile, &c)

	server := gateway.MustNewServer(c)

	service, exit := shutdown.NewService(server)

	go func() {
		fmt.Println("Starting GATEWAY server ...")
		service.Start()
	}()

	<-exit

	service.Stop()
}
