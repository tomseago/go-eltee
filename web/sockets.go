package web

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"strings"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewWebSocketHandler(server *WebServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Errorf("Unable to upgrade socket: %v", err)
			return
		}

		server.NewClientFor(conn)
	}
}

func (web *WebServer) NewClientFor(conn *websocket.Conn) {
	client := &WebServerSocketClient{
		web:  web,
		conn: conn,

		writeQ: make(chan *WSMessage, 10),
	}

	go client.readPump()

	go client.writePump()

	web.sockets = append(web.sockets, client)
	log.Infof("New websocket connection from %v", conn.RemoteAddr())
}

func (web *WebServer) RemoveClient(client *WebServerSocketClient) {
	ix := -1
	for i := 0; i < len(web.sockets); i++ {
		if web.sockets[i] == client {
			ix = i
			break
		}
	}
	if ix != -1 {
		web.sockets[ix] = web.sockets[len(web.sockets)-1]
		web.sockets = web.sockets[:len(web.sockets)-1]
	}
	// log.Infof("New websocket connection from %v", conn.RemoteAddr())
}

////////////////////////////////////////////////////////

type WSMessage struct {
	msgType int
	data    []byte
}

func (msg *WSMessage) String() string {
	var b strings.Builder
	if msg.msgType == websocket.BinaryMessage {
		b.WriteString(fmt.Sprintf("BIN:(%v) ", len(msg.data)))

		for ix := 0; ix < len(msg.data); ix++ {
			b.WriteString(fmt.Sprintf("%02x", msg.data[ix]))
		}
	} else {
		b.WriteString("TXT: ")
		b.WriteString(string(msg.data))
	}

	return b.String()
}

////////////////////////////////////////////////////////

type WebServerSocketClient struct {
	web  *WebServer
	conn *websocket.Conn

	writeQ chan *WSMessage

	// Local state of interest
	wantsDmx bool
}

func (client *WebServerSocketClient) readPump() {

	log.Debug("Started read pump")

	msg := &WSMessage{}
	var err error

	for {
		log.Info("Calling ReadMessage...")
		msg.msgType, msg.data, err = client.conn.ReadMessage()
		if err != nil {
			log.Infof("Closing connection from %v: %v", client.conn.RemoteAddr(), err)
			break
		}

		err = client.handleMessage(msg)
		if err != nil {
			log.Warningf("Error trying to handle message %v from %v : %v", msg, client.conn.RemoteAddr(), err)
			// keep going though. Not necessarily a close right???
		}
	}

	// Make sure the write shuts down
	close(client.writeQ)
}

type ClientMessage struct {
	Msg  string
	Body map[string]interface{}
}

// func (msg *ClientMessage) bvAsFloat(key string) float64 {
// 	if msg == nil {
// 		return 0.0
// 	}

// 	val := msg.Body[key]
// 	return ValAsFloat(val)
// }

// func (msg *ClientMessage) bvAsInt(key string) int {
// 	if msg == nil {
// 		return 0
// 	}

// 	val := msg.Body[key]
// 	return ValAsInt(val)
// }

// func (msg *ClientMessage) bvAsString(key string) string {
// 	if msg == nil {
// 		return ""
// 	}

// 	val := msg.Body[key]
// 	return ValAsString(val)
// }

func (client *WebServerSocketClient) handleMessage(msg *WSMessage) error {
	log.Infof("From %v got %v", client.conn.RemoteAddr(), msg)

	if msg.msgType == websocket.TextMessage {
		if len(msg.data) == 0 {
			return fmt.Errorf("Message has no length")
		}

		// Interpret it as a JSON message. This is probably most all messages yeah?
		var cMsg ClientMessage
		if err := json.Unmarshal(msg.data, &cMsg); err != nil {
			// TODO: Maybe understand it as raw text???
			return err
		}

		log.Debugf("Decoded message %v", cMsg)

		switch cMsg.Msg {
		case "hello":
			client.HandleHello(&cMsg)

		case "reqFixtures":
			client.HandleReqFixtures(&cMsg)

		case "reqProfiles":
			client.HandleReqProfiles(&cMsg)

		case "reqCurrentState":
			client.HandleReqCurrentState(&cMsg)

		case "setCP":
			client.HandleSetCP(&cMsg)

		case "dmxMonStart":
			client.HandleDmxMonStart(&cMsg)

		case "dmxMonStop":
			client.HandleDmxMonStop(&cMsg)

		default:
			log.Warningf("Don't know how to handle this: %v", cMsg)
		}
	} else {
		return fmt.Errorf("Don't currently handle binary messages!!!!")
	}

	return nil
}

////
func (client *WebServerSocketClient) writePump() {
	log.Debug("Started write pump")

	for {
		log.Warningf("Waiting to read from channel client=%v", client)
		msg := <-client.writeQ
		log.Warningf("Read from channel")

		if msg == nil {
			log.Warningf("Closing the writePump()")
			client.web.RemoveClient(client)
			return
		}

		log.Debugf("Writing %v", msg)
		if err := client.conn.WriteMessage(msg.msgType, msg.data); err != nil {
			log.Errorf("Error writing message: %v", err)

			// Is it close worthy though???? Yeah, probably
			client.conn.Close()
			client.web.RemoveClient(client)
			return
		}
	}
}

func (client *WebServerSocketClient) WriteText(data []byte) {
	log.Errorf("WriteText is not implemented")

}

func (client *WebServerSocketClient) WriteJSON(data interface{}) {
	log.Debugf("WriteJSON(%v)", data)

	msg := &WSMessage{
		msgType: websocket.TextMessage,
	}

	text, err := json.Marshal(data)
	if err != nil {
		log.Errorf("Unable to marshall message: %v", data)
		return
	}

	msg.data = []byte(text)

	log.Warningf("Writing to channel: %v  client=%v", msg, client)
	client.writeQ <- msg
}

func (client *WebServerSocketClient) WriteBinary(data []byte) {
	log.Errorf("WriteBinary is not implemented")
}
