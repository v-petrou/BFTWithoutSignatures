package main

import (
	"BFTWithoutSignatures/config"
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/messenger"
	"BFTWithoutSignatures/modules"
	"BFTWithoutSignatures/threshenc"
	"BFTWithoutSignatures/variables"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

// Initializer - Method that initializes all required processes
func initializer(id int, n int, t int, clients int, scenario config.Scenario) {
	variables.Initialize(id, n, t, clients)

	if variables.Remote {
		config.InitializeIP()
	} else {
		config.InitializeLocal()
	}
	config.InitializeScenario(scenario)

	logger.InitializeLogger("./logs/out/", "./logs/error/")
	logger.OutLogger.Print(
		"ID:", variables.ID, " | N:", variables.N, " | F:", variables.F,
		" | T:", variables.T, " | Clients:", variables.Clients, "\n\n",
	)

	threshenc.ReadKeys("./keys/")

	messenger.InitializeMessenger()
	messenger.Subscribe()
	messenger.TransmitMessages()

	modules.InitiateAtomicBroadcast()
	modules.RequestHandler()

	cleanup()
}

func cleanup() {
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		for range terminate {
			for i := 0; i < variables.N; i++ {
				if i == variables.ID {
					continue // Not myself
				}
				messenger.ReceiveSockets[i].Close()
				messenger.SendSockets[i].Close()
			}

			for i := 0; i < variables.Clients; i++ {
				messenger.ServerSockets[i].Close()
				messenger.ResponseSockets[i].Close()
			}
			os.Exit(0)
		}
	}()
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && string(args[0]) == "generate_keys" {
		N, _ := strconv.Atoi(args[1])
		threshenc.GenerateKeys(N, "./keys/")

	} else if len(args) == 5 {
		id, _ := strconv.Atoi(args[0])
		n, _ := strconv.Atoi(args[1])
		t, _ := strconv.Atoi(args[2])
		clients, _ := strconv.Atoi(args[3])
		tmp, _ := strconv.Atoi(args[4])
		scenario := config.Scenario(tmp)

		initializer(id, n, t, clients, scenario)

		// To keep the server running
		done := make(chan interface{})
		<-done

	} else {
		log.Fatal("Arguments should be '<id> <n> <t> <clients> <scenario>")
	}
}
