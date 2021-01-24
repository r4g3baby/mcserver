package protocol

import "github.com/r4g3baby/mcserver/pkg/util/bytes"

type (
	State     int
	Direction int
	Protocol  int32

	Packet interface {
		GetID(proto Protocol) (int32, error)
		Read(proto Protocol, buffer *bytes.Buffer) error
		Write(proto Protocol, buffer *bytes.Buffer) error
	}
)

const (
	Handshaking State = iota
	Status
	Login
	Play
)

func (state State) String() string {
	switch state {
	case Handshaking:
		return "Handshaking"
	case Status:
		return "Status"
	case Login:
		return "Login"
	case Play:
		return "Play"
	default:
		return "Unknown"
	}
}

const (
	ClientBound Direction = iota
	ServerBound
)

func (direction Direction) String() string {
	switch direction {
	case ClientBound:
		return "ClientBound"
	case ServerBound:
		return "ServerBound"
	default:
		return "Unknown"
	}
}

const (
	Unknown Protocol = -1
	V1_8    Protocol = 47
	V1_9    Protocol = 107
	V1_9_1  Protocol = 108
	V1_9_2  Protocol = 109
	V1_9_3  Protocol = 110
	V1_10   Protocol = 210
	V1_11   Protocol = 315
	V1_11_1 Protocol = 316
	V1_12   Protocol = 335
	V1_12_1 Protocol = 338
	V1_12_2 Protocol = 340
	V1_13   Protocol = 393
	V1_13_1 Protocol = 401
	V1_13_2 Protocol = 404
	V1_14   Protocol = 477
	V1_14_1 Protocol = 480
	V1_14_2 Protocol = 485
	V1_14_3 Protocol = 490
	V1_14_4 Protocol = 498
	V1_15   Protocol = 573
	V1_15_1 Protocol = 575
	V1_15_2 Protocol = 578
	V1_16   Protocol = 735
	V1_16_1 Protocol = 736
	V1_16_2 Protocol = 751
	V1_16_3 Protocol = 753
	V1_16_4 Protocol = 754
)

var (
	SupportedProtocols = []Protocol{
		V1_8,
		V1_9, V1_9_1, V1_9_2, V1_9_3,
		V1_10,
		V1_11, V1_11_1,
		V1_12, V1_12_1, V1_12_2,
		V1_13, V1_13_1, V1_13_2,
		V1_14, V1_14_1, V1_14_2, V1_14_3, V1_14_4,
		V1_15, V1_15_1, V1_15_2,
		V1_16, V1_16_1, V1_16_2, V1_16_3, V1_16_4,
	}
	LowestProtocol  = SupportedProtocols[0]
	HighestProtocol = SupportedProtocols[len(SupportedProtocols)-1]
)

func IsSupported(protocol Protocol) bool {
	for _, x := range SupportedProtocols {
		if x == protocol {
			return true
		}
	}
	return false
}
