package messenger

import (
	"SSBFT/config"
	"SSBFT/logger"
	"SSBFT/types"
	"SSBFT/variables"
	"bytes"
	"encoding/gob"
	"log"
	"strings"

	"github.com/pebbe/zmq4"
)

var (
	// Context to initialize sockets
	Context *zmq4.Context

	// RcvSockets -Receive messages from other
	RcvSockets map[int]*zmq4.Socket

	// SndSockets - Send messages to other servers
	SndSockets map[int]*zmq4.Socket

	// ServerSockets - Get the client requests
	ServerSockets map[int]*zmq4.Socket

	// ResponseSockets - Send responses to clients
	ResponseSockets map[int]*zmq4.Socket

	SendRecvSync map[int]chan interface{} // Probably not needed

	//requestsChan map[int]chan interface{}

	messageChan = make(chan struct {
		Message types.Message
		To      int
	})

	CoordChan = make(chan struct {
		Message *types.CoordinationMessage
		From    int
	}, 100)

	VcmChan = make(chan struct {
		Vcm  types.VCM
		From int
	}, 100)

	TokenChan = make(chan struct {
		Token types.Token
		From  int
	}, 100)

	RequestChan = make(chan *types.ClientMessage, 100)

	ReplicaChan = make(chan struct {
		Rep  *types.ReplicaStructure
		From int
	}, 100)

	count = 0
)

// InitializeMessenger - Initializes the 0MQ sockets for the servers to communicate with clients and each other
func InitializeMessenger() {
	SendRecvSync = make(map[int]chan interface{}, variables.K) // Probably not needed, cause its only used here
	for i := 0; i < variables.K; i++ {
		SendRecvSync[i] = make(chan interface{})
	}

	//requestsChan = make(map[int]chan interface{}, variables.K)
	//for i := 0; i < variables.K; i++ {
	//	requestsChan[i] = make(chan interface{}, 1)
	//}

	Context, err := zmq4.NewContext()
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	// Sockets REP & PUB to communicate with each one of the clients
	ServerSockets = make(map[int]*zmq4.Socket, variables.K)
	ResponseSockets = make(map[int]*zmq4.Socket, variables.K)
	for i := 0; i < variables.K; i++ {
		// ServerSockets initialization to get clients requests
		ServerSockets[i], err = Context.NewSocket(zmq4.REP)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		var serverAddr string
		if !variables.Remote {
			serverAddr = config.GetServerAddressLocal(i)
		} else {
			serverAddr = config.GetServerAddress(i)
		}
		err = ServerSockets[i].Bind(serverAddr)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		// ResponseSockets initialization to publish the response back to the clients.
		ResponseSockets[i], err = Context.NewSocket(zmq4.PUB)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		var responseAddr string
		if !variables.Remote {
			responseAddr = config.GetResponseAddressLocal(i)
		} else {
			responseAddr = config.GetResponseAddress(i)
		}
		err = ResponseSockets[i].Bind(responseAddr)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
	}

	// A socket pair (REP/REQ) to communicate with each one of the other servers
	RcvSockets = make(map[int]*zmq4.Socket)
	SndSockets = make(map[int]*zmq4.Socket)
	for i := 0; i < variables.N; i++ {
		// Not myself
		if i == variables.Id {
			continue
		}

		RcvSockets[i], err = Context.NewSocket(zmq4.REP)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		var rcvAddr string
		if !variables.Remote {
			rcvAddr = strings.Replace(
				config.GetRepAddressLocal(i),
				"localhost",
				"*",
				1)
		} else {
			rcvAddr = config.GetRepAddress(i)
		}
		err = RcvSockets[i].Bind(rcvAddr)
		if err != nil {
			logger.ErrLogger.Fatal(err, " "+rcvAddr)
		}
		logger.OutLogger.Println("Binded on ", rcvAddr)

		var sndAddr string
		if !variables.Remote {
			sndAddr = config.GetReqAddressLocal(i)
		} else {
			sndAddr = config.GetReqAddress(i)
		}
		SndSockets[i], err = Context.NewSocket(zmq4.REQ)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		err = SndSockets[i].Connect(sndAddr)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		logger.OutLogger.Println("Connected to ", sndAddr)
	}
}

