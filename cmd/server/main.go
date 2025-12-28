package main

import (
	"log"

	"smbsd.local/wwwrpc"
)

func main() {
	log.Println("start root handler")
	keep := make(wwwrpc.Keep)
	keep.Serve("127.0.0.1:4400")
}
