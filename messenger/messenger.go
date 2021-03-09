package messenger

import (
	"BFTWithoutSignatures/config"
	"BFTWithoutSignatures/logger"
	"BFTWithoutSignatures/threshenc"
	"BFTWithoutSignatures/types"
	"BFTWithoutSignatures/variables"
	"bytes"
	"encoding/gob"

	"github.com/pebbe/zmq4"
)

// Sockets
var (
	// Context to initialize sockets
	Context *zmq4.Context

	// SendSockets - Send messages to other servers
	SendSockets map[int]*zmq4.Socket

	// ReceiveSockets - Receive messages from other servers
	ReceiveSockets map[int]*zmq4.Socket

	// ServerSockets - Get the client requests
	ServerSockets map[int]*zmq4.Socket

	// ResponseSockets - Send responses to clients
	ResponseSockets map[int]*zmq4.Socket
)

// Channels for messages
var (
	// MessageChannel - Channel to put the messages that need to be transmitted in
	MessageChannel = make(chan struct {
		Message types.Message
		To      int
	})

	// BvbChannel - Channel to put the BVB messages in
	BvbChannel = make(map[int]chan struct {
		BcMessage types.BcMessage
		From      int
	})

	// BcChannel - Channel to put the BC messages in
	BcChannel = make(map[int]chan struct {
		BcMessage types.BcMessage
		From      int
	})

	// RbChannel - Channel to put the RB messages in
	RbChannel = make(map[string]map[int]chan struct {
		RbMessage types.RbMessage
		From      int
	})

	// RbAbcChannel - Channel to put the RB messages for ABC in
	RbAbcChannel = make(chan struct {
		RbMessage types.RbMessage
		From      int
	})

	// MvcChannel - Channel to put the MVC messages in
	MvcChannel = make(map[int]chan struct {
		MvcMessage types.MvcMessage
		From       int
	})

	// VcChannel - Channel to put the VC messages in
	VcChannel = make(map[int]chan struct {
		VcMessage types.VcMessage
		From      int
	})

	// AbcChannel - Channel to put the ABC messages in
	AbcChannel = make(chan struct {
		AbcMessage types.AbcMessage
		From       int
	})

	// RequestChannel - Channel to put the client requests in
	RequestChannel = make(chan types.ClientMessage, 100)
)

// InitializeMessenger - Initializes the 0MQ sockets (between Servers and Clients)
func InitializeMessenger() {
	Context, err := zmq4.NewContext()
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	// Initialization of a socket pair to communicate with each one of the other servers
	ReceiveSockets = make(map[int]*zmq4.Socket)
	SendSockets = make(map[int]*zmq4.Socket)
	for i := 0; i < variables.N; i++ {
		if i == variables.ID {
			continue // Not myself
		}

		// ReceiveSockets initialization to get information from other servers
		ReceiveSockets[i], err = Context.NewSocket(zmq4.REP)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		var receiveAddr string
		if !variables.Remote {
			receiveAddr = config.GetRepAddressLocal(i)
		} else {
			receiveAddr = config.GetRepAddress(i)
		}
		err = ReceiveSockets[i].Bind(receiveAddr)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		logger.OutLogger.Println("Receive from Server", i, "on", receiveAddr)

		// SendSockets initialization to send information to other servers
		SendSockets[i], err = Context.NewSocket(zmq4.REQ)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		var sendAddr string
		if !variables.Remote {
			sendAddr = config.GetReqAddressLocal(i)
		} else {
			sendAddr = config.GetReqAddress(i)
		}
		err = SendSockets[i].Connect(sendAddr)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		logger.OutLogger.Println("Send to Server", i, "on", sendAddr)
	}

	logger.OutLogger.Println("-----------------------------------------")

	// Initialization of a socket pair to communicate with each one of the clients
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
		logger.OutLogger.Println("Requests from Client", i, "on", serverAddr)

		// ResponseSockets initialization to publish the response back to the clients
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
		logger.OutLogger.Println("Response to Client", i, "on", responseAddr)
	}

	initRBChannels()

	logger.OutLogger.Print("-----------------------------------------\n\n")
}

// Initializes the RB channels
func initRBChannels() {
	RbChannel["MVC"] = make(map[int]chan struct {
		RbMessage types.RbMessage
		From      int
	})

	RbChannel["VC"] = make(map[int]chan struct {
		RbMessage types.RbMessage
		From      int
	})
}

// Broadcast - Broadcasts a message to all other servers
func Broadcast(message types.Message) {
	for i := 0; i < variables.N; i++ {
		if i == variables.ID {
			continue // Not myself
		}
		SendMessage(message, i)
	}
}

// SendMessage - Puts the messages in the message channel to be transmitted
func SendMessage(message types.Message, to int) {
	MessageChannel <- struct {
		Message types.Message
		To      int
	}{Message: message, To: to}
}

// TransmitMessages - Transmits the messages to the other servers [go started from main]
func TransmitMessages() {
	for messageTo := range MessageChannel {
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

		_, err = SendSockets[to].Recv(0)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}
		logger.OutLogger.Println("SENT", messageTo.Message.Type, "to", to)
	}
}

