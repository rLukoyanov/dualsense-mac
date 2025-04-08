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
		49, 112, 16, 255, 247, 0, 0, 0, 0, 0, 0,
		0, 16, 38, 144, 160, 255, 0, 0, 0, 0,
		0, 0, 0, 38, 144, 160, 255, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 2, 0, 2, 0, 0, byte(r), byte(g), byte(b), 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 79, 52, 31, 57,
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
	d.Mx.Lock()
	if err := hid.Init(); err != nil {
		log.Fatalf("хуй: %v", err)
	}

	device, err := hid.OpenFirst(SonyVendorID, DualSenseUSB)
	if err != nil {
		log.Fatal("Connect error:", err)
	}

	reportID := byte(0x05)
	report := make([]byte, 64)
	report[0] = reportID
	n, err := device.GetFeatureReport(report)
	if err != nil {
		log.Fatal("Connect error:", err)
	}

	d.device = device
	log.Printf("device connectted:%v bytes: %v\n", d.device, n)
	d.Mx.Unlock()
}

func (d *DualSenseDevice) Read(ch chan DualSenseState) {
	log.Println("device reading", d.device)

	buf := make([]byte, 64)
	for {
		n, err := d.device.Read(buf)
		log.Println(n, buf)
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		if n >= 10 && buf[0] == 0x01 {
			ch <- processInput(buf[:n])
		}

		time.Sleep(16 * time.Millisecond)
	}
}

func processInput(data []byte) DualSenseState {
	// 01 7d 7e 83 82 08 00 00 00 00 хуйня с блютуза
	report := DualSenseReport{
		ReportID:     data[0],      // 0x01
		LeftX:        data[1],      // 0x7d
		LeftY:        data[2],      // 0x7e
		RightX:       data[3],      // 0x83
		RightY:       data[4],      // 0x82
		DPadButtons:  data[5],      // 0x08
		Buttons1:     data[5] >> 4, // Верхние 4 бита
		Buttons2:     data[6],      // 0x00
		Buttons3:     data[7],      // 0x00
		LeftTrigger:  data[8],      // 0x00
		RightTrigger: data[9],      // 0x00
	}

	state := parseReport(report)
	printState(state)
	return state
}

func parseReport(report DualSenseReport) DualSenseState {
	var state DualSenseState

	state.LeftStick.X = int((float64(report.LeftX) - StickNeutral) / StickNeutral * 100)
	state.LeftStick.Y = int((float64(report.LeftY) - StickNeutral) / StickNeutral * 100)
	state.RightStick.X = int((float64(report.RightX) - StickNeutral) / StickNeutral * 100)
	state.RightStick.Y = int((float64(report.RightY) - StickNeutral) / StickNeutral * 100)

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

	state.Buttons.Square = report.Buttons1&0x01 > 0
	state.Buttons.Cross = report.Buttons1&0x02 > 0
	state.Buttons.Circle = report.Buttons1&0x04 > 0
	state.Buttons.Triangle = report.Buttons1&0x08 > 0

	state.Buttons.Options = report.Buttons2&0x10 > 0
	state.Buttons.R1 = report.Buttons2&0x20 > 0
	state.Buttons.L3 = report.Buttons2&0x40 > 0
	state.Buttons.R3 = report.Buttons2&0x80 > 0
	state.Buttons.L1 = report.Buttons2&0x01 > 0
	state.Buttons.R1 = report.Buttons2&0x02 > 0

	state.Buttons.PS = report.Buttons3&0x01 > 0
	state.Buttons.Touchpad = report.Buttons3&0x02 > 0

	state.Triggers.Left = uint8(float64(report.LeftTrigger) / 2.55)
	state.Triggers.Right = uint8(float64(report.RightTrigger) / 2.55)

	return state
}

func printState(state DualSenseState) {
	fmt.Printf("\033[H\033[2J")
	fmt.Println("=== DualSense State ===")

	fmt.Printf("\nстики:\n")
	fmt.Printf("Left:  X=%-4d%% Y=%-4d%%\n", state.LeftStick.X, state.LeftStick.Y)
	fmt.Printf("Right: X=%-4d%% Y=%-4d%%\n", state.RightStick.X, state.RightStick.Y)

	fmt.Printf("\nКнопки:\n")
	fmt.Printf("Фигуры: △=%v ○=%v ×=%v □=%v\n",
		state.Buttons.Triangle, state.Buttons.Circle,
		state.Buttons.Cross, state.Buttons.Square)

	fmt.Printf("\nD-pad: %s\n", state.DPad)
	fmt.Printf("L1=%v R1=%v L3=%v R3=%v\n",
		state.Buttons.L1, state.Buttons.R1,
		state.Buttons.L3, state.Buttons.R3)
	fmt.Printf("Triggers: L2=%-3d%% R2=%-3d%%\n",
		state.Triggers.Left, state.Triggers.Right)

	fmt.Printf("System: PS=%v Touchpad=%v\n",
		state.Buttons.PS, state.Buttons.Touchpad)
}
