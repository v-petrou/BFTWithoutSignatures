package tests

// go test -v -run $TEST /home/vasilis/go/src/BFTWithoutSignatures/tests -args 0 4 1 1 0

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
	"time"
)

func TestRBroadcast(t *testing.T) {
	args := os.Args[5:]
	if len(args) == 5 {
		id, _ := strconv.Atoi(args[0])
		n, _ := strconv.Atoi(args[1])
		clients, _ := strconv.Atoi(args[2])
		scenario, _ := strconv.Atoi(args[3])
		remote, _ := strconv.Atoi(args[4])

		initializeForTestRb(id, n, clients, scenario, remote)
	} else {
		log.Fatal("Arguments should be '<id> <n> <clients> <scenario> <remote>'")
	}

	/*** Start Testing ***/

	go modules.ReliableBroadcast(1, "MVC", []byte("AEK"))

	if (variables.ID % 2) == 0 {
		go modules.ReliableBroadcast(2, "MVC", []byte("LFC"))
	} else {
		go modules.ReliableBroadcast(2, "MVC", []byte("lfc"))
	}

	go modules.ReliableBroadcast(3, "MVC", []byte("aek"))

	/*** End Testing ***/

	done := make(chan interface{}) // To keep the server running
	<-done
}

func TestAbcRBroadcast(t *testing.T) {
	args := os.Args[5:]
	if len(args) == 5 {
		id, _ := strconv.Atoi(args[0])
		n, _ := strconv.Atoi(args[1])
		clients, _ := strconv.Atoi(args[2])
		scenario, _ := strconv.Atoi(args[3])
		remote, _ := strconv.Atoi(args[4])

		initializeForTestRb(id, n, clients, scenario, remote)
	} else {
		log.Fatal("Arguments should be '<id> <n> <clients> <scenario> <remote>'")
	}

	/*** Start Testing ***/

	go modules.ReliableBroadcastAbc()

	time.Sleep(2 * time.Second)

	modules.SendRBInit(1, []byte("LFC"))

	if variables.ID == 0 {
		modules.SendRBInit(2, []byte("AEK"))
	}

	if variables.ID == 1 {
		modules.SendRBInit(2, []byte("TEST"))
	}

	/*** End Testing ***/

	done := make(chan interface{}) // To keep the server running
	<-done
}

// Initializes the environment for the test
func initializeForTestRb(id int, n int, clients int, scenario int, rem int) {
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
