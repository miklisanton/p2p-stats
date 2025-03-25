package main

import (
	"os"
	"p2p-stats/internal/config"
	"p2p-stats/internal/db/drivers"
	"p2p-stats/internal/db/repository"
	"p2p-stats/internal/telegram"
	"p2p-stats/internal/utils"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	cfg *config.Config
	db  *sqlx.DB
)

func init() {
	// Set up zerolog with the custom writer
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC1123Z}).
		Level(zerolog.DebugLevel).
		With().
		Timestamp().
		Caller().
		Logger().
		Hook(utils.GoroutineHook{})
	log.Info().Msg("Logger initialized")
	cfgPath, err := config.ParseCLI()
	// Parse CLI arguments
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to parse CLI")
	}
	// Read config
	cfg, err = config.NewConfig(cfgPath)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read config")
	}
	log.Info().Msg("Config loaded")
	// Connect to database
	db, err = drivers.Connect(cfg.Db.Name, "internal/db/migrations")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	log.Info().Msg("Database connected. Path: " + cfg.Db.Name)
}

func main() {
    from, to, err := utils.ParseList("11-11-2022 11-11-2023")
    if err != nil {
        log.Fatal().Err(err).Msg("Failed to parse list")
    }
    log.Info().Msgf("From: %v, To: %v", from, to)
    defer db.Close()
    // Create repositories
    userRepo := repository.NewUserRepository(db)
    recordRepo := repository.NewRecordRepository(db)
    client := telegram.NewClient(cfg.Telegram.Token, recordRepo, userRepo)
    // Start telegram client
    client.Start()
}

