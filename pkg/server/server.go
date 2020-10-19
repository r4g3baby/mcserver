package server

import (
	"errors"
	"github.com/rs/zerolog/log"
	"net"
	"strconv"
	"strings"
)

var (
	ErrServerRunning = errors.New("server already running")
	ErrServerStopped = errors.New("server already stopped")
)

type Server struct {
	config Config

	listener net.Listener
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

	go func() {
		closed := false
		for !closed {
			client, err := listener.Accept()
			if err != nil {
				// See https://github.com/golang/go/issues/4373 for info.
				if !strings.Contains(err.Error(), "use of closed network connection") {
					log.Warn().Err(err).Msg("error occurred while accepting s new connection")
				} else {
					// End loop once the listener is closed.
					closed = true
				}
				continue
			}

			go server.handleClient(client)
		}
	}()

	return nil
}

func (server *Server) Stop() error {
	if server.listener == nil {
		return ErrServerStopped
	}

	log.Info().Msg("stopping server")

	if err := server.listener.Close(); err != nil {
		log.Error().Err(err).Msg("failed to close listener")
	}
	server.listener = nil

	return nil
}

func (server *Server) handleClient(conn net.Conn) {
	log.Debug().Stringer("connection", conn.RemoteAddr()).Msg("client connected")

	if err := conn.Close(); err != nil {
		// See https://github.com/golang/go/issues/4373 for info.
		if !strings.Contains(err.Error(), "use of closed network connection") {
			log.Warn().Err(err).Stringer("connection", conn.RemoteAddr()).Msg("got error while closing connection")
			return
		}
	}

	log.Debug().Stringer("connection", conn.RemoteAddr()).Msg("client disconnected")
}

func NewServer(config Config) *Server {
	return &Server{
		config: config,
	}
}
