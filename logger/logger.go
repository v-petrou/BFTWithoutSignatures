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
func InitializeLogger() {
	outFolder := "./logs/"
	errFolder := "./logs/"
	t := time.Now().Format("01-02-2006 15:04:05")

	outFilePath := outFolder + "output_" + strconv.Itoa(variables.ID) + "_" + t + ".txt"
	errorFilePath := errFolder + "error_" + strconv.Itoa(variables.ID) + "_" + t + ".txt"

	outFile, err := os.OpenFile(outFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	errFile, err := os.OpenFile(errorFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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
