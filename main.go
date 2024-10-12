package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	host := flag.String("host", "0.0.0.0", "host to bind to")
	port := flag.Int("port", 8080, "port to listen on")
	flag.Parse()

	addr := fmt.Sprintf("%s:%d", *host, *port)

	server := NewServer()

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(os.Stdout)

	if err := server.Start(addr); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
