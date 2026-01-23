package env

// sender
const (
	KeyDstAddr = "HOLE_DST_ADDR"
	KeyDstIp   = "HOLE_DST_IP"
	KeyDstPort = "HOLE_DST_PORT"
)

// receiver
const (
	KeyHolePort = "HOLE_PORT"
	KeyHoleDir  = "HOLE_DIR"
)

var (
	DstAddr  string
	DstIp    string
	DstPort  int
	HolePort int
	HoleDir  string
)

// default
var (
	DefaultHolePort = 10801
	DefaultHoleDir  = "."
)
