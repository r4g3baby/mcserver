package server

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/r4g3baby/mcserver/pkg/log"
	"github.com/r4g3baby/mcserver/pkg/protocol"
	"github.com/r4g3baby/mcserver/pkg/protocol/packets"
	"github.com/r4g3baby/mcserver/pkg/util/bytes"
	"github.com/r4g3baby/mcserver/pkg/util/chat"
	"github.com/r4g3baby/mcserver/pkg/util/eventbus"
	"github.com/r4g3baby/mcserver/pkg/util/schematic"
	"io"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

var (
	ErrServerRunning = errors.New("server already running")
	ErrServerStopped = errors.New("server already stopped")
)

type (
	Server interface {
		Start() error
		Stop() error

		GetConfig() Config
		GetWorld() World

		GetPlayerCount() int
		GetPlayers() []Player
		GetPlayer(uniqueID uuid.UUID) Player
		ForEachPlayer(fn func(player Player) bool)

		FireEvent(event string, args ...interface{})
		On(event string, fn interface{}, priority ...eventbus.Priority) error
		OnAsync(event string, fn interface{}) error

		createPlayer(conn Connection) (player Player, online bool)
		removePlayer(uniqueID uuid.UUID)
	}

	server struct {
		config Config
		world  World

		players  sync.Map
		eventbus eventbus.EventBus

		running  bool
		shutdown func()
	}
)

func (server *server) Start() error {
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

	log.Log.WithValues(
		"addr", listener.Addr(),
	).Info("server listening for new connections")

	var wait sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	server.shutdown = func() {
		cancel()
		if err := listener.Close(); err != nil {
			log.Log.Error(err, "got error while closing listener")
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
					if !errors.Is(err, net.ErrClosed) {
						log.Log.Error(err, "error occurred while accepting a new connection")
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

func (server *server) Stop() error {
	if !server.running {
		return ErrServerStopped
	}

	log.Log.Info("stopping server")

	server.shutdown()
	server.ForEachPlayer(func(player Player) bool {
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

func (server *server) GetConfig() Config {
	return server.config
}

func (server *server) GetWorld() World {
	return server.world
}

func (server *server) GetPlayerCount() int {
	var count int
	server.players.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

func (server *server) GetPlayers() []Player {
	var players []Player
	server.players.Range(func(_, value interface{}) bool {
		players = append(players, value.(Player))
		return true
	})
	return players
}

func (server *server) GetPlayer(uniqueID uuid.UUID) Player {
	if value, ok := server.players.Load(uniqueID); ok {
		return value.(Player)
	}
	return nil
}

func (server *server) ForEachPlayer(fn func(player Player) bool) {
	server.players.Range(func(_, value interface{}) bool {
		return fn(value.(Player))
	})
}

func (server *server) FireEvent(event string, args ...interface{}) {
	server.eventbus.Publish(event, args...)
}

func (server *server) On(event string, fn interface{}, priority ...eventbus.Priority) error {
	return server.eventbus.Subscribe(event, fn, priority...)
}

func (server *server) OnAsync(event string, fn interface{}) error {
	return server.eventbus.SubscribeAsync(event, fn)
}

func (server *server) createPlayer(conn Connection) (Player, bool) {
	value, loaded := server.players.LoadOrStore(conn.GetUniqueID(), newPlayer(conn))
	player := value.(Player)
	if !loaded {
		log.Log.WithValues(
			"name", player.GetUsername(),
			"uuid", player.GetUniqueID(),
		).Info("player joined the server")
	}
	return player, loaded
}

func (server *server) removePlayer(uniqueID uuid.UUID) {
	if player, ok := server.players.LoadAndDelete(uniqueID); ok {
		player := player.(Player)
		log.Log.WithValues(
			"name", player.GetUsername(),
			"uuid", player.GetUniqueID(),
		).Info("player left the server")
	}
}

func (server *server) handleClient(conn net.Conn) {
	log.Log.WithValues(
		"connection", conn.RemoteAddr(),
	).V(1).Info("client connected")

	connection := newConnection(conn, server)
	for {
		if err := connection.ReadPacket(); err != nil {
			if !errors.Is(err, net.ErrClosed) && !errors.Is(err, io.EOF) {
				log.Log.WithValues(
					"connection", conn.RemoteAddr(),
				).Error(err, "got error during packet read")
				// todo: should we disconnect?
			}
			break
		}
	}

	if err := connection.Close(); err != nil {
		if !errors.Is(err, net.ErrClosed) {
			log.Log.WithValues(
				"connection", conn.RemoteAddr(),
			).Error(err, "got error while closing connection")
			return
		}
	}

	log.Log.WithValues(
		"connection", conn.RemoteAddr(),
	).V(1).Info("client disconnected")
}

func (server *server) sendKeepAlive() {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	server.ForEachPlayer(func(player Player) bool {
		if time.Since(player.GetLastKeepAliveTime()) >= 15*time.Second {
			if !player.IsKeepAlivePending() {
				if err := player.SendPacket(&packets.PacketPlayOutKeepAlive{
					KeepAliveID: random.Int31n(math.MaxInt32),
				}); err != nil {
					log.Log.WithValues(
						"name", player.GetUsername(),
						"uuid", player.GetUniqueID(),
					).Error(err, "failed to send keep alive packet")
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
					log.Log.WithValues(
						"name", player.GetUsername(),
						"uuid", player.GetUniqueID(),
					).Error(err, "failed to kick player")
				}
			}
		}
		return true
	})
}

func NewServer(config Config) Server {
	var world = NewWorld("overworld", protocol.Overworld)
	renderDistance := config.World.RenderDistance + 1
	for x := -renderDistance; x <= renderDistance; x++ {
		for z := -renderDistance; z <= renderDistance; z++ {
			world.GetChunk(x, z)
		}
	}

	if schemFileName := config.World.Schematic; schemFileName != "" {
		if fileBytes, err := os.ReadFile(schemFileName); err == nil {
			if schem, err := schematic.Read(bytes.NewBuffer(fileBytes)); err == nil {
				log.Log.WithValues(
					"width", schem.GetWidth(),
					"height", schem.GetHeight(),
					"length", schem.GetLength(),
					"size", schem.GetWidth()*schem.GetHeight()*schem.GetLength(),
				).Info("loading schematic")
				for x := 0; x < schem.GetWidth(); x++ {
					for y := 0; y < schem.GetHeight(); y++ {
						for z := 0; z < schem.GetLength(); z++ {
							world.SetBlock(x, y, z, schem.GetBlocks()[x][y][z])
						}
					}
				}
				log.Log.Info("done loading schematic")
			} else {
				log.Log.Error(err, "failed to read schematic file")
			}
		} else {
			log.Log.Error(err, "failed to read schematic file")
		}
	}

	return &server{
		config:   config,
		world:    world,
		players:  sync.Map{},
		eventbus: eventbus.New(),
	}
}
