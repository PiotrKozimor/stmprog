package stmprog

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"periph.io/x/periph/host"
)

func init() {
	logrus.SetLevel(logrus.DebugLevel)
}

func handle(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

func TestInit(t *testing.T) {
	dut, err := NewSerialFlasher(DefaultOptions(TestPort))
	handle(err, t)
	handle(dut.Init(), t)
}

func TestGet_(t *testing.T) {
	dut, err := NewSerialFlasher(DefaultOptions(TestPort))
	handle(err, t)
	bootloader, err := dut.Get()
	handle(err, t)
	t.Logf("bootloader: %+v", bootloader)
}

func TestGetVersionAndReadProtectionStatus(t *testing.T) {
	dut, err := NewSerialFlasher(DefaultOptions(TestPort))
	handle(err, t)
	rdp, err := dut.GetVersionAndReadProtectionStatus()
	handle(err, t)
	t.Logf("rdp: %+v", rdp)
}

func TestReadMemory(t *testing.T) {
	dut, err := NewSerialFlasher(DefaultOptions(TestPort))
	handle(err, t)
	b := make([]byte, 128)
	handle(dut.readMemory(b, APPLICATION_ADDRESS+0x400), t)
	t.Logf("%x", b)
}

func TestWriteMemory(t *testing.T) {
	_, b, err := ReadFile(TEST_APP_FILE)
	handle(err, t)
	dut, err := NewSerialFlasher(DefaultOptions(TestPort))
	handle(err, t)
	handle(dut.writeMemory(b[0:256], APPLICATION_ADDRESS), t)
}

func TestEraseGlobal(t *testing.T) {
	dut, err := NewSerialFlasher(DefaultOptions(TestPort))
	handle(err, t)
	handle(dut.EraseGlobal(), t)
}

func TestExtendedEraseGlobal(t *testing.T) {
	dut, err := NewSerialFlasher(DefaultOptions(TestPort))
	handle(err, t)
	handle(dut.ExtendedEraseGlobal(), t)
}

func TestWriteUnprotect(t *testing.T) {
	dut, err := NewSerialFlasher(DefaultOptions(TestPort))
	handle(err, t)
	handle(dut.WriteUnprotect(), t)
}

func TestReadoutProtect(t *testing.T) {
	dut, err := NewSerialFlasher(DefaultOptions(TestPort))
	handle(err, t)
	handle(dut.ReadoutProtect(), t)
}
func TestReadoutUnprotect(t *testing.T) {
	dut, err := NewSerialFlasher(DefaultOptions(TestPort))
	handle(err, t)
	handle(dut.ReadoutUnprotect(), t)
}

func TestIntegrate(t *testing.T) {
	_, err := host.Init()
	handle(err, t)
	booter := TestBooter
	dut, err := NewSerialFlasher(DefaultOptions(TestPort))
	handle(err, t)
	_, b, err := ReadFile(TEST_APP_FILE)
	handle(err, t)
	handle(booter.Boot(), t)
	handle(dut.Init(), t)
	bootloader, err := dut.Get()
	handle(err, t)
	if bootloader.Version == 0 {
		t.Error("got 0 bootloader version")
	}
	rdp, err := dut.GetVersionAndReadProtectionStatus()
	handle(err, t)
	if rdp.Version == 0 {
		t.Error("got 0 bootloader version")
	}
	handle(dut.writeMemory(b[0:256], APPLICATION_ADDRESS), t)
	read := make([]byte, 128)
	handle(dut.readMemory(read, APPLICATION_ADDRESS), t)
	if bytes.Compare(read, b[0:128]) != 0 {
		log.Fatal("Read data nto equal to programmed")
	}
	handle(dut.ExtendedEraseGlobal(), t)
	handle(dut.WriteUnprotect(), t)
	time.Sleep(time.Millisecond * 100)
	handle(booter.Boot(), t)
	handle(dut.Init(), t)
	handle(dut.ReadoutProtect(), t)
	time.Sleep(time.Millisecond * 100)
	handle(booter.Boot(), t)
	handle(dut.Init(), t)
	handle(dut.ReadoutUnprotect(), t)
}
