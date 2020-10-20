package server

import (
	"errors"
	"fmt"
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/protocol/packets"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"reflect"
	"strconv"
	"strings"
)

type Connection struct {
	net.Conn

	protocol protocol.Protocol
	state    protocol.State
}

func (conn *Connection) ReadPacket() error {
	length, err := conn.readLength()
	if err != nil {
		return err
	}

	if length < 0 || length > 2147483647 {
		return errors.New("received invalid packet length")
	}

	var payload = make([]byte, length)
	if _, err = io.ReadFull(conn, payload); err != nil {
		return err
	}

	packetData := bytes.NewBuffer(payload)
	packetID, err := packetData.ReadVarInt()
	if err != nil {
		return err
	}

	packet, err := packets.GetPacketByID(conn.state, protocol.ServerBound, packetID)
	if err != nil {
		return err
	}

	if e := log.Debug(); e.Enabled() {
		e.Str("id", fmt.Sprintf("%#0X", packetID))
		e.Stringer("type", reflect.TypeOf(packet))
		e.Stringer("state", conn.state)
		e.Int32("protocol", int32(conn.protocol))
		e.Msg("received packet")
	}

	if err = packet.Read(packetData); err != nil {
		return err
	}

	return conn.handlePacketRead(packet)
}

func (conn *Connection) handlePacketRead(packet protocol.Packet) error {
	switch conn.state {
	case protocol.Handshaking:
		switch p := packet.(type) {
		case *packets.PacketHandshakingStart:
			conn.protocol = protocol.Protocol(p.ProtocolVersion)

			switch p.NextState {
			case 1:
				conn.state = protocol.Status
			case 2:
				conn.state = protocol.Login
			default:
				return errors.New("received invalid nextState")
			}
		}
	case protocol.Status:
		switch p := packet.(type) {
		case *packets.PacketStatusInRequest:
			proto := strconv.Itoa(int(conn.protocol))
			return conn.WritePacket(&packets.PacketStatusOutResponse{
				Response: `{"version":{"name":"§cHello World!","protocol":` + proto + `},"players":{"max":100,"online":0,"sample":[]},"description":{"text":"§9Hello World!"}}`,
			})
		case *packets.PacketStatusInPing:
			return conn.WritePacket(&packets.PacketStatusOutPong{
				Payload: p.Payload,
			})
		}
	}
	return nil
}

func (conn *Connection) WritePacket(packet protocol.Packet) error {
	packetData := bytes.NewBuffer(nil)

	if err := packetData.WriteVarInt(packet.GetID()); err != nil {
		return err
	}

	if err := packet.Write(packetData); err != nil {
		return err
	}

	buffer := bytes.NewBuffer(nil)
	if err := buffer.WriteVarInt(int32(packetData.Len())); err != nil {
		return err
	}

	if _, err := packetData.WriteTo(buffer); err != nil {
		return err
	}

	if _, err := buffer.WriteTo(conn); err != nil {
		return err
	}

	if e := log.Debug(); e.Enabled() {
		e.Str("id", fmt.Sprintf("%#0X", packet.GetID()))
		e.Stringer("type", reflect.TypeOf(packet))
		e.Stringer("state", conn.state)
		e.Int32("protocol", int32(conn.protocol))
		e.Msg("sent packet")
	}

	return conn.handlePacketWrite(packet)
}

func (conn *Connection) handlePacketWrite(packet protocol.Packet) error {
	switch conn.state {
	case protocol.Status:
		switch packet.(type) {
		case *packets.PacketStatusOutPong:
			if err := conn.Close(); err != nil {
				// See https://github.com/golang/go/issues/4373 for info.
				if !strings.Contains(err.Error(), "use of closed network connection") {
					return err
				}
			}
		}
	}
	return nil
}

func (conn *Connection) readLength() (int32, error) {
	var result int32 = 0
	for numRead := 0; ; numRead++ {
		var read = make([]byte, 1)
		if _, err := conn.Read(read); err != nil {
			return result, err
		}

		result |= int32(read[0]&0x7F) << (7 * numRead)

		if numRead >= 5 {
			return result, errors.New("VarInt too big")
		}

		if (read[0] & 0x80) != 0x80 {
			break
		}
	}
	return result, nil
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn,
		protocol.Unknown,
		protocol.Handshaking,
	}
}
