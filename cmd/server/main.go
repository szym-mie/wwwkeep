package main

import (
	"log"

	"smbsd.local/wwwkeep"
)

func main() {
	log.Println("start root handler")
	keep := make(wwwkeep.Keep)
	keep.Serve("127.0.0.1:4400")
}
