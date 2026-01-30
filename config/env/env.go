package env

import "github.com/xuyang-lee/wormhole/model"

// sender
const (
	KeyDstAddr       = "HOLE_DST_ADDR"
	KeyDstIp         = "HOLE_DST_IP"
	KeyDstPort       = "HOLE_DST_PORT"
	KeyPreferredMode = "HOLE_PREFERRED_MODE"
)

// receiver
const (
	KeyHolePort = "HOLE_PORT"
	KeyHoleDir  = "HOLE_DIR"
)

var (
	DstAddr       string
	DstIp         string
	DstPort       int
	PreferredMode model.SendModeFlag
	HolePort      int
	HoleDir       string
)

// default
var (
	DefaultHolePort = 10801
	DefaultHoleDir  = "."
)
