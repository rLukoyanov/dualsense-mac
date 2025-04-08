package dualsense

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/sstallion/go-hid"
)

const (
	SonyVendorID        = 0x054c
	DualSenseConn       = 0x0ce6
	InputReportSizeUSB  = 64
	InputReportSizeBT   = 78
	OutputReportSizeUSB = 47
	OutputReportSizeBT  = 78
)

type ConnectionType int

const (
	ConnectionUnknown ConnectionType = iota
	ConnectionUSB
	ConnectionBluetooth
)

type DualSense struct {
	device         *hid.Device
	connectionType ConnectionType
}

func (d *DualSense) Connect() error {
	if err := hid.Init(); err != nil {
		return err
	}

	device, err := hid.OpenFirst(SonyVendorID, DualSenseConn)
	if err == nil {
		d.device = device
		report := make([]byte, 64)
		n, err := device.Read(report)
		if err != nil {
			return err
		}
		if n == InputReportSizeUSB {
			d.connectionType = ConnectionUSB
		} else {
			d.connectionType = ConnectionBluetooth
		}
		log.Println("Connected via ", d.connectionType)
		return nil
	}

	return err
}

func (d *DualSense) Close() error {
	if d.device != nil {
		return d.device.Close()
	}
	return nil
}

func (d *DualSense) ReadInput() ([]byte, error) {
	var buf []byte

	switch d.connectionType {
	case ConnectionUSB:
		buf = make([]byte, InputReportSizeUSB)
	case ConnectionBluetooth:
		buf = make([]byte, InputReportSizeBT)
		buf[0] = 0x05
		_, err := d.device.GetFeatureReport(buf)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown connection type")
	}

	n, err := d.device.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}

func (d *DualSense) SendOutput(data []byte) error {
	var expectedSize int

	switch d.connectionType {
	case ConnectionUSB:
		expectedSize = OutputReportSizeUSB
	case ConnectionBluetooth:
		expectedSize = OutputReportSizeBT
	default:
		return fmt.Errorf("unknown connection type")
	}

	if len(data) != expectedSize {
		return fmt.Errorf("invalid output report size, expected %d, got %d", expectedSize, len(data))
	}

	_, err := d.device.Write(data)
	return err
}

func (d *DualSense) GetConnectionType() ConnectionType {
	return d.connectionType
}

func (d *DualSense) ReadFeatureReport05() error {
	if d.connectionType != ConnectionBluetooth {
		return nil
	}

	if d.device == nil {
		return errors.New("device not connected")
	}

	report := make([]byte, 1)
	report[0] = 0x05

	start := time.Now()
	n, err := d.device.SendFeatureReport(report)
	if err != nil {
		return fmt.Errorf("failed to send feature report 0x05: %v", err)
	}

	if n != len(report) {
		return fmt.Errorf("incomplete feature report sent: %d/%d bytes", n, len(report))
	}

	log.Printf("Feature report 0x05 sent successfully (took %v)", time.Since(start))
	return nil
}
