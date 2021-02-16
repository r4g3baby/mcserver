package server

var (
	OnPlayerChatEvent = "onPlayerChatEvent"
)

type PlayerChatEvent struct {
	player  *Player
	message string
}

func (e *PlayerChatEvent) GetPlayer() *Player {
	return e.player
}

func (e *PlayerChatEvent) SetMessage(message string) {
	e.message = message
}

func (e *PlayerChatEvent) GetMessage() string {
	return e.message
}

func NewPlayerChatEvent(player *Player, message string) *PlayerChatEvent {
	return &PlayerChatEvent{player: player, message: message}
}
