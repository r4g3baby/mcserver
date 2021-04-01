package cmd

import (
	"fmt"
	"github.com/r4g3baby/mcserver/internal"
	"github.com/r4g3baby/mcserver/pkg/protocol/packets"
	"github.com/r4g3baby/mcserver/pkg/server"
	"github.com/r4g3baby/mcserver/pkg/util/chat"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"os/signal"
	"syscall"
)

var (
	// set at build time
	Version = "1.0.0-dev"
)

var rootCmd = &cobra.Command{
	Use:     "mcserver",
	Short:   "A lightweight and performant minecraft server made in Go",
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		var config = setupConfig()
		setupLogger(config)

		serv := server.NewServer(config.Server)
		if err := serv.Start(); err != nil {
			log.Fatal().Err(err).Msg("failed to start server")
		}

		_ = serv.OnAsync(server.OnPacketReadEvent, func(e server.PacketEvent) {
			if chatPacket, ok := e.GetPacket().(*packets.PacketPlayInChatMessage); ok {
				serv.ForEachPlayer(func(player server.Player) bool {
					_ = player.SendPacket(&packets.PacketPlayOutChatMessage{
						Message: []chat.Component{
							&chat.TranslatableComponent{
								Translate: "chat.type.text",
								With: []chat.Component{
									&chat.TextComponent{
										Text: player.GetUsername(),
										BaseComponent: chat.BaseComponent{
											ClickEvent: &chat.ClickEvent{
												Action: chat.SuggestCommandClickAction,
												Value:  "/tell " + player.GetUsername(),
											},
											HoverEvent: &chat.HoverEvent{
												Action: chat.ShowEntityHoverAction,
												Contents: fmt.Sprintf(
													"{id:%s,type:minecraft:player,name:%s}",
													player.GetUniqueID(), player.GetUsername(),
												),
											},
											Insertion: player.GetUsername(),
										},
									},
									&chat.TextComponent{
										Text: chatPacket.Message,
									},
								},
							},
						},
						Position: 0,
						Sender:   player.GetUniqueID(),
					})
					return true
				})
			}
		})

		shutdownSignal := make(chan os.Signal, 1)
		signal.Notify(shutdownSignal, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP, os.Interrupt)
		sig := <-shutdownSignal

		log.Debug().Str("signal", sig.String()).Msg("received shutdown signal")
		if err := serv.Stop(); err != nil {
			log.Fatal().Err(err).Msg("failed to stop server")
		}
	},
}

func init() {
	rootCmd.Flags().Bool("debug", false, "sets log level to debug")
	rootCmd.Flags().String("host", "0.0.0.0", "sets the server host")
	rootCmd.Flags().Int("port", 25565, "sets the server port")

	_ = viper.BindPFlag("debug", rootCmd.Flags().Lookup("debug"))
	_ = viper.BindPFlag("server.host", rootCmd.Flags().Lookup("host"))
	_ = viper.BindPFlag("server.port", rootCmd.Flags().Lookup("port"))
}

func setupLogger(config internal.Config) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if config.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	consoleLogger := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "02 Jan 15:04"}
	if config.Logger.Enabled {
		fileLogger := consoleLogger
		fileLogger.Out = &lumberjack.Logger{
			Filename:   config.Logger.Filename,
			MaxSize:    config.Logger.MaxSize,
			MaxAge:     config.Logger.MaxAge,
			MaxBackups: config.Logger.MaxBackups,
			LocalTime:  config.Logger.LocalTime,
			Compress:   config.Logger.Compress,
		}
		fileLogger.NoColor = true
		log.Logger = log.Output(zerolog.MultiLevelWriter(consoleLogger, fileLogger))
	} else {
		log.Logger = log.Output(consoleLogger)
	}
}

func setupConfig() internal.Config {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("toml")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.SafeWriteConfig(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
				os.Exit(1)
			}
		} else {
			_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}
	}

	var config internal.Config
	if err := viper.Unmarshal(&config); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	return config
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
