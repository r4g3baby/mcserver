package server

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/klauspost/compress/zlib"
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/protocol/packets"
	"github.com/r4g3baby/mcserver/pkg/util"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
	"github.com/r4g3baby/mcserver/pkg/util/chat"
	"github.com/r4g3baby/mcserver/pkg/util/nbt"
	"github.com/rs/zerolog/log"
	"io"
	"net"
	"reflect"
	"sync"
	"time"
)

type (
	Connection interface {
		RemoteAddr() net.Addr

		GetServer() Server
		SetUniqueID(uniqueID uuid.UUID)
		GetUniqueID() uuid.UUID
		SetUsername(username string)
		GetUsername() string
		SetProtocol(protocol protocol.Protocol)
		GetProtocol() protocol.Protocol
		SetState(state protocol.State)
		GetState() protocol.State
		UseCompression() bool
		setCompressionThreshold(threshold int)
		GetCompressionThreshold() int

		Close() error
		DelayedClose(delay time.Duration) error

		ReadPacket() error
		WritePacket(packet protocol.Packet) error
	}

	connection struct {
		net.Conn

		server Server

		mutex       sync.RWMutex
		uniqueID    uuid.UUID
		username    string
		protocol    protocol.Protocol
		state       protocol.State
		compression struct {
			enabled   bool
			threshold int
		}
	}
)

var (
	zlibReaders = sync.Pool{New: func() interface{} {
		reader, _ := zlib.NewReader(nil)
		return reader
	}}

	zlibWriters = sync.Pool{New: func() interface{} {
		return zlib.NewWriter(nil)
	}}
)

func (conn *connection) GetServer() Server {
	return conn.server
}

func (conn *connection) SetUniqueID(uniqueID uuid.UUID) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.uniqueID = uniqueID
}

func (conn *connection) GetUniqueID() uuid.UUID {
	conn.mutex.RLock()
	defer conn.mutex.RUnlock()
	return conn.uniqueID
}

func (conn *connection) SetUsername(username string) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.username = username
}

func (conn *connection) GetUsername() string {
	conn.mutex.RLock()
	defer conn.mutex.RUnlock()
	return conn.username
}

func (conn *connection) SetProtocol(protocol protocol.Protocol) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.protocol = protocol
}

func (conn *connection) GetProtocol() protocol.Protocol {
	conn.mutex.RLock()
	defer conn.mutex.RUnlock()
	return conn.protocol
}

func (conn *connection) SetState(state protocol.State) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.state = state
}

func (conn *connection) GetState() protocol.State {
	conn.mutex.RLock()
	defer conn.mutex.RUnlock()
	return conn.state
}

func (conn *connection) UseCompression() bool {
	conn.mutex.RLock()
	defer conn.mutex.RUnlock()
	return conn.compression.enabled
}

func (conn *connection) setCompressionThreshold(threshold int) {
	conn.mutex.Lock()
	defer conn.mutex.Unlock()
	conn.compression.threshold = threshold
	conn.compression.enabled = threshold >= 0
}

func (conn *connection) GetCompressionThreshold() int {
	conn.mutex.RLock()
	defer conn.mutex.RUnlock()
	return conn.compression.threshold
}

func (conn *connection) Close() error {
	conn.server.removePlayer(conn.GetUniqueID())
	return conn.Conn.Close()
}

func (conn *connection) DelayedClose(delay time.Duration) error {
	time.Sleep(delay)
	return conn.Close()
}

func (conn *connection) ReadPacket() error {
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
	if conn.UseCompression() {
		dataLength, err := packetData.ReadVarInt()
		if err != nil {
			return err
		}

		if dataLength > int32(conn.GetCompressionThreshold()) {
			return errors.New("fat")
		} else if dataLength > 0 {
			var uncompressedData = make([]byte, length)
			if err := func(zlibReader io.ReadCloser) error {
				defer zlibReaders.Put(zlibReader)
				if err = zlibReader.(zlib.Resetter).Reset(packetData, nil); err != nil {
					return err
				}
				if _, err = io.ReadFull(zlibReader, uncompressedData); err != nil {
					return err
				}
				if err = zlibReader.Close(); err != nil {
					return err
				}
				return nil
			}(zlibReaders.Get().(io.ReadCloser)); err != nil {
				return err
			}

			packetData = bytes.NewBuffer(uncompressedData)
		}
	}

	packetID, err := packetData.ReadVarInt()
	if err != nil {
		return err
	}

	packet, err := packets.GetPacketByID(conn.GetProtocol(), conn.GetState(), protocol.ServerBound, packetID)
	if err != nil {
		if e := log.Debug(); e.Enabled() {
			e.Str("id", fmt.Sprintf("%#0X", packetID))
			e.Stringer("state", conn.GetState())
			e.Int32("protocol", int32(conn.GetProtocol()))
			e.Bool("compression", conn.UseCompression())
			e.Msg("received unknown packet")
		}
		return nil
	}

	if e := log.Debug(); e.Enabled() {
		e.Str("id", fmt.Sprintf("%#0X", packetID))
		e.Stringer("type", reflect.TypeOf(packet))
		e.Stringer("state", conn.GetState())
		e.Int32("protocol", int32(conn.GetProtocol()))
		e.Bool("compression", conn.UseCompression())
		e.Msg("received packet")
	}

	if err = packet.Read(conn.GetProtocol(), packetData); err != nil {
		return err
	}

	player := conn.server.GetPlayer(conn.GetUniqueID())
	event := NewPacketEvent(conn, player, packet)
	conn.server.FireEvent(OnPacketReadEvent, event)
	if event.IsCancelled() {
		return nil
	}

	return conn.handlePacketRead(packet)
}

