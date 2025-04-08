package main

import (
	"driver/internal/dualsense"
	"log"
	"net/http"
	"sync"
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

type handler struct {
	dualsenseChannel chan dualsense.DualSenseState
}

func NewHandler(dualsenseChannel chan dualsense.DualSenseState) *handler {
	return &handler{dualsenseChannel}
}

func (h *handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	log.Println("Client connected")

	ticker := time.NewTicker(16 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			dsState := <-h.dualsenseChannel

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
				log.Println("Write error:", err)
				return
			}
		}
	}
}

func main() {
	dsStateCh := make(chan dualsense.DualSenseState, 1)
	var mx sync.Mutex
	d := dualsense.DualSenseDevice{Mx: &mx}
	d.ConnectDualsense()
	// d.SetDualSenseColor(2, 225, 2)
	// go d.Read(dsStateCh)

	h := NewHandler(dsStateCh)

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		h.wsHandler(w, r)
	})
	http.Handle("/", http.FileServer(http.Dir(".")))

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
