package logger

import (
	"BFTWithoutSignatures/variables"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	// OutLogger - Log the outputs
	OutLogger *log.Logger

	// ErrLogger - Log the errors
	ErrLogger *log.Logger
)

// InitializeLogger - Initializes the Out and Err loggers
func InitializeLogger(outFolder string, errFolder string) {
	t := time.Now().Format("01-02-2006 15:04:05")

	outFilePath := outFolder + strconv.Itoa(variables.ID) + "_output_" + t + ".txt"
	errFilePath := errFolder + strconv.Itoa(variables.ID) + "_error_" + t + ".txt"

	outFile, err := os.OpenFile(outFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	errFile, err := os.OpenFile(errFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	OutLogger = log.New(
		outFile,
		"INFO:\t",
		log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	ErrLogger = log.New(
		errFile,
		"ERROR:\t",
		log.LstdFlags|log.Lmicroseconds|log.Lshortfile)
}
