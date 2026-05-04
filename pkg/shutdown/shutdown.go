package shutdown

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

type service[S Service] struct {
	instance *S
}

func NewService[S Service](instance *S) (*service[S], <-chan os.Signal) {
	exit := make(chan os.Signal, 1)

	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)

	return &service[S]{instance}, exit
}

func (self *service[S]) Start() {
	log.Println("Starting the service")

	(*self.instance).Start()
}

func (self *service[S]) Stop() {
	log.Println("Stopping the service")

	(*self.instance).Stop()
}
