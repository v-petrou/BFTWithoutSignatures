package main

import (
	"BFTWithoutSignatures/config"
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/messenger"
	"BFTWithoutSignatures/threshenc"
	"BFTWithoutSignatures/types"
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

	messenger.InitializeMessenger()
	// app.InitializeReplication()

	// Initialize the message transmition and handling for the servers
	messenger.Subscribe()
	go messenger.TransmitMessages()
	// go app.ByzantineReplication()

	threshenc.ReadKeys()

	// Start Testing
	if variables.ID == 1 {
		messenger.SendMessage(types.NewMessage([]byte("TEST"), "Test"), 2)

		logger.OutLogger.Println(threshenc.SecretKey)
		logger.OutLogger.Println(threshenc.VerificationKeys[0])
	}
	// End Testing

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt)
	go func() {
		for range terminate {
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
	if len(args) == 2 && string(args[0]) == "generate_keys" {
		N, _ := strconv.Atoi(args[1])
		threshenc.GenerateKeys(N)

	} else if len(args) == 5 {
		id, _ := strconv.Atoi(args[0])
		n, _ := strconv.Atoi(args[1])
		t, _ := strconv.Atoi(args[2])
		clients, _ := strconv.Atoi(args[3])
		tmp, _ := strconv.Atoi(args[4])
		scenario := config.Scenario(tmp)

		initializer(id, n, t, clients, scenario)

		_ = <-done

	} else {
		log.Fatal("Arguments should be '<id> <n> <f> <k> <scenario>")
	}
}
