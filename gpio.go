package stmprog

import (
	"time"

	"periph.io/x/periph/conn/gpio"
)

type GpioBooter struct {
	BootPin gpio.PinOut
	RstPin  gpio.PinOut
}

func (dp *GpioBooter) Reset() error {
	dp.RstPin.Out(gpio.High)
	time.Sleep(200 * time.Millisecond)
	err := dp.RstPin.Out(gpio.Low)
	time.Sleep(200 * time.Millisecond)
	return err
}

func (dp *GpioBooter) Boot() error {
	dp.BootPin.Out(gpio.High)
	time.Sleep(200 * time.Millisecond)
	dp.Reset()
	return dp.BootPin.Out(gpio.Low)
}
