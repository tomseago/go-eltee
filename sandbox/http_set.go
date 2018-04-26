package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func SendDmxX(universe int, data []byte) {
	vals := make(url.Values)

	vals.Set("u", strconv.Itoa(universe))
	vals.Set("d", fmt.Sprintf("%v", data))

	fmt.Printf("vals=%v", vals)

	resp, err := http.PostForm("http://localhost:9090/set_dmx", vals)
	fmt.Printf("resp=%v\nerr=%v", resp, err)
}

func SendDmx(universe int, data []byte) {
	var b strings.Builder

	b.WriteString("u=")
	b.WriteString(strconv.Itoa(universe))
	b.WriteString("&d=")
	for i := 0; i < len(data); i++ {
		b.WriteString(strconv.Itoa(int(data[i])))

		if i < len(data)-1 {
			b.WriteString(",")
		}
	}

	dataStr := b.String()
	fmt.Printf("data = %v\n", dataStr)
	reader := strings.NewReader(dataStr)
	resp, err := http.Post("http://localhost:9090/set_dmx", "application/x-www-form-urlencoded", reader)
	fmt.Printf("resp=%v\nerr=%v", resp, err)
}

func main() {
	universe := 1
	data := make([]byte, 0)

	for i := 0; i < 10; i++ {
		data = append(data, byte(i+1))
	}

	SendDmx(universe, data)
}
