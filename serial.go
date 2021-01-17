package stmprog

import (
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/tarm/serial"
)

type SerialFlasher struct {
	p    *serial.Port
	addr uint32
}

type SerialFlasherOpts struct {
	serial.Config
	Addr      uint32
	ReadChunk uint8
}

// DefaultOptions sets maximum baundrate for port (115200 bps). Use WithX to modify settings.
func DefaultOptions(port string) SerialFlasherOpts {
	return SerialFlasherOpts{
		Config: serial.Config{
			Name:        port,
			Baud:        115200,
			Parity:      serial.ParityEven,
			Size:        8,
			StopBits:    serial.Stop1,
			ReadTimeout: time.Millisecond * 500,
		},
		Addr:      APPLICATION_ADDRESS,
		ReadChunk: 128,
	}
}

// WithBaund returns new options with given baundrate.
func (o SerialFlasherOpts) WithBaund(baund int) SerialFlasherOpts {
	o.Config.Baud = baund
	return o
}

// WithApplicationAdress returns new options with given application address.
func (o SerialFlasherOpts) WithApplicationAdress(addr uint32) SerialFlasherOpts {
	o.Addr = addr
	return o
}

// WithReadChunk returns new options with selected read chunk size. Up to 256 bytes can be read in one command. When encountering issues with EEPROM reads, decrease chunk size.
func (o SerialFlasherOpts) WithReadChunk(chunk uint8) SerialFlasherOpts {
	o.ReadChunk = chunk
	return o
}

// NewSerialFlasher open selected port and flushes it.
func NewSerialFlasher(o SerialFlasherOpts) (*SerialFlasher, error) {
	conn, err := serial.OpenPort(&o.Config)
	if err != nil {
		return nil, fmt.Errorf("open port: %v", err)
	}
	err = conn.Flush()
	if err != nil {
		return nil, fmt.Errorf("flush port: %v", err)
	}
	flasher := SerialFlasher{p: conn, addr: o.Addr}
	return &flasher, err
}

func (p *SerialFlasher) Write(b []byte) (int, error) {
	writeCnt := uint32(0)
	for {
		start := int(writeCnt) * 256
		stop := start + 256
		if start > len(b) {
			return len(b), nil
		}
		writeAddr := p.addr + writeCnt*256
		if stop > len(b) {
			stop = len(b)
		}
		buf := b[start:stop]
		n := stop - start
		log.Infof("write %d bytes at: 0x%x", n, writeAddr)
		err := p.writeMemory(buf, writeAddr)
		if err != nil {
			return start, fmt.Errorf("writememory: %v", err)
		}
		writeCnt++
	}
}
func (p *SerialFlasher) Read(b []byte) (int, error) {
	readCnt := uint32(0)
	readChunk := 128
	for {
		start := int(readCnt) * readChunk
		stop := start + readChunk
		if start > len(b) {
			return len(b), nil
		}
		readAddr := p.addr + readCnt*uint32(readChunk)
		if stop > len(b) {
			stop = len(b)
		}
		buf := b[start:stop]
		n := stop - start
		log.Infof("read %d bytes at: 0x%x", n, readAddr)
		err := p.readMemory(buf, readAddr)
		if err != nil {
			return start, fmt.Errorf("writememory: %v", err)
		}
		readCnt++
	}
}

func (p *SerialFlasher) Erase() error {
	bootloader, err := p.Get()
	if err != nil {
		return fmt.Errorf("get: %v", err)
	}
	for _, command := range bootloader.Commands {
		if command == ERASE {
			return p.EraseGlobal()
		}
		if command == EXTENDED_ERASE {
			return p.ExtendedEraseGlobal()
		}
	}
	return errors.New("erase command not available")
}
