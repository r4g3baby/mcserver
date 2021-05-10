package packets

import (
	"errors"
	"fmt"
	"github.com/imdario/mergo"
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"reflect"
)

var (
	packets = map[protocol.Protocol]map[protocol.State]map[protocol.Direction]map[reflect.Type]int32{
		protocol.Unknown: {
			protocol.Handshaking: {
				protocol.ServerBound: {
					reflect.TypeOf((*PacketHandshakingStart)(nil)).Elem(): 0x00,
				},
			},
			protocol.Status: {
				protocol.ClientBound: {
					reflect.TypeOf((*PacketStatusOutResponse)(nil)).Elem(): 0x00,
					reflect.TypeOf((*PacketStatusOutPong)(nil)).Elem():     0x01,
				},
				protocol.ServerBound: {
					reflect.TypeOf((*PacketStatusInRequest)(nil)).Elem(): 0x00,
					reflect.TypeOf((*PacketStatusInPing)(nil)).Elem():    0x01,
				},
			},
			protocol.Login: {
				protocol.ClientBound: {
					reflect.TypeOf((*PacketLoginOutDisconnect)(nil)).Elem():  0x00,
					reflect.TypeOf((*PacketLoginOutSuccess)(nil)).Elem():     0x02,
					reflect.TypeOf((*PacketLoginOutCompression)(nil)).Elem(): 0x03,
				},
				protocol.ServerBound: {
					reflect.TypeOf((*PacketLoginInStart)(nil)).Elem(): 0x00,
				},
			},
		},
	}

	packetsByID = map[protocol.Protocol]map[protocol.State]map[protocol.Direction]map[int32]reflect.Type{
		protocol.Unknown: {
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
		},
	}
)

