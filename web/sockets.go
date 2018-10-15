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

		server.AddWebSocket(conn)
	}
}

func (web *WebServer) AddWebSocket(conn *websocket.Conn) {
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
			return
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
			return
		}

		log.Debugf("Writing %v", msg)
		if err := client.conn.WriteMessage(msg.msgType, msg.data); err != nil {
			log.Errorf("Error writing message: %v", err)

			// Is it close worthy though???? Yeah, probably
			client.conn.Close()
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
