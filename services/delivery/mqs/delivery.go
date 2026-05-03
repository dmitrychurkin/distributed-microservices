package main

import (
	"context"
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

	ctx, svcContext := context.Background(), svc.NewServiceContext(c)

	for _, mq := range []service.Service{
		kq.MustNewQueue(c.KqConsumerConf, logic.NewOrderCreated(ctx, svcContext)),
	} {
		serviceGroup.Add(mq)
	}

	serviceGroup.Add(server)

	go func() {
		fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

		serviceGroup.Start()
	}()

	graceService, exit := shutdown.NewService(server)

	<-exit

	graceService.Stop()
	serviceGroup.Stop()
}
