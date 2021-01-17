[![Go Report Card](https://goreportcard.com/badge/github.com/PiotrKozimor/stmprog?style=flat-square)](https://goreportcard.com/report/github.com/PiotrKozimor/stmprog)
[![LICENSE](https://img.shields.io/github/license/PiotrKozimor/stmprog.svg?style=flat-square)](https://github.com/PiotrKozimor/stmprog/blob/master/LICENSE)
[![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/PiotrKozimor/stmprog)
[![Releases](https://img.shields.io/github/release/PiotrKozimor/stmprog.svg?style=flat-square)](https://github.com/PiotrKozimor/stmprog/releases)
# stmprog

Stmprog is implementation of [USART protocol used in the STM32 bootloader](https://www.st.com/resource/en/application_note/cd00264342-usart-protocol-used-in-the-stm32-bootloader-stmicroelectronics.pdf) in Golang.

## Usage
Please see [example](example/main.go). Assumptions:
 - Boot0 pin of STM connected to pin 65 in RPi Compute Module (pin 19 according to BCM enumeration, please consult [datasheet](https://www.raspberrypi.org/documentation/hardware/computemodule/datasheets/rpi_DATA_CM3plus_1p0.pdf), chapter 5),
 - Rst pin of STM connected to pin 65 in compute module (pin 18 according to BCM enumeration),
 - `115200` baundrate,
 - `/dev/ttyUSB0` serial port,
 - `app_firmware.bin` application file,
 - `0x8000000` application address