package main

import (
	"driver/internal/dualsense"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Разрешаем все origin (для тестов)
	},
}

type Direction struct {
	Up    bool `json:"up"`
	Down  bool `json:"down"`
	Left  bool `json:"left"`
	Right bool `json:"right"`
}

type ColorMessage struct {
	Color string `json:"color"`
}

type handler struct {
	dualsenseChannel chan dualsense.DualSenseState
	d                dualsense.DualSenseDevice
}

func NewHandler(dualsenseChannel chan dualsense.DualSenseState, d dualsense.DualSenseDevice) *handler {
	return &handler{
		dualsenseChannel: dualsenseChannel,
		d:                d,
	}
}

func (h *handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// Канал для сообщений от клиента
	clientMessages := make(chan ColorMessage)
	defer close(clientMessages)

	// Горутина для чтения сообщений от клиента
	go func() {
		for {
			var msg ColorMessage
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println("Read error:", err)
				return
			}
			clientMessages <- msg
		}
	}()

	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	log.Println("Client connected")
	for {
		select {
		case v := <-h.dualsenseChannel:
			dsState := v

			dir := Direction{}

			switch dsState.DPad {
			case "up":
				dir.Up = true
			case "down":
				dir.Down = true
			case "left":
				dir.Left = true
			case "right":
				dir.Right = true
			}

			err := conn.WriteJSON(dir)
			if err != nil {
				log.Println("Write direction error:", err)
				return
			}

		case msg := <-clientMessages:
			// Обработка сообщения с цветом от клиента
			log.Printf("Received color: %s", msg.Color)

			if r, g, b, err := ParseHexColor(msg.Color); err == nil {
				h.d.SetDualSenseColor(r, g, b)
			}

		case <-ticker.C:
			// Тикер для поддержания соединения
		}
	}
}

// ParseHexColor преобразует hex-строку (#RRGGBB) в RGB значения
func ParseHexColor(s string) (r, g, b uint8, err error) {
	if len(s) != 7 || s[0] != '#' {
		return 0, 0, 0, fmt.Errorf("invalid color format")
	}
	hexToByte := func(b byte) uint8 {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return 10 + b - 'a'
		case b >= 'A' && b <= 'F':
			return 10 + b - 'A'
		}
		return 0
	}
	r = hexToByte(s[1])<<4 + hexToByte(s[2])
	g = hexToByte(s[3])<<4 + hexToByte(s[4])
	b = hexToByte(s[5])<<4 + hexToByte(s[6])
	return r, g, b, nil
}

func main() {
	dsStateCh := make(chan dualsense.DualSenseState, 1)
	d := dualsense.DualSenseDevice{}
	d.ConnectDualsense()
	go d.Read(dsStateCh)

	h := NewHandler(dsStateCh, d)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		h.wsHandler(w, r)
	})
	http.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
