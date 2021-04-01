package server

import (
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"sync"
)

var (
	OnPacketReadEvent  = "onPacketRead"
	OnPacketWriteEvent = "onPacketWrite"
)

type (
	PacketEvent interface {
		GetConnection() Connection
		GetPlayer() Player
		GetPacket() protocol.Packet
		SetCancelled(cancelled bool)
		IsCancelled() bool
	}

	packetEvent struct {
		connection Connection
		player     Player
		packet     protocol.Packet

		mutex     sync.RWMutex
		cancelled bool
	}
)

func (e *packetEvent) GetConnection() Connection {
	return e.connection
}

func (e *packetEvent) GetPlayer() Player {
	return e.player
}

func (e *packetEvent) GetPacket() protocol.Packet {
	return e.packet
}

func (e *packetEvent) SetCancelled(cancelled bool) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.cancelled = cancelled
}

func (e *packetEvent) IsCancelled() bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	return e.cancelled
}

func NewPacketEvent(connection Connection, player Player, packet protocol.Packet) PacketEvent {
	return &packetEvent{
		connection: connection,
		player:     player,
		packet:     packet,
		cancelled:  false,
	}
}
