package web

type DmxData struct {
	Msg  string
	Body []int
}

func (client *WebServerSocketClient) HandleDmxMonStart(msg *ClientMessage) {
	client.wantsDmx = true
}

func (client *WebServerSocketClient) HandleDmxMonStop(msg *ClientMessage) {
	client.wantsDmx = false
}

func (client *WebServerSocketClient) SendDmx(frame []byte) {
	msg := &DmxData{
		Msg:  "dmxData",
		Body: make([]int, 512),
	}

	for i := 0; i < 512; i++ {
		msg.Body[i] = int(frame[i])
	}

	client.WriteJSON(msg)
}
