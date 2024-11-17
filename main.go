package main

import (
	"flag"

	"it1shka.com/checkers-server/testapp"
)

func main() {
	testRun := flag.Bool("test-run", false, "runs application in test mode")
	flag.Parse()
	if *testRun {
		testapp.RunTestApp()
		return
	}

	// TODO: implement actual application
	panic("TODO: not implemented!")
}
