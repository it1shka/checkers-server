package main

import (
	"flag"
	"fmt"
	"sync"

	"it1shka.com/checkers-server/multiplayer"
	"it1shka.com/checkers-server/testapp"
)

func main() {
	testRun := flag.Bool("test-run", false, "runs application in test mode")
	flag.Parse()
	if *testRun {
		testapp.RunTestApp()
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		fmt.Println("Starting up realtime server...")
		multiplayer.GetServer().Start()
	}()

	go func() {
		defer wg.Done()
		fmt.Println("Starting up API server...")
		// TODO: startup API server here
	}()

	wg.Wait()
}