// Subscribe - Handles the inputs from both clients and other servers
func Subscribe() {
	// Gets messages from other servers and handles them
	for i := 0; i < variables.N; i++ {
		if i == variables.ID {
			continue // Not myself
		}
		go func(i int) { // Initializes them with a goroutine and waits forever
			for {
				message, err := ReceiveSockets[i].RecvBytes(0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}

				go HandleMessage(message)

				_, err = ReceiveSockets[i].Send("OK.", 0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}
			}
		}(i)
	}

	// Gets requests from clients and handles them
	for i := 0; i < variables.Clients; i++ {
		go func(i int) { // Initialize them with a goroutine and waits forever
			for {
				message, err := ServerSockets[i].RecvBytes(0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}

				go handleRequest(message, i)

				_, err = ServerSockets[i].Send("", 0)
				if err != nil {
					logger.ErrLogger.Fatal(err)
				}
			}
		}(i)
	}
}

// Put client's message in RequestChannel to be handled
func handleRequest(msg []byte, from int) {
	var message rune
	buffer := bytes.NewBuffer(msg)
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&message)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	logger.OutLogger.Println("RECEIVED REQ from", from)
	RequestChannel <- types.NewClientMessage(from, message)
}

// HandleMessage - Handles the messages from the other servers
func HandleMessage(msg []byte) {
	message := new(types.Message)
	buffer := bytes.NewBuffer([]byte(msg))
	decoder := gob.NewDecoder(buffer)
	err := decoder.Decode(&message)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	if !(threshenc.VerifyMessage(message.Payload, message.Signature, message.From)) {
		logger.OutLogger.Println("INVALID", message.Type, "from", message.From)
		return
	}

	logger.OutLogger.Println("RECEIVED", message.Type, "from", message.From)

	switch message.Type {
	case "BVB":
		bcMessage := new(types.BcMessage)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&bcMessage)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		tag := bcMessage.Tag
		if _, in := BvbChannel[tag]; !in {
			BvbChannel[tag] = make(chan struct {
				BcMessage types.BcMessage
				From      int
			})
		}

		BvbChannel[tag] <- struct {
			BcMessage types.BcMessage
			From      int
		}{BcMessage: *bcMessage, From: message.From}

	case "BC":
		bcMessage := new(types.BcMessage)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&bcMessage)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		tag := bcMessage.Tag
		if _, in := BcChannel[tag]; !in {
			BcChannel[tag] = make(chan struct {
				BcMessage types.BcMessage
				From      int
			})
		}

		BcChannel[tag] <- struct {
			BcMessage types.BcMessage
			From      int
		}{BcMessage: *bcMessage, From: message.From}

	case "RB":
		rbMessage := new(types.RbMessage)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&rbMessage)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		rbid := rbMessage.Rbid
		rbType := rbMessage.Type
		if _, in := RbChannel[rbType][rbid]; !in {
			RbChannel[rbType][rbid] = make(chan struct {
				RbMessage types.RbMessage
				From      int
			})
		}

		RbChannel[rbType][rbid] <- struct {
			RbMessage types.RbMessage
			From      int
		}{RbMessage: *rbMessage, From: message.From}

	case "RB_ABC":
		rbMessage := new(types.RbMessage)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&rbMessage)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		RbAbcChannel <- struct {
			RbMessage types.RbMessage
			From      int
		}{RbMessage: *rbMessage, From: message.From}

	case "MVC":
		mvcMessage := new(types.MvcMessage)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&mvcMessage)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		cid := mvcMessage.Cid
		if _, in := MvcChannel[cid]; !in {
			MvcChannel[cid] = make(chan struct {
				MvcMessage types.MvcMessage
				From       int
			})
		}

		MvcChannel[cid] <- struct {
			MvcMessage types.MvcMessage
			From       int
		}{MvcMessage: *mvcMessage, From: message.From}

	case "VC":
		vcMessage := new(types.VcMessage)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&vcMessage)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		vcid := vcMessage.Vcid
		if _, in := VcChannel[vcid]; !in {
			VcChannel[vcid] = make(chan struct {
				VcMessage types.VcMessage
				From      int
			})
		}

		VcChannel[vcid] <- struct {
			VcMessage types.VcMessage
			From      int
		}{VcMessage: *vcMessage, From: message.From}

	case "ABC":
		abcMessage := new(types.AbcMessage)
		buf := bytes.NewBuffer(message.Payload)
		dec := gob.NewDecoder(buf)
		err = dec.Decode(&abcMessage)
		if err != nil {
			logger.ErrLogger.Fatal(err)
		}

		AbcChannel <- struct {
			AbcMessage types.AbcMessage
			From       int
		}{AbcMessage: *abcMessage, From: message.From}
	}
}

// ReplyClient - Sends back a response to the client
func ReplyClient(reply []byte, to int) {
	w := new(bytes.Buffer)
	encoder := gob.NewEncoder(w)
	err := encoder.Encode(reply)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}

	_, err = ResponseSockets[to].SendBytes(w.Bytes(), 0)
	if err != nil {
		logger.ErrLogger.Fatal(err)
	}
	logger.OutLogger.Println("Replied to Client", to, "(", string(reply), ")")
}
