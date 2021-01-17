package stmprog

import (
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func ReadFile(path string) (int, []byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return 0, nil, err
	}
	b := make([]byte, int(fi.Size()))
	n, err := f.Read(b)
	return n, b, err
}

func crc(data []byte, magic byte) byte {
	for _, b := range data {
		magic ^= b
	}
	return magic
}

func (p *SerialFlasher) read(b []byte) error {
	n, err := p.p.Read(b)
	if n != len(b) {
		return fmt.Errorf("read %d bytes, wanted %d", n, len(b))
	}
	if err != nil {
		return err
	}
	log.Debugf("read  %d bytes: %x", n, b)
	return nil
}

func (p *SerialFlasher) readAck(b []byte) error {
	err := p.read(b)
	if err != nil {
		return err
	}
	if b[len(b)-1] != ACK {
		return fmt.Errorf("readAck: got %x", b[len(b)-1])
	}
	return nil
}

func (p *SerialFlasher) readAnyAck() ([]byte, error) {
	n := []byte{0}
	err := p.read(n)
	if err != nil {
		return nil, err
	}
	b := make([]byte, int(n[0]+2))
	err = p.readAck(b)
	return b, err
}

func (p *SerialFlasher) write(b []byte) error {
	n, err := p.p.Write(b)
	if n != len(b) {
		return fmt.Errorf("wrote %d bytes, wanted %d", n, len(b))
	}
	if err != nil {
		return err
	}
	log.Debugf("wrote %d bytes: %x", n, b)
	return nil
}

func (p *SerialFlasher) writeCmdAck(c Command) error {
	b := []byte{byte(c), byte(c) ^ 0xFF}
	return p.writeAck(b)
}

func (p *SerialFlasher) writeAck(b []byte) error {
	err := p.write(b)
	if err != nil {
		return err
	}
	err = p.getACK()
	if err != nil {
		return fmt.Errorf("writeAck: %v", err)
	}
	return nil
}

func (p *SerialFlasher) getACK() error {
	r := []byte{0}
	err := p.read(r)
	if err != nil {
		return fmt.Errorf("read: %v", err)
	}
	if r[0] == NACK {
		return errors.New("NACK received")
	}
	if r[0] != ACK {
		return fmt.Errorf("received unrecognized: 0x%x", r[0])
	}
	return nil
}

// func (p *SerialFlasher) sendSectors(sectors []byte) error {
// 	toSend := append([]byte{byte(len(sectors)) - 1}, sectors...)
// 	toSend = append(toSend, crc(toSend, 0x00))

// }
