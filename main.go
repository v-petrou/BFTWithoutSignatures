package main

import (
	messenger "BFTWithoutSignatures/app"
	"BFTWithoutSignatures/config"
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/variables"
	"SSBFT/app"
	"log"
	"os"
	"os/signal"
	"strconv"
)

// Initializer method
func initializer(id int, n int, t int, clients int, scenario config.Scenario) {
	variables.Initialize(id, n, t, clients)

	config.InitializeLocal()
	config.InitializeIP()
	config.InitializeScenario(scenario)

	logger.InitializeLogger()
	logger.OutLogger.Println(
		"N", variables.N,
		"ID", variables.ID,
		"F", variables.F,
		"Threshold T", variables.T,
		"Client Size", variables.Clients,
	)

	messenger.InitializeMessenger()

	app.InitializeAutomaton()
	app.InitializeViewChange()
	app.InitializeFailureDetector()
	app.InitializeEstablishment()
	app.InitializeReplication()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			for i := 0; i < n; i++ {
				if i == id {
					continue
				}
				messenger.ReceiveSockets[i].Close()
				messenger.SendSockets[i].Close()
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

	// Initialize the message transmition and handling for the servers
	messenger.Subscribe()
	go messenger.TransmitMessages()
	go app.ByzantineReplication()

	_ = <-done
}
