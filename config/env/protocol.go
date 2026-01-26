package env

import "strings"

const (
	ProtocolSplitMark = "|"
)

var (
	ProtocolHelloMsg string = strings.Join([]string{ServerName, Version}, ProtocolSplitMark)
	ProtocolHiMsg    string = strings.Join([]string{ServerName, Version}, ProtocolSplitMark)

	ProtocolHelloSize uint16 = uint16(len(ProtocolHelloMsg))
	ProtocolHiSize    uint16 = uint16(len(ProtocolHiMsg))
)
