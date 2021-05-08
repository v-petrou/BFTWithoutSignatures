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
	"time"
)

// Initializer - Method that initializes all required processes
func initializer(id int, n int, clients int, scenario int, rem int) {
	variables.Initialize(id, n, clients, rem)
	logger.InitializeLogger("./logs/out/", "./logs/error/")

	config.InitializeScenario(scenario)
	if variables.Remote {
		config.InitializeIP()
	} else {
		config.InitializeLocal()
	}

	logger.OutLogger.Print(
		"ID:", variables.ID, " | N:", variables.N, " | F:", variables.F, " | Clients:",
		variables.Clients, " | Scenario:", config.Scenario, " | Remote:", variables.Remote, "\n\n",
	)

	threshenc.ReadKeys("./keys/")

	messenger.InitializeMessenger()
	messenger.Subscribe()

	if (config.Scenario == "IDLE") && (variables.Byzantine) {
		logger.ErrLogger.Println(config.Scenario)
		return
	}

	messenger.TransmitMessages()
	modules.InitiateAtomicBroadcast()
	time.Sleep(2 * time.Second) // Wait 2s before start accepting requests to initiate all maps
	modules.RequestHandler()
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
			if (config.Scenario == "IDLE") && (variables.Byzantine) {
				logger.OutLogger.Printf("\n\nMessage Complexity: 0.00 msgs\nMessage Size: 0.000 MB\n\n")
			} else {
				if modules.Aid == 0 {
					logger.OutLogger.Printf("\n\nMessage Complexity: 0.00 msgs\nMessage Size: 0.000 MB\n\n")
				} else {
					logger.OutLogger.Printf(
						"\n\nMessage Complexity: %.2f msgs\nMessage Size: %.3f MB\n\n",
						float64(variables.MsgComplexity/modules.Aid),
						float64(float64(variables.MsgSize/int64(modules.Aid))/1000000))
				}
			}

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
		clients, _ := strconv.Atoi(args[2])
		scenario, _ := strconv.Atoi(args[3])
		remote, _ := strconv.Atoi(args[4])

		initializer(id, n, clients, scenario, remote)
		cleanup()

		done := make(chan interface{}) // To keep the server running
		<-done

	} else {
		log.Fatal("Arguments should be '<ID> <N> <Clients> <Scenario> <Remote>'")
	}
}