func init() {
	if err := Register(protocol.V1_8, map[protocol.State]map[protocol.Direction]map[reflect.Type]int32{
		protocol.Play: {
			protocol.ClientBound: {
				reflect.TypeOf((*PacketPlayOutKeepAlive)(nil)).Elem():        0x00,
				reflect.TypeOf((*PacketPlayOutChatMessage)(nil)).Elem():      0x02,
				reflect.TypeOf((*PacketPlayOutJoinGame)(nil)).Elem():         0x01,
				reflect.TypeOf((*PacketPlayOutPositionAndLook)(nil)).Elem():  0x08,
				reflect.TypeOf((*PacketPlayOutDisconnect)(nil)).Elem():       0x40,
				reflect.TypeOf((*PacketPlayOutServerDifficulty)(nil)).Elem(): 0x41,
			},
			protocol.ServerBound: {
				reflect.TypeOf((*PacketPlayInChatMessage)(nil)).Elem(): 0x01,
				reflect.TypeOf((*PacketPlayInKeepAlive)(nil)).Elem():   0x00,
			},
		},
	}); err != nil {
		panic(err)
	}

	if err := Register(protocol.V1_9, map[protocol.State]map[protocol.Direction]map[reflect.Type]int32{
		protocol.Play: {
			protocol.ClientBound: {
				reflect.TypeOf((*PacketPlayOutServerDifficulty)(nil)).Elem(): 0x0D,
				reflect.TypeOf((*PacketPlayOutChatMessage)(nil)).Elem():      0x0F,
				reflect.TypeOf((*PacketPlayOutDisconnect)(nil)).Elem():       0x1A,
				reflect.TypeOf((*PacketPlayOutKeepAlive)(nil)).Elem():        0x1F,
				reflect.TypeOf((*PacketPlayOutJoinGame)(nil)).Elem():         0x23,
				reflect.TypeOf((*PacketPlayOutPositionAndLook)(nil)).Elem():  0x2E,
			},
			protocol.ServerBound: {
				reflect.TypeOf((*PacketPlayInChatMessage)(nil)).Elem(): 0x02,
				reflect.TypeOf((*PacketPlayInKeepAlive)(nil)).Elem():   0x0B,
			},
		},
	}); err != nil {
		panic(err)
	}

	if err := Register(protocol.V1_12, map[protocol.State]map[protocol.Direction]map[reflect.Type]int32{
		protocol.Play: {
			protocol.ClientBound: {
				reflect.TypeOf((*PacketPlayOutServerDifficulty)(nil)).Elem(): 0x0D,
				reflect.TypeOf((*PacketPlayOutChatMessage)(nil)).Elem():      0x0F,
				reflect.TypeOf((*PacketPlayOutDisconnect)(nil)).Elem():       0x1A,
				reflect.TypeOf((*PacketPlayOutKeepAlive)(nil)).Elem():        0x1F,
				reflect.TypeOf((*PacketPlayOutJoinGame)(nil)).Elem():         0x23,
				reflect.TypeOf((*PacketPlayOutPositionAndLook)(nil)).Elem():  0x2E,
			},
			protocol.ServerBound: {
				reflect.TypeOf((*PacketPlayInChatMessage)(nil)).Elem(): 0x03,
				reflect.TypeOf((*PacketPlayInKeepAlive)(nil)).Elem():   0x0C,
			},
		},
	}); err != nil {
		panic(err)
	}

	if err := Register(protocol.V1_12_1, map[protocol.State]map[protocol.Direction]map[reflect.Type]int32{
		protocol.Play: {
			protocol.ClientBound: {
				reflect.TypeOf((*PacketPlayOutServerDifficulty)(nil)).Elem(): 0x0D,
				reflect.TypeOf((*PacketPlayOutChatMessage)(nil)).Elem():      0x0F,
				reflect.TypeOf((*PacketPlayOutDisconnect)(nil)).Elem():       0x1A,
				reflect.TypeOf((*PacketPlayOutKeepAlive)(nil)).Elem():        0x1F,
				reflect.TypeOf((*PacketPlayOutJoinGame)(nil)).Elem():         0x23,
				reflect.TypeOf((*PacketPlayOutPositionAndLook)(nil)).Elem():  0x2F,
			},
			protocol.ServerBound: {
				reflect.TypeOf((*PacketPlayInChatMessage)(nil)).Elem(): 0x02,
				reflect.TypeOf((*PacketPlayInKeepAlive)(nil)).Elem():   0x0B,
			},
		},
	}); err != nil {
		panic(err)
	}

	if err := Register(protocol.V1_13, map[protocol.State]map[protocol.Direction]map[reflect.Type]int32{
		protocol.Play: {
			protocol.ClientBound: {
				reflect.TypeOf((*PacketPlayOutServerDifficulty)(nil)).Elem(): 0x0D,
				reflect.TypeOf((*PacketPlayOutChatMessage)(nil)).Elem():      0x0E,
				reflect.TypeOf((*PacketPlayOutDisconnect)(nil)).Elem():       0x1B,
				reflect.TypeOf((*PacketPlayOutKeepAlive)(nil)).Elem():        0x21,
				reflect.TypeOf((*PacketPlayOutJoinGame)(nil)).Elem():         0x25,
				reflect.TypeOf((*PacketPlayOutPositionAndLook)(nil)).Elem():  0x32,
			},
			protocol.ServerBound: {
				reflect.TypeOf((*PacketPlayInChatMessage)(nil)).Elem(): 0x02,
				reflect.TypeOf((*PacketPlayInKeepAlive)(nil)).Elem():   0x0E,
			},
		},
	}); err != nil {
		panic(err)
	}

	if err := Register(protocol.V1_14, map[protocol.State]map[protocol.Direction]map[reflect.Type]int32{
		protocol.Play: {
			protocol.ClientBound: {
				reflect.TypeOf((*PacketPlayOutServerDifficulty)(nil)).Elem(): 0x0D,
				reflect.TypeOf((*PacketPlayOutChatMessage)(nil)).Elem():      0x0E,
				reflect.TypeOf((*PacketPlayOutDisconnect)(nil)).Elem():       0x1A,
				reflect.TypeOf((*PacketPlayOutKeepAlive)(nil)).Elem():        0x20,
				reflect.TypeOf((*PacketPlayOutJoinGame)(nil)).Elem():         0x25,
				reflect.TypeOf((*PacketPlayOutPositionAndLook)(nil)).Elem():  0x35,
			},
			protocol.ServerBound: {
				reflect.TypeOf((*PacketPlayInChatMessage)(nil)).Elem(): 0x03,
				reflect.TypeOf((*PacketPlayInKeepAlive)(nil)).Elem():   0x0F,
			},
		},
	}); err != nil {
		panic(err)
	}

	if err := Register(protocol.V1_15, map[protocol.State]map[protocol.Direction]map[reflect.Type]int32{
		protocol.Play: {
			protocol.ClientBound: {
				reflect.TypeOf((*PacketPlayOutServerDifficulty)(nil)).Elem(): 0x0E,
				reflect.TypeOf((*PacketPlayOutChatMessage)(nil)).Elem():      0x0F,
				reflect.TypeOf((*PacketPlayOutDisconnect)(nil)).Elem():       0x1B,
				reflect.TypeOf((*PacketPlayOutKeepAlive)(nil)).Elem():        0x21,
				reflect.TypeOf((*PacketPlayOutJoinGame)(nil)).Elem():         0x26,
				reflect.TypeOf((*PacketPlayOutPositionAndLook)(nil)).Elem():  0x36,
			},
			protocol.ServerBound: {
				reflect.TypeOf((*PacketPlayInChatMessage)(nil)).Elem(): 0x03,
				reflect.TypeOf((*PacketPlayInKeepAlive)(nil)).Elem():   0x0F,
			},
		},
	}); err != nil {
		panic(err)
	}

	if err := Register(protocol.V1_16, map[protocol.State]map[protocol.Direction]map[reflect.Type]int32{
		protocol.Play: {
			protocol.ClientBound: {
				reflect.TypeOf((*PacketPlayOutServerDifficulty)(nil)).Elem(): 0x0D,
				reflect.TypeOf((*PacketPlayOutChatMessage)(nil)).Elem():      0x0E,
				reflect.TypeOf((*PacketPlayOutDisconnect)(nil)).Elem():       0x1A,
				reflect.TypeOf((*PacketPlayOutKeepAlive)(nil)).Elem():        0x20,
				reflect.TypeOf((*PacketPlayOutJoinGame)(nil)).Elem():         0x25,
				reflect.TypeOf((*PacketPlayOutPositionAndLook)(nil)).Elem():  0x35,
			},
			protocol.ServerBound: {
				reflect.TypeOf((*PacketPlayInChatMessage)(nil)).Elem(): 0x03,
				reflect.TypeOf((*PacketPlayInKeepAlive)(nil)).Elem():   0x10,
			},
		},
	}); err != nil {
		panic(err)
	}

	if err := Register(protocol.V1_16_2, map[protocol.State]map[protocol.Direction]map[reflect.Type]int32{
		protocol.Play: {
			protocol.ClientBound: {
				reflect.TypeOf((*PacketPlayOutServerDifficulty)(nil)).Elem(): 0x0D,
				reflect.TypeOf((*PacketPlayOutChatMessage)(nil)).Elem():      0x0E,
				reflect.TypeOf((*PacketPlayOutDisconnect)(nil)).Elem():       0x19,
				reflect.TypeOf((*PacketPlayOutKeepAlive)(nil)).Elem():        0x1F,
				reflect.TypeOf((*PacketPlayOutChunkData)(nil)).Elem():        0x20,
				reflect.TypeOf((*PacketPlayOutJoinGame)(nil)).Elem():         0x24,
				reflect.TypeOf((*PacketPlayOutPositionAndLook)(nil)).Elem():  0x34,
			},
			protocol.ServerBound: {
				reflect.TypeOf((*PacketPlayInChatMessage)(nil)).Elem(): 0x03,
				reflect.TypeOf((*PacketPlayInKeepAlive)(nil)).Elem():   0x10,
			},
		},
	}); err != nil {
		panic(err)
	}

	Copy(protocol.V1_9, protocol.V1_9_1, protocol.V1_9_2, protocol.V1_9_3, protocol.V1_10, protocol.V1_11, protocol.V1_11_1)
	Copy(protocol.V1_12_1, protocol.V1_12_2)
	Copy(protocol.V1_13, protocol.V1_13_1, protocol.V1_13_2)
	Copy(protocol.V1_14, protocol.V1_14_1, protocol.V1_14_2, protocol.V1_14_3, protocol.V1_14_4)
	Copy(protocol.V1_15, protocol.V1_15_1, protocol.V1_15_2)
	Copy(protocol.V1_16, protocol.V1_16_1)
	Copy(protocol.V1_16_2, protocol.V1_16_3, protocol.V1_16_4)
}