// Broadcast - Calls SendMessage
func Broadcast(message types.Message) {
	for i := 0; i < variables.N; i++ {
		// Not myself
		if i == variables.Id {
			continue
		}
		SendMessage(message, i)
	}
}

// TransmitMessages - Echos the message to the other servers
func TransmitMessages() {
	for messageTo := range messageChan {
		to := messageTo.To
		message := messageTo.Message
		w := new(bytes.Buffer)
		encoder := gob.NewEncoder(w)
		err := encoder.Encode(message)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		_, err = SndSockets[to].SendBytes(w.Bytes(), 0)
		//logger.OutLogger.Println("SENT Message to ", to)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		_, err = SndSockets[to].Recv(0)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		//logger.OutLogger.Println("OKAY FROM ", to)
	}
}

// SendMessage - Puts the messages in the message channel to be transmited with TransmitMessages
func SendMessage(message types.Message, to int) {
	messageChan <- struct {
		Message types.Message
		To      int
	}{Message: message, To: to}
}

// SendReplica - (I think) Creates the replicas
func SendReplica(replica *types.ReplicaStructure, to int) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(replica)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	message := types.NewMessage(w.Bytes(), "ReplicaStructure")
	SendMessage(message, to)
}

// Subscribe - Handles the inputs from clients and other servers
func Subscribe() {
	// Gets requests from clients and handles them by putting them in RequestChan
	for i := 0; i < variables.K; i++ {
		go func(i int) {
			// Initialize them with a goroutine and waits forever
			for {
				message, err := ServerSockets[i].RecvBytes(0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}
				logger.OutLogger.Println("Request Received")
				handleRequest(message)
				_, err = ServerSockets[i].Send("", 0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}
				//_ = <-SendRecvSync[i]
				//requestsChan[i] <- struct{}{}
			}
		}(i)
	}

	// Gets messages from other servers and handles them by putting them in the appropriate channel
	for i := 0; i < variables.N; i++ {
		// Not myself
		if i == variables.Id {
			continue
		}
		// Initializes them with a goroutine and waits forever
		go func(i int) {
			for {
				message, err := RcvSockets[i].RecvBytes(0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}
				go handleMessage(message)
				_, err = RcvSockets[i].Send("OK.", 0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}
			}
		}(i)
	}
}

// Handles the requests from clients
func handleRequest(msg []byte) {
	cm := new(types.ClientMessage)
	buffer := bytes.NewBuffer(msg)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&cm)
	if err != nil {
		logger.ErrLogger.Println(len(msg))
		logger.ErrLogger.Fatal(err)
	}
	RequestChan <- cm
}

// Handles the messages from the other servers (i think only ReplicaStructure consern us)
func handleMessage(msg []byte) {
	count++
	logger.OutLogger.Println("Message Count", count)
	message := new(types.Message)
	buffer := bytes.NewBuffer([]byte(msg))
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&message)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	switch message.Type {
	case "CoordinationMessage":
		coordination := new(types.CoordinationMessage)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&coordination)
		if err != nil {
			logger.ErrLogger.Println(len(message.Payload))
			logger.ErrLogger.Fatal(err)
		}
		CoordChan <- struct {
			Message *types.CoordinationMessage
			From    int
		}{Message: coordination, From: message.From}
		break
	case "VCM":
		vcm := new(types.VCM)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&vcm)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		VcmChan <- struct {
			Vcm  types.VCM
			From int
		}{Vcm: *vcm, From: message.From}
		break
	case "Token":
		token := new(types.Token)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&token)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		TokenChan <- struct {
			Token types.Token
			From  int
		}{Token: *token, From: message.From}
		break
	case "ReplicaStructure":
		replica := new(types.ReplicaStructure)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&replica)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		ReplicaChan <- struct {
			Rep  *types.ReplicaStructure
			From int
		}{Rep: replica, From: message.From}
	}
}

// ReplyClient - Sends back a response to the client
func ReplyClient(reply *types.Reply) {
	logger.OutLogger.Println("Replying to Client.")
	to := reply.Client
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(reply)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	message := types.NewMessage(w.Bytes(), "Reply")
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(message)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	//_ = <-requestsChan[to]
	log.Printf("%s\n", string(reply.Result))

	_, err = ResponseSockets[to].SendBytes(w.Bytes(), 0)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	logger.OutLogger.Println("Replied to Client.")
	//SendRecvSync[to] <- struct{}{}
}
