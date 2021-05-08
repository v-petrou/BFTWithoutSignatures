package config

import (
	"BFTWithoutSignatures/variables"
	"strconv"
)

var address = []string{
	"192.168.0.72",
}

var (
	// RepAddressesIP - Initialize the address of IP REP sockets
	RepAddressesIP map[int]string

	// ReqAddressesIP - Initialize the address of IP REQ sockets
	ReqAddressesIP map[int]string

	// ServerAddressesIP - Initialize the address of IP Server sockets
	ServerAddressesIP map[int]string

	// ResponseAddressesIP - Initialize the address of IP Response sockets
	ResponseAddressesIP map[int]string
)

// InitializeIP - Initializes system with ips.
func InitializeIP() {
	RepAddressesIP = make(map[int]string, variables.N)
	ReqAddressesIP = make(map[int]string, variables.N)
	ServerAddressesIP = make(map[int]string, variables.Clients)
	ResponseAddressesIP = make(map[int]string, variables.Clients)

	for i := 0; i < variables.N; i++ {
		ad := i % len(address)

		RepAddressesIP[i] = "tcp://*:" + strconv.Itoa(4000+(variables.ID*100)+i)
		ReqAddressesIP[i] = "tcp://" + address[ad] + ":" + strconv.Itoa(4000+(i*100)+variables.ID)
	}
	for i := 0; i < variables.Clients; i++ {
		ServerAddressesIP[i] = "tcp://*:" + strconv.Itoa(7000+(variables.ID*100)+i)
		ResponseAddressesIP[i] = "tcp://*:" + strconv.Itoa(10000+(variables.ID*100)+i)
	}
}

// GetRepAddress - Returns the IP REP address of the server with that id
func GetRepAddress(id int) string {
	return RepAddressesIP[id]
}

// GetReqAddress - Returns the IP REQ address of the server with that id
func GetReqAddress(id int) string {
	return ReqAddressesIP[id]
}

// GetServerAddress - Returns the IP Server address of the server with that id
func GetServerAddress(id int) string {
	return ServerAddressesIP[id]
}

// GetResponseAddress - Returns the IP Response address of the server with that id
func GetResponseAddress(id int) string {
	return ResponseAddressesIP[id]
}
