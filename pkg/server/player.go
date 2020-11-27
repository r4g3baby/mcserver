package server

import (
	"github.com/google/uuid"
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/protocol/packets"
)

type Player struct {
	conn *Connection

	keepAlivePending  bool
	lastKeepAliveTime int64
	lastKeepAliveID   int64
}

func (player *Player) GetServer() *Server {
	return player.conn.server
}

func (player *Player) GetUniqueID() uuid.UUID {
	return player.conn.uniqueID
}

func (player *Player) GetUsername() string {
	return player.conn.username
}

func (player *Player) GetProtocol() protocol.Protocol {
	return player.conn.protocol
}

func (player *Player) GetState() protocol.State {
	return player.conn.state
}

func (player *Player) SetKeepAlivePending(keepAlivePending bool) {
	player.keepAlivePending = keepAlivePending
}

func (player *Player) IsKeepAlivePending() bool {
	return player.keepAlivePending
}

func (player *Player) SetLastKeepAliveTime(lastKeepAliveTime int64) {
	player.lastKeepAliveTime = lastKeepAliveTime
}

func (player *Player) GetLastKeepAliveTime() int64 {
	return player.lastKeepAliveTime
}

func (player *Player) SetLastKeepAliveID(lastKeepAliveID int64) {
	player.lastKeepAliveID = lastKeepAliveID
}

func (player *Player) GetLastKeepAliveID() int64 {
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
