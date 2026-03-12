module github.com/SamW94/src/client

go 1.25.2

replace github.com/SamW94/get-lyrics/tracks => ../tracks

require github.com/SamW94/get-lyrics/tracks v0.0.0-00010101000000-000000000000

require (
	github.com/hbollon/go-edlib v1.7.0
	github.com/tetratelabs/wazero v1.10.1 // indirect
	go.senan.xyz/taglib v0.11.1 // indirect
)
