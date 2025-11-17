package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"net"
)

const (
	TEXT_SPLIT = byte('\n')
)

func main() {
	listner, err := net.Listen("tcp", ":42069")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listner.Accept()
		if err != nil {
			break
		}
		request, err := request.RequestFromReader(conn)
		if err != nil {
			panic(err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", request.RequestLine.Method)
		fmt.Printf("- Target: %s\n", request.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", request.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		for key, value := range request.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}

	}
	listner.Close()
}
