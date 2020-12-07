package server

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/r4g3baby/mcserver/pkg/protocol/packets"
	"github.com/r4g3baby/mcserver/pkg/util/chat"
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

	listener net.Listener
	wait     sync.WaitGroup
	shutdown context.CancelFunc
	players  sync.Map
}

func (server *Server) Start() error {
	if server.listener != nil {
		return ErrServerRunning
	}

	bind := net.JoinHostPort(server.config.Host, strconv.Itoa(server.config.Port))
	listener, err := net.Listen("tcp", bind)
	if err != nil {
		return err
	}
	server.listener = listener

	log.Info().Stringer("addr", listener.Addr()).Msg("server listening for new connections")

	ctx, cancel := context.WithCancel(context.Background())
	server.shutdown = cancel

	server.wait.Add(1)
	go func() {
		defer server.wait.Done()

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

	server.wait.Add(1)
	go func() {
		defer server.wait.Done()

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
	if server.listener == nil {
		return ErrServerStopped
	}

	log.Info().Msg("stopping server")
	server.ForEachPlayer(func(player *Player) {
		_ = player.Kick([]chat.Component{
			&chat.TextComponent{
				Text: "Server is shutting down",
				BaseComponent: chat.BaseComponent{
					Color: &chat.Red,
				},
			},
		})
	})

	server.shutdown()

	if err := server.listener.Close(); err != nil {
		log.Error().Err(err).Msg("got error while closing listener")
	}

	server.wait.Wait()
	server.listener = nil

	return nil
}

func (server *Server) GetPlayer(uniqueID uuid.UUID) *Player {
	if value, ok := server.players.Load(uniqueID); ok {
		return value.(*Player)
	}
	return nil
}

func (server *Server) ForEachPlayer(fn func(player *Player)) {
	server.players.Range(func(key, value interface{}) bool {
		fn(value.(*Player))
		return true
	})
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
	server.ForEachPlayer(func(player *Player) {
		currentTime := time.Now().UnixNano()
		if currentTime-player.GetLastKeepAliveTime() >= 15*int64(time.Second) {
			if !player.IsKeepAlivePending() {
				rand.Seed(currentTime)
				if err := player.SendPacket(&packets.PacketPlayOutKeepAlive{
					KeepAliveID: rand.Int63n(math.MaxInt32),
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
	})
}

func (server *Server) createPlayer(conn *Connection) *Player {
	player := NewPlayer(conn)
	server.players.Store(player.GetUniqueID(), player)
	return player
}

func (server *Server) removePlayer(uniqueID uuid.UUID) {
	server.players.Delete(uniqueID)
}

func NewServer(config Config) *Server {
	return &Server{
		config:  config,
		wait:    sync.WaitGroup{},
		players: sync.Map{},
	}
}
