package main

import (
	"net"
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		run()
	}()
	go func() {
		defer wg.Done()
		run()
	}()
	wg.Wait()
}

func run() {
	log.Info("reading file")
	f, err := os.ReadFile("./test/data/test.txt")
	if err != nil {
		log.Error(err)
		return
	}
	conn, err := net.Dial("tcp", ":4000")
	if err != nil {
		log.Error("unable to establish connection.")
	}
	_, err = conn.Write(f)
	if err != nil {
		log.Error(err)
	}
}
