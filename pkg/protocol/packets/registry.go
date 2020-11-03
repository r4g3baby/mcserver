package packets

import (
	"fmt"
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"reflect"
)

var packetsByID = map[protocol.State]map[protocol.Direction]map[int32]reflect.Type{
	protocol.Handshaking: {
		protocol.ServerBound: {
			0x00: reflect.TypeOf((*PacketHandshakingStart)(nil)).Elem(),
		},
	},
	protocol.Status: {
		protocol.ClientBound: {
			0x00: reflect.TypeOf((*PacketStatusOutResponse)(nil)).Elem(),
			0x01: reflect.TypeOf((*PacketStatusOutPong)(nil)).Elem(),
		},
		protocol.ServerBound: {
			0x00: reflect.TypeOf((*PacketStatusInRequest)(nil)).Elem(),
			0x01: reflect.TypeOf((*PacketStatusInPing)(nil)).Elem(),
		},
	},
	protocol.Login: {
		protocol.ClientBound: {
			0x00: reflect.TypeOf((*PacketLoginOutDisconnect)(nil)).Elem(),
			0x02: reflect.TypeOf((*PacketLoginOutSuccess)(nil)).Elem(),
		},
		protocol.ServerBound: {
			0x00: reflect.TypeOf((*PacketLoginInStart)(nil)).Elem(),
		},
	},
}

func GetPacketByID(state protocol.State, direction protocol.Direction, id int32) (protocol.Packet, error) {
	if directions, ok := packetsByID[state]; ok {
		if pIDs, ok := directions[direction]; ok {
			if pType, ok := pIDs[id]; ok {
				return reflect.New(pType).Interface().(protocol.Packet), nil
			}
		}
	}
	return nil, fmt.Errorf("no packet found with id %d", id)
}