func GetID(proto protocol.Protocol, state protocol.State, direction protocol.Direction, packet protocol.Packet) (int32, error) {
	if states, ok := packets[proto]; ok {
		if directions, ok := states[state]; ok {
			if pTypes, ok := directions[direction]; ok {
				if id, ok := pTypes[reflect.TypeOf(packet).Elem()]; ok {
					return id, nil
				}
			}
		}
	} else {
		return GetID(protocol.Unknown, state, direction, packet)
	}
	return 0, errors.New("no packet id found for the given options")
}

func Get(proto protocol.Protocol, state protocol.State, direction protocol.Direction, id int32) (protocol.Packet, error) {
	if states, ok := packetsByID[proto]; ok {
		if directions, ok := states[state]; ok {
			if pIDs, ok := directions[direction]; ok {
				if pType, ok := pIDs[id]; ok {
					return reflect.New(pType).Interface().(protocol.Packet), nil
				}
			}
		}
	} else {
		return Get(protocol.Unknown, state, direction, id)
	}
	return nil, fmt.Errorf("no packet found with id %d", id)
}

func Register(proto protocol.Protocol, packetsMap map[protocol.State]map[protocol.Direction]map[reflect.Type]int32) error {
	var newPacketsMap map[protocol.State]map[protocol.Direction]map[reflect.Type]int32
	if currentMap, ok := packets[proto]; ok {
		newPacketsMap = currentMap
	} else {
		newPacketsMap = make(map[protocol.State]map[protocol.Direction]map[reflect.Type]int32)
		for state, directions := range packets[protocol.Unknown] {
			newPacketsMap[state] = make(map[protocol.Direction]map[reflect.Type]int32)
			for direction, pTypes := range directions {
				newPacketsMap[state][direction] = make(map[reflect.Type]int32)
				for pType, pID := range pTypes {
					newPacketsMap[state][direction][pType] = pID
				}
			}
		}
	}

	if err := mergo.Merge(&newPacketsMap, packetsMap, mergo.WithOverride); err != nil {
		return err
	}

	packets[proto] = newPacketsMap

	packetsByID[proto] = make(map[protocol.State]map[protocol.Direction]map[int32]reflect.Type)
	for state, directions := range newPacketsMap {
		packetsByID[proto][state] = make(map[protocol.Direction]map[int32]reflect.Type)
		for direction, pTypes := range directions {
			packetsByID[proto][state][direction] = make(map[int32]reflect.Type)
			for pType, pID := range pTypes {
				packetsByID[proto][state][direction][pID] = pType
			}
		}
	}

	return nil
}

func Copy(src protocol.Protocol, destinations ...protocol.Protocol) {
	for _, dst := range destinations {
		packets[dst] = packets[src]
		packetsByID[dst] = packetsByID[src]
	}
}
