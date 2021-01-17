package stmprog

import (
	"fmt"
	"io"

	log "github.com/sirupsen/logrus"
)

type Programmer struct {
	Flasher
	Booter
}

// Flasher initializes chip and reads, and writes to memory. It can erase memory too. Sector erase is not yet implemented.
type Flasher interface {
	io.ReadWriter
	Init() error
	Erase() error
}

// Booter switches STM to bootloader mode and resets CPU.
type Booter interface {
	Boot() error
	Reset() error
}

type GpioSerialProgrammer struct {
	*GpioBooter
	*SerialFlasher
}

func (p *Programmer) program(b []byte) error {
	err := p.Boot()
	if err != nil {
		return fmt.Errorf("boot: %v", err)
	}
	err = p.Init()
	if err != nil {
		return fmt.Errorf("init: %v", err)
	}
	err = p.Erase()
	if err != nil {
		return fmt.Errorf("erase: %v", err)
	}
	n, err := p.Write(b)
	log.Infof("program: wrote %x bytes", n)
	return err
}

func (p *Programmer) Program(b []byte) error {
	err := p.program(b)
	if err != nil {
		return err
	}
	return p.Reset()
}

func (p *Programmer) ProgramVerify(b []byte) error {
	err := p.program(b)
	if err != nil {
		return err
	}
	read := make([]byte, len(b))
	n, err := p.Read(read)
	if err != nil {
		return err
	}
	if n != len(b) {
		return fmt.Errorf("read: wanted %d, got %d bytes", len(b), n)
	}
	log.Infof("program: read %x bytes", n)
	for i := range b {
		if b[i] != read[i] {
			return fmt.Errorf("programmed: %x, read: %x at offset %x", b[i], read[i], i*8)
		}
	}
	log.Info("verify succeeded")
	return p.Reset()
}
