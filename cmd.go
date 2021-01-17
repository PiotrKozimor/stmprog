package stmprog

import (
	"encoding/binary"
	"fmt"
)

const (
	ACK                 byte    = 0x79
	NACK                byte    = 0x1F
	INIT                        = 0x7F
	GET                 Command = 0x00
	GET_VERSION_AND_RPS Command = 0x01
	READ_MEMORY         Command = 0x11
	WRITE_MEMORY        Command = 0x31
	ERASE               Command = 0x43
	EXTENDED_ERASE      Command = 0x44
	WRITE_PROTECT       Command = 0x63
	WRITE_UNPROTECT     Command = 0x73
	READOUT_PROTECT     Command = 0x82
	READOUT_UNPROTECT   Command = 0x92
)
const APPLICATION_ADDRESS = 0x8000000

type Sector byte
type Command byte

//Rdp - read protection
type Rdp byte

//RdpStatus is not correctly reported, only Version is valid
type RdpStatus struct {
	Rdp          byte
	RdpToggleCnt byte
	Version      byte
}
type Bootloader struct {
	Version  byte
	Commands []Command
}

func (p *SerialFlasher) Init() error {
	return p.writeAck([]byte{INIT, INIT ^ 0xFF})
}

func (p *SerialFlasher) Get() (Bootloader, error) {
	err := p.writeCmdAck(GET)
	if err != nil {
		return Bootloader{}, err
	}
	b, err := p.readAnyAck()
	if err != nil {
		return Bootloader{}, err
	}
	commands := make([]Command, len(b)-1)
	for i, b := range b[1:] {
		commands[i] = Command(b)
	}
	return Bootloader{
		Version:  b[0],
		Commands: commands,
	}, err
}
func (p *SerialFlasher) GetVersionAndReadProtectionStatus() (RdpStatus, error) {
	err := p.writeCmdAck(GET_VERSION_AND_RPS)
	if err != nil {
		return RdpStatus{}, err
	}
	b := make([]byte, 4)
	p.readAck(b)
	return RdpStatus{
		Version:      b[0],
		Rdp:          b[1],
		RdpToggleCnt: b[2],
	}, err
}

func (p *SerialFlasher) readMemory(buf []byte, addr uint32) error {
	if len(buf) > 256 {
		return fmt.Errorf("maximum size of buffer is 256, got: %d", len(buf))
	}
	err := p.writeCmdAck(READ_MEMORY)
	if err != nil {
		return fmt.Errorf("readMemory: %v", err)
	}
	readAddr := make([]byte, 4)
	binary.BigEndian.PutUint32(readAddr, addr)
	readAddr = append(readAddr, crc(readAddr, 0x00))

	err = p.writeAck(readAddr)
	if err != nil {
		return err
	}
	readLen := byte(len(buf) - 1)
	toWrite := []byte{readLen, crc([]byte{readLen}, 0xFF)}
	err = p.writeAck(toWrite)
	if err != nil {
		return err
	}
	err = p.read(buf)
	return err
}

func (p *SerialFlasher) writeMemory(buf []byte, addr uint32) error {
	if len(buf) > 256 {
		return fmt.Errorf("maximum size of buffer is 256, got: %d", len(buf))
	}
	err := p.writeCmdAck(WRITE_MEMORY)
	if err != nil {
		return fmt.Errorf("cmd: %v", err)
	}
	addrB := make([]byte, 4)
	binary.BigEndian.PutUint32(addrB, addr)
	addrB = append(addrB, crc(addrB, 0))
	err = p.writeAck(addrB)
	if err != nil {
		return fmt.Errorf("addr: %v", err)
	}
	toWrite := append([]byte{byte(len(buf) - 1)}, buf...)
	toWrite = append(toWrite, crc(toWrite, 0x00))
	err = p.writeAck(toWrite)
	if err != nil {
		return fmt.Errorf("data at 0x%x: %v", addr, err)
	}
	return nil
}

func (p *SerialFlasher) EraseGlobal() error {
	err := p.writeCmdAck(ERASE)
	if err != nil {
		return fmt.Errorf("cmd: %v", err)
	}
	err = p.writeAck([]byte{0xFF})
	if err != nil {
		return fmt.Errorf("EraseGlobal: %v", err)
	}
	return nil
}

func (p *SerialFlasher) ExtendedEraseGlobal() error {
	err := p.writeCmdAck(EXTENDED_ERASE)
	if err != nil {
		return fmt.Errorf("cmd: %v", err)
	}
	err = p.writeAck([]byte{0xFF, 0xFF, 0x00})
	if err != nil {
		return fmt.Errorf("ExtendedEraseGlobal: %v", err)
	}
	return nil
}

func (p *SerialFlasher) WriteProtect(sectors []byte) error {
	// err := p.writeCmdAck(WRITE_PROTECT)
	// if err != nil {
	// 	return fmt.Errorf("cmd: %v", err)
	// }
	return nil

}
func (p *SerialFlasher) WriteUnprotect() error {
	err := p.writeCmdAck(WRITE_UNPROTECT)
	if err != nil {
		return err
	}
	return p.writeAck([]byte{})
}
func (p *SerialFlasher) ReadoutProtect() error {
	err := p.writeCmdAck(READOUT_PROTECT)
	if err != nil {
		return fmt.Errorf("cmd: %v", err)
	}
	return p.writeAck([]byte{})
}
func (p *SerialFlasher) ReadoutUnprotect() error {
	err := p.writeCmdAck(READOUT_UNPROTECT)
	if err != nil {
		return fmt.Errorf("cmd: %v", err)
	}
	return p.writeAck([]byte{})
}
