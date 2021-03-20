package config

import "BFTWithoutSignatures/logger"

var (
	Scenario string

	scenarios = map[int]string{
		0: "NORMAL",
		1: "A",
		2: "B",
		3: "C",
	}
)

func InitializeScenario(s int) {
	if s >= len(scenarios) {
		logger.ErrLogger.Println("Scenario out of bounds! Running with NORMAL scenario ...")
		s = 0
	}

	Scenario = scenarios[s]
}
