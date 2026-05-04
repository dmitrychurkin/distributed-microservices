package main

import (
	"flag"
	"fmt"
	"pkg/env"
	"pkg/shutdown"

	"ordering/rpc/internal/config"
	"ordering/rpc/internal/server"
	"ordering/rpc/internal/svc"
	"ordering/rpc/ordering"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	_ "github.com/lib/pq"
)

var configFile = flag.String("f", "etc/ordering.yaml", "the config file")

func main() {
	env.LoadEnvWithExit("../.env")

	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	ctx := svc.NewServiceContext(c)

	server := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		ordering.RegisterOrderingServer(grpcServer, server.NewOrderingServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	service, exit := shutdown.NewService(&server)

	go func() {
		fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
		service.Start()
	}()

	<-exit

	service.Stop()
}
