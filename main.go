package main

import (
	"driver/internal/dualsense"
	"encoding/binary"
	"log"
	"time"
)

func main() {
	ds := &dualsense.DualSense{}

	if err := ds.Connect(); err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer ds.Close()

	log.Println("Enabling extended reports...")

	ds.ReadFeatureReport05()

	for {
		data, err := ds.ReadInput()
		if err != nil {
			log.Printf("Read error: %v", err)
			continue
		}

		if len(data) > 0 && data[0] == 0x31 {
			processSensorData(data)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func processSensorData(data []byte) {
	if len(data) < 50 {
		return
	}

	ax := int16(binary.LittleEndian.Uint16(data[16:18]))
	ay := int16(binary.LittleEndian.Uint16(data[18:20]))
	az := int16(binary.LittleEndian.Uint16(data[20:22]))

	log.Printf("Accel: X=%.2f Y=%.2f Z=%.2f",
		float32(ax)/16384.0,
		float32(ay)/16384.0,
		float32(az)/16384.0)
}
