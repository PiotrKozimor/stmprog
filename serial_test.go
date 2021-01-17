package stmprog

import (
	"bytes"
	"testing"

	"github.com/sirupsen/logrus"
	"periph.io/x/periph/host"
)

const TestPort = "/dev/point2"

func TestWrite_(t *testing.T) {
	// logrus.SetLevel(logrus.DebugLevel)
	logrus.SetLevel(logrus.InfoLevel)
	_, b, err := ReadFile(TEST_APP_FILE)
	handle(err, t)
	dut, err := NewSerialFlasher(DefaultOptions(TestPort))
	handle(err, t)
	n, err := dut.Write(b)
	handle(err, t)
	t.Logf("write %d bytes", n)
}

func TestRead_(t *testing.T) {
	logrus.SetLevel(logrus.InfoLevel)
	dut, err := NewSerialFlasher(DefaultOptions(TestPort))
	handle(err, t)
	n, b, err := ReadFile(TEST_APP_FILE)
	handle(err, t)
	b = make([]byte, n)
	n, err = dut.Read(b)
	handle(err, t)
	t.Logf("read %x", b)
	t.Logf("read %d bytes", n)
}

func TestSerialIntegrate(t *testing.T) {
	logrus.SetLevel(logrus.InfoLevel)
	_, err := host.Init()
	handle(err, t)
	logrus.SetLevel(logrus.InfoLevel)
	booter := TestBooter
	dut, err := NewSerialFlasher(DefaultOptions(TestPort))
	handle(err, t)
	handle(booter.Boot(), t)
	handle(dut.Init(), t)
	n, b, err := ReadFile(TEST_APP_FILE)
	handle(err, t)
	n, err = dut.Write(b)
	handle(err, t)
	read := make([]byte, n)
	n, err = dut.Read(read)
	handle(err, t)
	t.Logf("read %x", read)
	t.Logf("read %d bytes", n)
	t.Logf("write %d bytes", n)
	if bytes.Compare(read, b) != 0 {
		t.Fatal("read data different than written")
	}
}
