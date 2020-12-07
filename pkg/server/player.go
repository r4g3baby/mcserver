package server

import (
	"github.com/google/uuid"
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/protocol/packets"
	"sync"
)

type Player struct {
	conn *Connection

	mutex             sync.RWMutex
	keepAlivePending  bool
	lastKeepAliveTime int64
	lastKeepAliveID   int64
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

func (player *Player) SetKeepAlivePending(keepAlivePending bool) {
	player.mutex.Lock()
	defer player.mutex.Unlock()
	player.keepAlivePending = keepAlivePending
}

func (player *Player) IsKeepAlivePending() bool {
	player.mutex.RLock()
	defer player.mutex.RUnlock()
	return player.keepAlivePending
}

func (player *Player) SetLastKeepAliveTime(lastKeepAliveTime int64) {
	player.mutex.Lock()
	defer player.mutex.Unlock()
	player.lastKeepAliveTime = lastKeepAliveTime
}

func (player *Player) GetLastKeepAliveTime() int64 {
	player.mutex.RLock()
	defer player.mutex.RUnlock()
	return player.lastKeepAliveTime
}

func (player *Player) SetLastKeepAliveID(lastKeepAliveID int64) {
	player.mutex.Lock()
	defer player.mutex.Unlock()
	player.lastKeepAliveID = lastKeepAliveID
}

func (player *Player) GetLastKeepAliveID() int64 {
	player.mutex.RLock()
	defer player.mutex.RUnlock()
	return player.lastKeepAliveID
}

func (player *Player) SendPacket(packet protocol.Packet) error {
	return player.conn.WritePacket(packet)
}

func (player *Player) Kick(reason string) error {
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

func NewPlayer(conn *Connection) *Player {
	return &Player{
		conn: conn,
	}
}
