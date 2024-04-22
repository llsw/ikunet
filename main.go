package main

import (
	"log"

	transport "github.com/llsw/ikunet/internal/kitex_gen/transport/transportservice"
)

func main() {
	svr := transport.NewServer(new(TransportServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
