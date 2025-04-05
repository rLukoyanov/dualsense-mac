package main

import (
	"fmt"
	"log"
	"time"

	"github.com/sstallion/go-hid"
)

const (
	SonyVendorID       = 0x054c
	DualSenseUSB       = 0x0ce6
	DualSenseBluetooth = 0x0df2
)

type DualSense struct {
	LeftStick  [2]uint8
	RightStick [2]uint8
}

func main() {
	if err := hid.Init(); err != nil {
		log.Fatalf("HID init failed: %v", err)
	}

	device, err := connectDualSense()
	if err != nil {
		log.Fatal(err)
	}
	defer device.Close()

	fmt.Println("DualSense подключен успешно!")
	fmt.Println("Нажмите кнопки или используйте стики...")

	buf := make([]byte, 256)
	for {
		n, err := device.Read(buf)
		if err != nil {
			log.Printf("Read error: %v", err)
			time.Sleep(time.Second)
			continue
		}

		log.Println(buf[:n])
		// connectDataToStruct(buf[:n])

		time.Sleep(16 * time.Millisecond)
	}
}

func connectDualSense() (*hid.Device, error) {
	return hid.OpenFirst(SonyVendorID, DualSenseUSB)
}

func connectDataToStruct(data []byte) {
	// законектить байты к кнопкам
	// 2025/04/06 00:30:05 [1 левый стик - 129 127 правый стик - 126 128 0 0 95 кнопки слева и справа - 8 0 0 0 66 73 114 198 2 0 0 0 4 0 235 255 121 31 137 5 125 58 78 80 22 141 232 64 65 128 0 0 0 228 9 9 0 0 0 0 0 59 77 78 80 18 24 0 67 154 3 69 27 124 116 190]
}
