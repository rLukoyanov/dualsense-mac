package dualsense

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/sstallion/go-hid"
)

const (
	SonyVendorID = 0x054c
	DualSenseUSB = 0x0ce6
	StickNeutral = 0x80
)

type DualSenseDevice struct {
	device *hid.Device
	Mx     *sync.Mutex
}

type DualSenseReport struct {
	ReportID     uint8
	LeftX        uint8
	LeftY        uint8
	RightX       uint8
	RightY       uint8
	DPadButtons  uint8
	Buttons1     uint8
	Buttons2     uint8
	Buttons3     uint8
	Reserved     uint8
	LeftTrigger  uint8
	RightTrigger uint8
}

type DualSenseState struct {
	LeftStick struct {
		X, Y int
	}
	RightStick struct {
		X, Y int
	}
	Buttons struct {
		Square, Cross, Circle, Triangle bool
		L1, R1, L2, R2                  bool
		Create, Options, L3, R3         bool
		PS, Touchpad                    bool
	}
	DPad     string
	Triggers struct {
		Left, Right uint8
	}
}

func (d *DualSenseDevice) SetDualSenseColor(r, g, b byte) error {
	data := []byte{
		0x02, 0xff, 0xf7, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, r, g, b,
	}

	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, &data); err != nil {
		return fmt.Errorf("binary write failed: %v", err)
	}

	n, err := d.device.Write(buf.Bytes())
	if err != nil {
		return fmt.Errorf("device write failed: %v", err)
	}
	if n != len(buf.Bytes()) {
		return fmt.Errorf("incomplete write: %d/%d bytes", n, len(buf.Bytes()))
	}

	log.Printf("Отправлено %d байт: [% X]", n, buf.Bytes())
	return nil
}

func (d *DualSenseDevice) ConnectDualsense() {
	if err := hid.Init(); err != nil {
		log.Fatalf("хуй: %v", err)
	}

	device, err := hid.OpenFirst(SonyVendorID, DualSenseUSB)
	if err != nil {
		log.Fatal("Connect error:", err)
	}

	if err != nil {
		log.Fatal("Connect error:", err)
	}

	d.device = device
}

func (d *DualSenseDevice) Read(ch chan DualSenseState) {
	log.Println("device reading", d.device)

	buf := make([]byte, 78)
	for {
		n, err := d.device.Read(buf)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		ch <- processInput(buf[:n])

		time.Sleep(16 * time.Millisecond)
	}
}

func processInput(data []byte) DualSenseState {
	report := DualSenseReport{
		ReportID:    data[0],
		DPadButtons: data[8],
	}

	state := parseReport(report)
	return state
}

func parseReport(report DualSenseReport) DualSenseState {
	var state DualSenseState

	switch report.DPadButtons & 0x0F {
	case 0x0:
		state.DPad = "up"
	case 0x1:
		state.DPad = "up-rigth"
	case 0x2:
		state.DPad = "right"
	case 0x3:
		state.DPad = "down-rigth"
	case 0x4:
		state.DPad = "down"
	case 0x5:
		state.DPad = "down-left"
	case 0x6:
		state.DPad = "left"
	case 0x7:
		state.DPad = "up-left"
	case 0x8:
		state.DPad = "Neutral"
	}

	return state
}
