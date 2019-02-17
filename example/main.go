package main

import (
	"fmt"

	"github.com/superp00t/go-web-daemons/service"
)

func main() {
	svc := service.New()

	svc.OnPort(func(port *service.Port) {
		fmt.Println("port opened")
	})

	svc.Run()
}
