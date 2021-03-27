package config

import "BFTWithoutSignatures/logger"

var (
	Scenario string

	scenarios = map[int]string{
		0: "NORMAL",      // Normal execution
		1: "IDLE",        // Byzantine processes remain idle (send nothing)
		2: "BC_ATTACK0",  // Byzantine processes only send 0 in BC
		3: "HALF_N_HALF", // Byzantine processes send correct messages to half and empty to others
	}
)

func InitializeScenario(s int) {
	if s >= len(scenarios) {
		logger.ErrLogger.Println("Scenario out of bounds! Executing with NORMAL scenario ...")
		s = 0
	}

	Scenario = scenarios[s]
}
