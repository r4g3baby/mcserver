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

type Player struct {
	conn *Connection

	latency atomic.Value

	mutex             sync.RWMutex
	keepAlivePending  bool
	lastKeepAliveTime time.Time
	lastKeepAliveID   int32
}

func (player *Player) GetServer() *Server {
	return player.conn.server
}

func (player *Player) GetUniqueID() uuid.UUID {
	return player.conn.GetUniqueID()
}

func (player *Player) GetUsername() string {
	return player.conn.GetUsername()
}

func (player *Player) GetProtocol() protocol.Protocol {
	return player.conn.GetProtocol()
}

func (player *Player) GetState() protocol.State {
	return player.conn.GetState()
}

func (player *Player) setLatency(latency time.Duration) {
	player.latency.Store(latency)
}

func (player *Player) GetLatency() time.Duration {
	return player.latency.Load().(time.Duration)
}

func (player *Player) setKeepAlivePending(keepAlivePending bool) {
	player.mutex.Lock()
	defer player.mutex.Unlock()
	player.keepAlivePending = keepAlivePending
}

func (player *Player) IsKeepAlivePending() bool {
	player.mutex.RLock()
	defer player.mutex.RUnlock()
	return player.keepAlivePending
}

func (player *Player) setLastKeepAliveTime(lastKeepAliveTime time.Time) {
	player.mutex.Lock()
	defer player.mutex.Unlock()
	player.lastKeepAliveTime = lastKeepAliveTime
}

func (player *Player) GetLastKeepAliveTime() time.Time {
	player.mutex.RLock()
	defer player.mutex.RUnlock()
	return player.lastKeepAliveTime
}

func (player *Player) setLastKeepAliveID(lastKeepAliveID int32) {
	player.mutex.Lock()
	defer player.mutex.Unlock()
	player.lastKeepAliveID = lastKeepAliveID
}

func (player *Player) GetLastKeepAliveID() int32 {
	player.mutex.RLock()
	defer player.mutex.RUnlock()
	return player.lastKeepAliveID
}

func (player *Player) SendPacket(packet protocol.Packet) error {
	return player.conn.WritePacket(packet)
}

func (player *Player) Kick(reason []chat.Component) error {
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

func newPlayer(conn *Connection) *Player {
	player := &Player{
		conn: conn,
	}
	player.setLatency(-1)
	return player
}
