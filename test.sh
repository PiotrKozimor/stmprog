set -xe
go test -run TestSerialIntegrate
go test -run TestProgramVerify
go test -run TestIntegrate
go test -run TestReset
go test -run TestBoot