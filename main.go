package main

import (
	messenger "BFTWithoutSignatures/app"
	"BFTWithoutSignatures/config"
	"BFTWithoutSignatures/logger"
	"SSBFT/app"
	"SSBFT/variables"
	"log"
	"os"
	"os/signal"
	"strconv"
)

// Initialize - Initializer method
func Initialize(id int, n int, t int, k int, scenario config.Scenario) {
	variables.Initialise(id, n, t, k)

	config.InitializeLocal(n)
	config.InitializeIp(n)
	config.InitializeScenario(scenario)

	logger.InitializeLogger()
	logger.OutLogger.Println(
		"N", variables.N,
		"ID", variables.Id,
		"F", variables.F,
		"Threshold T", variables.T,
		"Client Size", variables.K,
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
				messenger.RcvSockets[i].Close()
				messenger.SndSockets[i].Close()
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
	k, _ := strconv.Atoi(args[3])
	tmp, _ := strconv.Atoi(args[4])
	scenario := config.Scenario(tmp)

	Initialize(id, n, t, k, scenario)
	messenger.Subscribe()

	if config.TestCase != config.NON_SS {
		go app.FailDetector()
		go app.ViewChangeMonitor()
		go app.CoordinatingAutomaton()

	}
	go messenger.TransmitMessages()
	go app.ByzantineReplication()

	_ = <-done
}
