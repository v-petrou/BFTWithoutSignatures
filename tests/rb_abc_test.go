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
	"testing"
)

func TestAbcRBroadcast(t *testing.T) {
	args := os.Args[5:]
	if len(args) == 5 {
		id, _ := strconv.Atoi(args[0])
		n, _ := strconv.Atoi(args[1])
		t, _ := strconv.Atoi(args[2])
		clients, _ := strconv.Atoi(args[3])
		tmp, _ := strconv.Atoi(args[4])
		scenario := config.Scenario(tmp)

		initializeForTestAbcRb(id, n, t, clients, scenario)
	} else {
		log.Fatal("Arguments should be '<id> <n> <f> <k> <scenario>")
	}

	/*** Start Testing ***/

	go modules.ReliableBroadcastAbc()

	modules.SendRBInit(1, []byte("LFC"))

	if variables.ID == 0 {
		modules.SendRBInit(2, []byte("AEK"))
	}

	if variables.ID == 1 {
		modules.SendRBInit(2, []byte("TEST"))
	}

	/*** End Testing ***/

	done := make(chan interface{})
	_ = <-done
}

// Initializes the environment for the test
func initializeForTestAbcRb(id int, n int, t int, clients int, scenario config.Scenario) {
	variables.Initialize(id, n, t, clients)

	if variables.Remote {
		config.InitializeIP()
	} else {
		config.InitializeLocal()
	}
	config.InitializeScenario(scenario)

	logger.InitializeLogger("/home/vasilis/tests/out/", "/home/vasilis/tests/error/")
	logger.OutLogger.Print(
		"ID:", variables.ID, " | N:", variables.N, " | F:", variables.F,
		" | T:", variables.T, " | Clients:", variables.Clients, "\n\n",
	)

	threshenc.ReadKeys("/home/vasilis/keys/")

	messenger.InitializeMessenger()
	messenger.Subscribe()
	go messenger.TransmitMessages()

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
