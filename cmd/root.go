package cmd

import (
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/r4g3baby/mcserver/internal/config"
	"github.com/r4g3baby/mcserver/pkg/log"
	"github.com/r4g3baby/mcserver/pkg/protocol/packets"
	"github.com/r4g3baby/mcserver/pkg/server"
	"github.com/r4g3baby/mcserver/pkg/util/chat"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
			log.Log.Error(err, "failed to start server")
			os.Exit(1)
		}

		world := serv.GetWorld()
		world.SetBlock(0, 65, 0, "minecraft:torch")
		world.SetBlock(0, 64, 0, "minecraft:dirt")
		world.SetBlock(1, 64, 0, "minecraft:stone")
		world.SetBlock(1, 64, 1, "minecraft:stone")
		world.SetBlock(0, 64, 1, "minecraft:stone")
		world.SetBlock(-1, 64, 1, "minecraft:stone")
		world.SetBlock(-1, 64, 0, "minecraft:stone")
		world.SetBlock(-1, 64, -1, "minecraft:stone")
		world.SetBlock(0, 64, -1, "minecraft:stone")
		world.SetBlock(1, 64, -1, "minecraft:stone")

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

		log.Log.V(1).Info("received shutdown signal", "signal", sig)
		if err := serv.Stop(); err != nil {
			log.Log.Error(err, "failed to stop server")
			os.Exit(1)
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

func setupConfig() config.Config {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := config.WriteDefaultConfig(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			return setupConfig()
		} else {
			_, _ = fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
	}

	var conf config.Config
	if err := viper.Unmarshal(&conf); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	return conf
}

func setupLogger(config config.Config) {
	var zapConfig = zap.NewProductionEncoderConfig()
	var level = zap.InfoLevel
	if config.Debug {
		level = zap.DebugLevel
	}

	zapConfig.ConsoleSeparator = "\u0020"
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("02 Jan 15:04")
	zapConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(zapConfig)

	var core zapcore.Core
	if config.Logger.Enabled {
		zapConfig.EncodeTime = zapcore.EpochNanosTimeEncoder
		zapConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		fileEncoder := zapcore.NewJSONEncoder(zapConfig)

		core = zapcore.NewTee(
			zapcore.NewCore(fileEncoder, zapcore.AddSync(&lumberjack.Logger{
				Filename:   config.Logger.Filename,
				MaxSize:    config.Logger.MaxSize,
				MaxAge:     config.Logger.MaxAge,
				MaxBackups: config.Logger.MaxBackups,
				LocalTime:  config.Logger.LocalTime,
				Compress:   config.Logger.Compress,
			}), level),
			zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stderr), level),
		)
	} else {
		core = zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stderr), level)
	}

	log.SetLogger(zapr.NewLogger(zap.New(core, zap.WithCaller(false))))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
