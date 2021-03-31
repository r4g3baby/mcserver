package server

import (
	"github.com/google/uuid"
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/protocol/packets"
	"github.com/r4g3baby/mcserver/pkg/util/chat"
	"sync"
	"sync/atomic"
	"time"
)

type (
	Player interface {
		GetServer() Server
		GetUniqueID() uuid.UUID
		GetUsername() string
		GetState() protocol.State
		setLatency(latency time.Duration)
		GetLatency() time.Duration
		setKeepAlivePending(keepAlivePending bool)
		IsKeepAlivePending() bool
		setLastKeepAliveTime(lastKeepAliveTime time.Time)
		GetLastKeepAliveTime() time.Time
		setLastKeepAliveID(lastKeepAliveID int32)
		GetLastKeepAliveID() int32
		SendPacket(packet protocol.Packet) error
		Kick(reason []chat.Component) error
	}

	player struct {
		conn Connection

		latency atomic.Value

		mutex             sync.RWMutex
		keepAlivePending  bool
		lastKeepAliveTime time.Time
		lastKeepAliveID   int32
	}
)

func (player *player) GetServer() Server {
	return player.conn.GetServer()
}

func (player *player) GetUniqueID() uuid.UUID {
	return player.conn.GetUniqueID()
}

func (player *player) GetUsername() string {
	return player.conn.GetUsername()
}

func (player *player) GetProtocol() protocol.Protocol {
	return player.conn.GetProtocol()
}

func (player *player) GetState() protocol.State {
	return player.conn.GetState()
}

func (player *player) setLatency(latency time.Duration) {
	player.latency.Store(latency)
}

func (player *player) GetLatency() time.Duration {
	return player.latency.Load().(time.Duration)
}

func (player *player) setKeepAlivePending(keepAlivePending bool) {
	player.mutex.Lock()
	defer player.mutex.Unlock()
	player.keepAlivePending = keepAlivePending
}

func (player *player) IsKeepAlivePending() bool {
	player.mutex.RLock()
	defer player.mutex.RUnlock()
	return player.keepAlivePending
}

func (player *player) setLastKeepAliveTime(lastKeepAliveTime time.Time) {
	player.mutex.Lock()
	defer player.mutex.Unlock()
	player.lastKeepAliveTime = lastKeepAliveTime
}

func (player *player) GetLastKeepAliveTime() time.Time {
	player.mutex.RLock()
	defer player.mutex.RUnlock()
	return player.lastKeepAliveTime
}

func (player *player) setLastKeepAliveID(lastKeepAliveID int32) {
	player.mutex.Lock()
	defer player.mutex.Unlock()
	player.lastKeepAliveID = lastKeepAliveID
}

func (player *player) GetLastKeepAliveID() int32 {
	player.mutex.RLock()
	defer player.mutex.RUnlock()
	return player.lastKeepAliveID
}

func (player *player) SendPacket(packet protocol.Packet) error {
	return player.conn.WritePacket(packet)
}

func (player *player) Kick(reason []chat.Component) error {
	if player.GetState() == protocol.Handshaking || player.GetState() == protocol.Login {
		return player.SendPacket(&packets.PacketLoginOutDisconnect{
			Reason: reason,
		})
	} else {
		return player.SendPacket(&packets.PacketPlayOutDisconnect{
			Reason: reason,
		})
	}
}

func newPlayer(conn Connection) Player {
	player := &player{
		conn: conn,
	}
	player.setLatency(-1)
	return player
}