func (conn *connection) handlePacketRead(packet protocol.Packet) error {
	switch conn.GetState() {
	case protocol.Handshaking:
		switch p := packet.(type) {
		case *packets.PacketHandshakingStart:
			conn.SetProtocol(protocol.Protocol(p.ProtocolVersion))

			switch p.NextState {
			case 1:
				conn.SetState(protocol.Status)
			case 2:
				conn.SetState(protocol.Login)
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
						Protocol: int(conn.GetProtocol()),
					},
					Players: packets.Players{
						Max:    conn.server.GetPlayerCount(),
						Online: conn.server.GetPlayerCount(),
						Sample: []packets.Sample{
							{chat.ColorChar + "bHello World!", uuid.Nil},
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
								Color: &chat.Color{Hex: "c33131"},
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
			conn.SetUsername(p.Username)
			conn.SetUniqueID(util.NameUUIDFromBytes([]byte("OfflinePlayer:" + conn.GetUsername())))

			player, online := conn.server.createPlayer(conn)
			if online {
				return conn.WritePacket(&packets.PacketLoginOutDisconnect{
					Reason: []chat.Component{
						&chat.TextComponent{
							Text: "You are already connected to this server!",
							BaseComponent: chat.BaseComponent{
								Color: &chat.Red,
							},
						},
					},
				})
			}

			if err := conn.WritePacket(&packets.PacketLoginOutCompression{
				Threshold: 256,
			}); err != nil {
				return err
			}

			if err := conn.WritePacket(&packets.PacketLoginOutSuccess{
				UniqueID: player.GetUniqueID(),
				Username: player.GetUsername(),
			}); err != nil {
				return err
			}

			if err := conn.WritePacket(&packets.PacketPlayOutJoinGame{
				EntityID:         1,
				Hardcore:         false,
				Gamemode:         0,
				PreviousGamemode: -1,
				WorldNames:       []string{"minecraft:overworld"},
				DimensionCodec:   protocol.DefaultDimensionCodec,
				Dimension:        protocol.Overworld,
				WorldName:        "minecraft:overworld",
				DimensionID:      0,
				HashedSeed:       0,
				MaxPlayers:       20,
				LevelType:        "default",
				ViewDistance:     10,
				ReducedDebug:     false,
				RespawnScreen:    true,
				IsDebug:          false,
				IsFlat:           false,
			}); err != nil {
				return err
			}

			if err := conn.WritePacket(&packets.PacketPlayOutServerDifficulty{
				Difficulty: 1,
				Locked:     true,
			}); err != nil {
				return err
			}

			if err := conn.WritePacket(&packets.PacketPlayOutPositionAndLook{}); err != nil {
				return err
			}

			var biomes []int32
			for i := 0; i < 1024; i++ {
				biomes = append(biomes, 1)
			}

			return conn.WritePacket(&packets.PacketPlayOutChunkData{
				ChunkX:        0,
				ChunkZ:        0,
				FullChunk:     true,
				PrimaryBit:    0,
				Heightmaps:    nbt.CompoundTag{},
				Biomes:        biomes,
				Data:          []byte{},
				BlockEntities: []nbt.Tag{},
			})
		}
	case protocol.Play:
		switch p := packet.(type) {
		case *packets.PacketPlayInKeepAlive:
			if player := conn.server.GetPlayer(conn.GetUniqueID()); player != nil {
				if player.IsKeepAlivePending() && p.KeepAliveID == player.GetLastKeepAliveID() {
					player.setLatency(time.Since(player.GetLastKeepAliveTime()))
					player.setKeepAlivePending(false)
				}
			}
		}
	}
	return nil
}

func (conn *connection) WritePacket(packet protocol.Packet) error {
	player := conn.server.GetPlayer(conn.GetUniqueID())
	event := NewPacketEvent(conn, player, packet)
	conn.server.FireEvent(OnPacketWriteEvent, event)
	if event.IsCancelled() {
		return nil
	}

	packetData := bytes.NewBuffer(nil)

	packetID, err := packets.GetPacketID(conn.GetProtocol(), conn.GetState(), protocol.ClientBound, packet)
	if err != nil {
		return err
	}

	if err := packetData.WriteVarInt(packetID); err != nil {
		return err
	}

	if err := packet.Write(conn.GetProtocol(), packetData); err != nil {
		return err
	}

	dataLength := packetData.Len()

	buffer := bytes.NewBuffer(nil)
	if conn.UseCompression() {
		data := bytes.NewBuffer(nil)
		if dataLength >= conn.GetCompressionThreshold() {
			if err := data.WriteVarInt(int32(dataLength)); err != nil {
				return err
			}

			if err := func(zlibWriter *zlib.Writer) error {
				defer zlibWriters.Put(zlibWriter)
				zlibWriter.Reset(data)
				if _, err := packetData.WriteTo(zlibWriter); err != nil {
					return err
				}
				if err := zlibWriter.Close(); err != nil {
					return err
				}
				return nil
			}(zlibWriters.Get().(*zlib.Writer)); err != nil {
				return err
			}
		} else {
			if err := data.WriteVarInt(0); err != nil {
				return err
			}

			if _, err := packetData.WriteTo(data); err != nil {
				return err
			}
		}

		if err := buffer.WriteVarInt(int32(data.Len())); err != nil {
			return err
		}

		if _, err := data.WriteTo(buffer); err != nil {
			return err
		}
	} else {
		if err := buffer.WriteVarInt(int32(dataLength)); err != nil {
			return err
		}

		if _, err := packetData.WriteTo(buffer); err != nil {
			return err
		}
	}

	if err := conn.handlePrePacketWrite(packet); err != nil {
		return err
	}

	if _, err := buffer.WriteTo(conn); err != nil {
		return err
	}

	if e := log.Debug(); e.Enabled() {
		e.Str("id", fmt.Sprintf("%#0X", packetID))
		e.Stringer("type", reflect.TypeOf(packet))
		e.Stringer("state", conn.GetState())
		e.Int32("protocol", int32(conn.GetProtocol()))
		e.Bool("compression", conn.UseCompression())
		e.Msg("sent packet")
	}

	return conn.handlePostPacketWrite(packet)
}

func (conn *connection) handlePrePacketWrite(packet protocol.Packet) error {
	switch conn.GetState() {
	case protocol.Play:
		switch p := packet.(type) {
		case *packets.PacketPlayOutKeepAlive:
			if player := conn.server.GetPlayer(conn.GetUniqueID()); player != nil {
				player.setLastKeepAliveTime(time.Now())
				player.setLastKeepAliveID(p.KeepAliveID)
				player.setKeepAlivePending(true)
			}
		}
	}
	return nil
}

func (conn *connection) handlePostPacketWrite(packet protocol.Packet) error {
	switch conn.GetState() {
	case protocol.Status:
		switch packet.(type) {
		case *packets.PacketStatusOutPong:
			if err := conn.Close(); err != nil {
				if !errors.Is(err, net.ErrClosed) {
					return err
				}
			}
		}
	case protocol.Login:
		switch p := packet.(type) {
		case *packets.PacketLoginOutDisconnect:
			if err := conn.DelayedClose(250 * time.Millisecond); err != nil {
				if !errors.Is(err, net.ErrClosed) {
					return err
				}
			}
		case *packets.PacketLoginOutSuccess:
			conn.SetState(protocol.Play)
		case *packets.PacketLoginOutCompression:
			conn.setCompressionThreshold(int(p.Threshold))
		}
	case protocol.Play:
		switch packet.(type) {
		case *packets.PacketPlayOutDisconnect:
			if err := conn.DelayedClose(250 * time.Millisecond); err != nil {
				if !errors.Is(err, net.ErrClosed) {
					return err
				}
			}
		}
	}
	return nil
}

func (conn *connection) readLength() (int32, error) {
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

func newConnection(conn net.Conn, server Server) Connection {
	return &connection{
		Conn:     conn,
		server:   server,
		uniqueID: uuid.Nil,
		protocol: protocol.Unknown,
		state:    protocol.Handshaking,
	}
}
