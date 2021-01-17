package stmprog

import (
	"testing"

	"github.com/sirupsen/logrus"
	"periph.io/x/periph/host"
)

const TEST_APP_FILE = "app_firmware.bin"
const TEST_DALI_POINT = 2

var testOptions = DefaultOptions(TestPort).WithReadChunk(128)

func TestProgram_(t *testing.T) {
	_, err := host.Init()
	handle(err, t)
	logrus.SetLevel(logrus.InfoLevel)
	flasher, err := NewSerialFlasher(testOptions)
	handle(err, t)
	prog := Programmer{
		Flasher: flasher,
		Booter:  &TestBooter,
	}
	_, b, err := ReadFile(TEST_APP_FILE)
	handle(err, t)
	handle(prog.Program(b), t)
}

func TestProgramVerify(t *testing.T) {
	_, err := host.Init()
	handle(err, t)
	logrus.SetLevel(logrus.InfoLevel)
	flasher, err := NewSerialFlasher(testOptions)
	handle(err, t)
	prog := Programmer{
		Flasher: flasher,
		Booter:  &TestBooter,
	}
	_, b, err := ReadFile(TEST_APP_FILE)
	handle(err, t)
	handle(prog.ProgramVerify(b), t)
}
