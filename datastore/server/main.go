package main

import (
	"flag"
	"log"
)

var (
	addr = flag.String("addr", ":3124", "Port to listen on")
	baseDir = flag.String("basedir", "data", "Directory to save data to")
)

func main() {
	flag.Parse()
	s, err := New(*addr, *baseDir)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Serve()
	if err != nil {
		log.Fatal(err)
	}
}
