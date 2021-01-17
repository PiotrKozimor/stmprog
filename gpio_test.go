package stmprog

import (
	"testing"

	"periph.io/x/periph/host"
	"periph.io/x/periph/host/rpi"
)

var TestBooter = GpioBooter{
	BootPin: rpi.SO_65,
	RstPin:  rpi.SO_63,
}

func TestReset(t *testing.T) {
	_, err := host.Init()
	handle(err, t)
	dut := TestBooter
	handle(dut.Reset(), t)
}

func TestBoot(t *testing.T) {
	_, err := host.Init()
	handle(err, t)
	dut := TestBooter
	handle(dut.Boot(), t)
}
