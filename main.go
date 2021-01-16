package main

import (
	"BFTWithoutSignatures/app"
	"BFTWithoutSignatures/config"
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/variables"
	"log"
	"os"
	"os/signal"
	"strconv"
)

// Initializer - Method that initializes all required processes
func initializer(id int, n int, t int, clients int, scenario config.Scenario) {
	variables.Initialize(id, n, t, clients)

	config.InitializeLocal()
	config.InitializeIP()
	config.InitializeScenario(scenario)

	logger.InitializeLogger()
	logger.OutLogger.Println(
		"ID:", variables.ID,
		"| N:", variables.N,
		"| F:", variables.F,
		"| T:", variables.T,
		"| Clients:", variables.Clients,
	)

	app.InitializeMessenger()
	// app.InitializeReplication()

	// Initialize the message transmition and handling for the servers
	app.Subscribe()
	go app.TransmitMessages()
	// go app.ByzantineReplication()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	go func() {
		for range terminate {
			for i := 0; i < n; i++ {
				if i == id {
					continue
				}
				app.ReceiveSockets[i].Close()
				app.SendSockets[i].Close()
			}
			os.Exit(0)
		}
	}()
}

func main() {
	done := make(chan interface{})

	args := os.Args[1:]
	if len(args) < 5 {
		log.Fatal("Arguments should be '<id> <n> <f> <k> <scenario>")
	}

	id, _ := strconv.Atoi(args[0])
	n, _ := strconv.Atoi(args[1])
	t, _ := strconv.Atoi(args[2])
	clients, _ := strconv.Atoi(args[3])
	tmp, _ := strconv.Atoi(args[4])
	scenario := config.Scenario(tmp)

	initializer(id, n, t, clients, scenario)

	_ = <-done
}
