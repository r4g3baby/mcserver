package server

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/protocol/packets"
	"github.com/r4g3baby/mcserver/pkg/util"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
	"github.com/r4g3baby/mcserver/pkg/util/chat"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"reflect"
	"strings"
	"time"
)

type Connection struct {
	net.Conn

	server *Server

	uniqueID uuid.UUID
	username string
	protocol protocol.Protocol
	state    protocol.State
}

func (conn *Connection) Close() error {
	conn.server.removePlayer(conn.uniqueID)
	return conn.Conn.Close()
}

func (conn *Connection) DelayedClose(delay time.Duration) error {
	time.Sleep(delay)
	return conn.Close()
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
		if e := log.Debug(); e.Enabled() {
			e.Str("id", fmt.Sprintf("%#0X", packetID))
			e.Stringer("state", conn.state)
			e.Int32("protocol", int32(conn.protocol))
			e.Msg("received unknown packet")
		}
		return nil
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
			return conn.WritePacket(&packets.PacketStatusOutResponse{
				Response: packets.Response{
					Version: packets.Version{
						Name:     chat.ColorChar + "cHello World!",
						Protocol: int(conn.protocol),
					},
					Players: packets.Players{
						Max:    100,
						Online: 0,
						Sample: []packets.Sample{
							{chat.ColorChar + "bHello World!", uuid.New()},
						},
					},
					Description: []chat.Component{
						&chat.TextComponent{
							Text: "Hello World!\n",
							BaseComponent: chat.BaseComponent{
								Color: &chat.Blue,
							},
						},
						&chat.TextComponent{
							Text: "Hello World!",
							BaseComponent: chat.BaseComponent{
								Color: &chat.Color{Hex: "c31331"},
							},
						},
					},
				},
			})
		case *packets.PacketStatusInPing:
			return conn.WritePacket(&packets.PacketStatusOutPong{
				Payload: p.Payload,
			})
		}
	case protocol.Login:
		switch p := packet.(type) {
		case *packets.PacketLoginInStart:
			conn.uniqueID = util.NameUUIDFromBytes([]byte("OfflinePlayer:" + conn.username))
			conn.username = p.Username

			_ = conn.server.createPlayer(conn)
			if err := conn.WritePacket(&packets.PacketLoginOutSuccess{
				UniqueID: conn.uniqueID,
				Username: conn.username,
			}); err != nil {
				return err
			}
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
	case protocol.Login:
		switch packet.(type) {
		case *packets.PacketLoginOutDisconnect:
			if err := conn.DelayedClose(250 * time.Millisecond); err != nil {
				// See https://github.com/golang/go/issues/4373 for info.
				if !strings.Contains(err.Error(), "use of closed network connection") {
					return err
				}
			}
		case *packets.PacketLoginOutSuccess:
			conn.state = protocol.Play
		}
	case protocol.Play:
		switch packet.(type) {
		case *packets.PacketPlayOutDisconnect:
			if err := conn.DelayedClose(250 * time.Millisecond); err != nil {
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

func NewConnection(conn net.Conn, server *Server) *Connection {
	return &Connection{
		conn,
		server,
		uuid.Nil,
		"",
		protocol.Unknown,
		protocol.Handshaking,
	}
}
