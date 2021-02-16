package server

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/r4g3baby/mcserver/pkg/protocol/packets"
	"github.com/r4g3baby/mcserver/pkg/util/chat"
	"github.com/r4g3baby/mcserver/pkg/util/eventbus"
	"github.com/rs/zerolog/log"
	"math"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrServerRunning = errors.New("server already running")
	ErrServerStopped = errors.New("server already stopped")
)

type Server struct {
	config Config

	players  sync.Map
	eventbus eventbus.EventBus

	running  bool
	shutdown func()
}

func (server *Server) Start() error {
	if server.running {
		return ErrServerRunning
	}
	server.running = true

	bind := net.JoinHostPort(server.config.Host, strconv.Itoa(server.config.Port))
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		server.running = false
		return err
	}

	log.Info().Stringer("addr", listener.Addr()).Msg("server listening for new connections")

	var wait sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	server.shutdown = func() {
		cancel()
		if err := listener.Close(); err != nil {
			log.Error().Err(err).Msg("got error while closing listener")
		}
		wait.Wait()
	}

	wait.Add(2)
	go func() {
		defer wait.Done()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				client, err := listener.Accept()
				if err != nil {
					// See https://github.com/golang/go/issues/4373 for info.
					if !strings.Contains(err.Error(), "use of closed network connection") {
						log.Warn().Err(err).Msg("error occurred while accepting s new connection")
					}
					continue
				}

				go server.handleClient(client)
			}
		}
	}()

	go func() {
		defer wait.Done()

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				go server.sendKeepAlive()
			}
		}
	}()

	return nil
}

func (server *Server) Stop() error {
	if !server.running {
		return ErrServerStopped
	}

	log.Info().Msg("stopping server")

	server.shutdown()
	server.ForEachPlayer(func(player *Player) bool {
		_ = player.Kick([]chat.Component{
			&chat.TextComponent{
				Text: "Server is shutting down",
				BaseComponent: chat.BaseComponent{
					Color: &chat.Red,
				},
			},
		})
		return true
	})

	server.running = false

	return nil
}

func (server *Server) GetPlayerCount() int {
	var count int
	server.players.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

func (server *Server) GetPlayers() []*Player {
	var players []*Player
	server.players.Range(func(_, value interface{}) bool {
		players = append(players, value.(*Player))
		return true
	})
	return players
}

func (server *Server) GetPlayer(uniqueID uuid.UUID) *Player {
	if value, ok := server.players.Load(uniqueID); ok {
		return value.(*Player)
	}
	return nil
}

func (server *Server) ForEachPlayer(fn func(player *Player) bool) {
	server.players.Range(func(_, value interface{}) bool {
		return fn(value.(*Player))
	})
}

func (server *Server) FireEvent(event string, args ...interface{}) {
	server.eventbus.Publish(event, args...)
}

func (server *Server) On(event string, fn interface{}) error {
	return server.eventbus.Subscribe(event, fn)
}

func (server *Server) OnAsync(event string, fn interface{}) error {
	return server.eventbus.SubscribeAsync(event, fn)
}

func (server *Server) handleClient(conn net.Conn) {
	log.Debug().Stringer("connection", conn.RemoteAddr()).Msg("client connected")

	connection := NewConnection(conn, server)
	for {
		if err := connection.ReadPacket(); err != nil {
			// See https://github.com/golang/go/issues/4373 for info.
			if !strings.Contains(err.Error(), "use of closed network connection") {
				log.Error().Err(err).Stringer("connection", connection.RemoteAddr()).Msg("got error during packet read")
				// todo: should we disconnect?
			}
			break
		}
	}

	if err := connection.Close(); err != nil {
		// See https://github.com/golang/go/issues/4373 for info.
		if !strings.Contains(err.Error(), "use of closed network connection") {
			log.Warn().Err(err).Stringer("connection", connection.RemoteAddr()).Msg("got error while closing connection")
			return
		}
	}

	log.Debug().Stringer("connection", conn.RemoteAddr()).Msg("client disconnected")
}

func (server *Server) sendKeepAlive() {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	server.ForEachPlayer(func(player *Player) bool {
		if time.Since(player.GetLastKeepAliveTime()) >= 15*time.Second {
			if !player.IsKeepAlivePending() {
				if err := player.SendPacket(&packets.PacketPlayOutKeepAlive{
					KeepAliveID: random.Int31n(math.MaxInt32),
				}); err != nil {
					log.Warn().Err(err).Str("player", player.GetUsername()).Msg("failed to send keep alive packet")
				}
			} else {
				if err := player.Kick([]chat.Component{
					&chat.TextComponent{
						Text: "Timed out",
						BaseComponent: chat.BaseComponent{
							Color: &chat.Red,
						},
					},
				}); err != nil {
					log.Warn().Err(err).Str("player", player.GetUsername()).Msg("failed to kick player")
				}
			}
		}
		return true
	})
}

func (server *Server) addPlayer(conn *Connection) (*Player, bool) {
	value, loaded := server.players.LoadOrStore(conn.GetUniqueID(), newPlayer(conn))
	player := value.(*Player)
	if !loaded {
		log.Info().
			Str("name", player.GetUsername()).
			Stringer("uuid", player.GetUniqueID()).
			Msg("player joined the server")
	}
	return player, loaded
}

func (server *Server) removePlayer(uniqueID uuid.UUID) {
	if player, ok := server.players.LoadAndDelete(uniqueID); ok {
		player := player.(*Player)
		log.Info().
			Str("name", player.GetUsername()).
			Stringer("uuid", player.GetUniqueID()).
			Msg("player left the server")
	}
}

func NewServer(config Config) *Server {
	return &Server{
		config:   config,
		players:  sync.Map{},
		eventbus: eventbus.New(),
	}
}
