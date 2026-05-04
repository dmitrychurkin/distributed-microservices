package main

import (
	"flag"
	"fmt"
	"pkg/env"
	"pkg/shutdown"

	"delivery/mqs/internal/config"
	"delivery/mqs/internal/logic"
	"delivery/mqs/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"

	_ "github.com/lib/pq"
)

var configFile = flag.String("f", "etc/delivery.yaml", "the config file")

func main() {
	env.LoadEnvWithExit("../.env")

	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	server, serviceGroup := rest.MustNewServer(c.RestConf), service.NewServiceGroup()

	svcContext := svc.NewServiceContext(c)

	for _, mq := range []service.Service{
		kq.MustNewQueue(c.KqConsumerConf, logic.NewOrderCreated(svcContext)),
	} {
		serviceGroup.Add(mq)
	}

	serviceGroup.Add(server)

	group, exit := shutdown.NewService(&serviceGroup)

	go func() {
		fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

		group.Start()
	}()

	<-exit

	group.Stop()
}
