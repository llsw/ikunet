package main

import (
	transport "github.com/llsw/ikunnet/kitex_gen/transport/transportservice"
	"log"
)

func main() {
	svr := transport.NewServer(new(TransportServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
