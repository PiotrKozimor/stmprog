package main

import (
	"log"

	"github.com/PiotrKozimor/stmprog"
	"github.com/sirupsen/logrus"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/rpi"
)

var booter = stmprog.GpioBooter{
	BootPin: rpi.SO_65,
	RstPin:  rpi.SO_63,
}

func handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	logrus.SetLevel(logrus.InfoLevel)
	_, err := host.Init()
	handle(err)
	flasher, err := stmprog.NewSerialFlasher(stmprog.DefaultOptions("/dev/ttyUSB2").WithApplicationAdress(0x8000000).WithReadChunk(128))
	handle(err)
	prog := &stmprog.Programmer{Booter: &booter, Flasher: flasher}
	_, b, err := stmprog.ReadFile("app_firmware.bin")
	handle(err)
	err = prog.ProgramVerify(b)
	handle(err)
}
