package web

func (client *WebServerSocketClient) HandleHello(msg *ClientMessage) {
	log.Infof("Client %v said hello to '%v'", client.conn.RemoteAddr(), msg.Body["target"])
}
