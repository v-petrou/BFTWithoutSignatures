package messenger

import (
	"BFTWithoutSignatures/config"
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/variables"
	"SSBFT/types"
	"bytes"
	"encoding/gob"
	"strings"

	"github.com/pebbe/zmq4"
)

var (
	// Context to initialize sockets
	Context *zmq4.Context

	// ReceiveSockets - Receive messages from other servers
	ReceiveSockets map[int]*zmq4.Socket

	// SendSockets - Send messages to other servers
	SendSockets map[int]*zmq4.Socket

	// ServerSockets - Get the client requests
	ServerSockets map[int]*zmq4.Socket

	// ResponseSockets - Send responses to clients
	ResponseSockets map[int]*zmq4.Socket

	// SendRecvSync - Probably not needed
	SendRecvSync map[int]chan interface{}

	messageChannel = make(chan struct {
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

	RequestChannel = make(chan *types.ClientMessage, 100)

	ReplicaChan = make(chan struct {
		Rep  *types.ReplicaStructure
		From int
	}, 100)

	count = 0
)

// InitializeMessenger - Initializes the 0MQ sockets (servers communication with clients/servers)
func InitializeMessenger() {
	SendRecvSync = make(map[int]chan interface{}, variables.Clients) // Probably not needed, cause its only used here
	for i := 0; i < variables.Clients; i++ {
		SendRecvSync[i] = make(chan interface{})
	}

	Context, err := zmq4.NewContext()
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	// Sockets REP & PUB to communicate with each one of the clients
	ServerSockets = make(map[int]*zmq4.Socket, variables.Clients)
	ResponseSockets = make(map[int]*zmq4.Socket, variables.Clients)
	for i := 0; i < variables.Clients; i++ {

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
	ReceiveSockets = make(map[int]*zmq4.Socket)
	SendSockets = make(map[int]*zmq4.Socket)
	for i := 0; i < variables.N; i++ {
		// Not myself
		if i == variables.ID {
			continue
		}

		// ReceiveSockets initialization to get information from other servers
		ReceiveSockets[i], err = Context.NewSocket(zmq4.REP)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		var rcvAddr string
		if !variables.Remote {
			rcvAddr = strings.Replace(config.GetRepAddressLocal(i), "localhost", "*", 1)
		} else {
			rcvAddr = config.GetRepAddress(i)
		}
		err = ReceiveSockets[i].Bind(rcvAddr)
		if err != nil {
			logger.ErrLogger.Fatal(err, " "+rcvAddr)
		}
		logger.OutLogger.Println("Binded on ", rcvAddr)

		// SendSockets initialization to send information to other servers
		SendSockets[i], err = Context.NewSocket(zmq4.REQ)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		var sndAddr string
		if !variables.Remote {
			sndAddr = config.GetReqAddressLocal(i)
		} else {
			sndAddr = config.GetReqAddress(i)
		}
		err = SendSockets[i].Connect(sndAddr)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		logger.OutLogger.Println("Connected to ", sndAddr)
	}
}

// Broadcast - Broadcasts a message to other servers
func Broadcast(message types.Message) {
	for i := 0; i < variables.N; i++ {
		// Not myself
		if i == variables.ID {
			continue
		}
		SendMessage(message, i)
	}
}

// TransmitMessages - Echos the message to the other servers
func TransmitMessages() {
	for messageTo := range messageChannel {
		to := messageTo.To
		message := messageTo.Message
		w := new(bytes.Buffer)
		encoder := gob.NewEncoder(w)
		err := encoder.Encode(message)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		_, err = SendSockets[to].SendBytes(w.Bytes(), 0)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		//logger.OutLogger.Println("SENT Message to ", to)

		_, err = SendSockets[to].Recv(0)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		//logger.OutLogger.Println("OKAY FROM ", to)
	}
}

// SendMessage - Puts the messages in the message channel to be transmitted with TransmitMessages
func SendMessage(message types.Message, to int) {
	messageChannel <- struct {
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

// Subscribe - Handles the inputs from both clients and other servers
func Subscribe() {
	// Gets requests from clients and handles them
	for i := 0; i < variables.Clients; i++ {
		go func(i int) { // Initialize them with a goroutine and waits forever
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
			}
		}(i)
	}

	// Gets messages from other servers and handles them
	for i := 0; i < variables.N; i++ {
		// Not myself
		if i == variables.ID {
			continue
		}
		go func(i int) { // Initializes them with a goroutine and waits forever
			for {
				message, err := ReceiveSockets[i].RecvBytes(0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}

				go handleMessage(message)

				_, err = ReceiveSockets[i].Send("OK.", 0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}
			}
		}(i)
	}
}

// Put client's message in RequestChannel
func handleRequest(msg []byte) {
	message := new(types.ClientMessage)
	buffer := bytes.NewBuffer(msg)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&message)
	if err != nil {
		logger.ErrLogger.Println(len(msg))
		logger.ErrLogger.Fatal(err)
	}
	RequestChannel <- message
}

// Handles the messages from the other servers (i think only ReplicaStructure concern us)
func handleMessage(msg []byte) {
	count++
	logger.OutLogger.Println("Message Count:", count)

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
	to := reply.Client
	logger.OutLogger.Println("Replying to Client", to, "...")

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
	logger.OutLogger.Printf("%s\n", string(reply.Result))

	_, err = ResponseSockets[to].SendBytes(w.Bytes(), 0)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	logger.OutLogger.Println("Replied to Client", to)
}
