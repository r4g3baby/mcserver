package server

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/klauspost/compress/zlib"
	"github.com/r4g3baby/mcserver/pkg/log"
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/protocol/packets"
	"github.com/r4g3baby/mcserver/pkg/util"
	"github.com/r4g3baby/mcserver/pkg/util/chat"
	"github.com/r4g3baby/mcserver/pkg/util/pools"
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
		SetCompressionLevel(level int)
		GetCompressionLevel() int

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
		compression compression
	}

	compression struct {
		mutex     sync.RWMutex
		enabled   bool
		threshold int
		level     int
	}
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
	conn.compression.mutex.RLock()
	defer conn.compression.mutex.RUnlock()
	return conn.compression.enabled
}

func (conn *connection) setCompressionThreshold(threshold int) {
	conn.compression.mutex.Lock()
	defer conn.compression.mutex.Unlock()
	conn.compression.threshold = threshold
	conn.compression.enabled = threshold >= 0
}

func (conn *connection) GetCompressionThreshold() int {
	conn.compression.mutex.RLock()
	defer conn.compression.mutex.RUnlock()
	return conn.compression.threshold
}

func (conn *connection) SetCompressionLevel(level int) {
	conn.compression.mutex.Lock()
	defer conn.compression.mutex.Unlock()
	conn.compression.level = level
}

func (conn *connection) GetCompressionLevel() int {
	conn.compression.mutex.RLock()
	defer conn.compression.mutex.RUnlock()
	return conn.compression.level
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
	if _, err := io.ReadFull(conn, payload); err != nil {
		return err
	}

	packetData := pools.Buffer.Get(payload)
	defer pools.Buffer.Put(packetData)

	if conn.UseCompression() {
		dataLength, err := packetData.ReadVarInt()
		if err != nil {
			return err
		}

		if dataLength > int32(conn.GetCompressionThreshold()) {
			return errors.New("fat")
		} else if dataLength > 0 {
			var uncompressedData = make([]byte, length)
			if err := func(zlibReader io.ReadCloser, err error) error {
				if err != nil {
					return err
				}
				defer pools.Zlib.PutReader(zlibReader)
				if _, err := io.ReadFull(zlibReader, uncompressedData); err != nil {
					return err
				}
				return nil
			}(pools.Zlib.GetReader(packetData)); err != nil {
				return err
			}

			packetData.Reset()
			_, _ = packetData.Write(uncompressedData)
		}
	}

	packetID, err := packetData.ReadVarInt()
	if err != nil {
		return err
	}

	packet, err := packets.Get(conn.GetProtocol(), conn.GetState(), protocol.ServerBound, packetID)
	if err != nil {
		if debugLog := log.Log.V(1); debugLog.Enabled() {
			debugLog.WithValues(
				"id", fmt.Sprintf("%#0X", packetID),
				"state", conn.GetState(),
				"protocol", int32(conn.GetProtocol()),
				"compression", conn.UseCompression(),
			).V(1).Info("received unknown packet")
		}
		return nil
	}

	if debugLog := log.Log.V(1); debugLog.Enabled() {
		debugLog.WithValues(
			"id", fmt.Sprintf("%#0X", packetID),
			"type", reflect.TypeOf(packet),
			"state", conn.GetState(),
			"protocol", int32(conn.GetProtocol()),
			"compression", conn.UseCompression(),
		).V(1).Info("received packet")
	}

	if err := packet.Read(conn.GetProtocol(), packetData); err != nil {
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
				Threshold: int32(conn.server.GetConfig().Compression.Threshold),
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
				Gamemode:         1,
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

			if err := conn.WritePacket(&packets.PacketPlayOutPositionAndLook{
				X: 0.5,
				Y: 65,
				Z: 0.5,
			}); err != nil {
				return err
			}

			world := NewWorld("overworld")
			renderDistance := 10 // load all chunks in render distance
			for x := -renderDistance; x <= renderDistance; x++ {
				for z := -renderDistance; z <= renderDistance; z++ {
					world.GetChunk(x, z)
				}
			}

			world.SetBlock(0, 65, 0, "minecraft:torch")
			world.SetBlock(0, 64, 0, "minecraft:dirt")
			world.SetBlock(1, 64, 0, "minecraft:stone")
			world.SetBlock(1, 64, 1, "minecraft:stone")
			world.SetBlock(0, 64, 1, "minecraft:stone")
			world.SetBlock(-1, 64, 1, "minecraft:stone")
			world.SetBlock(-1, 64, 0, "minecraft:stone")
			world.SetBlock(-1, 64, -1, "minecraft:stone")
			world.SetBlock(0, 64, -1, "minecraft:stone")
			world.SetBlock(1, 64, -1, "minecraft:stone")

			return world.SendChunks(player)
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

	packetData := pools.Buffer.Get(nil)
	defer pools.Buffer.Put(packetData)

	packetID, err := packets.GetID(conn.GetProtocol(), conn.GetState(), protocol.ClientBound, packet)
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

	buffer := pools.Buffer.Get(nil)
	defer pools.Buffer.Put(buffer)

	if conn.UseCompression() {
		data := pools.Buffer.Get(nil)
		defer pools.Buffer.Put(data)

		if dataLength >= conn.GetCompressionThreshold() {
			if err := data.WriteVarInt(int32(dataLength)); err != nil {
				return err
			}

			level := conn.GetCompressionLevel()
			if err := func(zlibWriter *zlib.Writer, err error) error {
				if err != nil {
					return err
				}
				defer pools.Zlib.PutWriter(zlibWriter, level)
				if _, err := packetData.WriteTo(zlibWriter); err != nil {
					return err
				}
				return nil
			}(pools.Zlib.GetWriter(data, level)); err != nil {
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

	if debugLog := log.Log.V(1); debugLog.Enabled() {
		debugLog.WithValues(
			"id", fmt.Sprintf("%#0X", packetID),
			"type", reflect.TypeOf(packet),
			"state", conn.GetState(),
			"protocol", int32(conn.GetProtocol()),
			"compression", conn.UseCompression(),
		).V(1).Info("sent packet")
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
		compression: compression{
			level: server.GetConfig().Compression.Level,
		},
	}
}
