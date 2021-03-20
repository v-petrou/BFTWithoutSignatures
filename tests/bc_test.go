package tests

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
	"testing"
)

func TestBvBroadcast(t *testing.T) {
	args := os.Args[5:]
	if len(args) == 5 {
		id, _ := strconv.Atoi(args[0])
		n, _ := strconv.Atoi(args[1])
		clients, _ := strconv.Atoi(args[2])
		scenario, _ := strconv.Atoi(args[3])
		remote, _ := strconv.Atoi(args[4])

		initializeForTestBc(id, n, clients, scenario, remote)
	} else {
		log.Fatal("Arguments should be '<id> <n> <clients> <scenario> <remote>'")
	}

	/*** Start Testing ***/

	go modules.BvBroadcast(1, 0)

	go modules.BvBroadcast(2, 1)

	go modules.BvBroadcast(3, uint(variables.ID%2))

	/*** End Testing ***/

	done := make(chan interface{}) // To keep the server running
	<-done
}

func TestBConsensus(t *testing.T) {
	args := os.Args[5:]
	if len(args) == 5 {
		id, _ := strconv.Atoi(args[0])
		n, _ := strconv.Atoi(args[1])
		clients, _ := strconv.Atoi(args[2])
		scenario, _ := strconv.Atoi(args[3])
		remote, _ := strconv.Atoi(args[4])

		initializeForTestBc(id, n, clients, scenario, remote)
	} else {
		log.Fatal("Arguments should be '<id> <n> <clients> <scenario> <remote>'")
	}

	/*** Start Testing ***/

	go modules.BinaryConsensus(1, 0)

	go modules.BinaryConsensus(2, 1)

	go modules.BinaryConsensus(3, uint(variables.ID%2))

	/*** End Testing ***/

	done := make(chan interface{}) // To keep the server running
	<-done
}

// Initializes the environment for the test
func initializeForTestBc(id int, n int, clients int, scenario int, rem int) {
	variables.Initialize(id, n, clients, rem)
	logger.InitializeLogger("/home/vasilis/tests/out/", "/home/vasilis/tests/error/")

	if variables.Remote {
		config.InitializeIP()
	} else {
		config.InitializeLocal()
	}
	config.InitializeScenario(scenario)

	logger.OutLogger.Print(
		"ID:", variables.ID, " | N:", variables.N, " | F:", variables.F, " | Clients:",
		variables.Clients, " | Scenario:", config.Scenario, " | Remote:", variables.Remote, "\n\n",
	)

	threshenc.ReadKeys("/home/vasilis/keys/")

	messenger.InitializeMessenger()
	messenger.Subscribe()
	messenger.TransmitMessages()

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
