package config

import "BFTWithoutSignatures/logger"

var (
	Scenario string

	scenarios = map[int]string{
		0: "NORMAL",      // Normal execution
		1: "IDLE",        // Byzantine processes remain idle (send nothing)
		2: "BC_ATTACK",   // Byzantine processes send wrong bytes for BC
		3: "HALF_&_HALF", // Byzantine processes send correct messages to half and empty to others
	}
)

func InitializeScenario(s int) {
	if s >= len(scenarios) {
		logger.ErrLogger.Println("Scenario out of bounds! Executing with NORMAL scenario ...")
		s = 0
	}

	Scenario = scenarios[s]
}
