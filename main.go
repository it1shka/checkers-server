package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"it1shka.com/checkers-server/multiplayer"
	"it1shka.com/checkers-server/testapp"
)

const ENV_NAME_TESTAPP_PORT = "TESTAPP_PORT"
const ENV_NAME_MULTIPLAYER_PORT = "MULTIPLAYER_PORT"

const PORT_DEFAULT_TESTAPP = ":3030"
const PORT_DEFAULT_MULTIPLAYER = ":8080"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
		return
	}

	testRun := flag.Bool("test-run", false, "runs application in test mode")
	flag.Parse()

	if *testRun {
		testport := os.Getenv(ENV_NAME_TESTAPP_PORT)
		if testport == "" {
			fmt.Printf("%s is not set. Using the default one...\n", ENV_NAME_TESTAPP_PORT)
			testport = PORT_DEFAULT_TESTAPP
		}
		testapp.RunTestApp(testport)
		return
	}

	port := os.Getenv(ENV_NAME_MULTIPLAYER_PORT)
	if port == "" {
		fmt.Printf("%s is not set. Using the default one...\n", ENV_NAME_MULTIPLAYER_PORT)
		port = PORT_DEFAULT_MULTIPLAYER
	}
	server := multiplayer.NewServer()
	server.Start(port)
}
