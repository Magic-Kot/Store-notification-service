package application

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // postgres driver
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"

	"github.com/Magic-Kot/Store-notification-service/internal/config"
	"github.com/Magic-Kot/Store-notification-service/pkg/logging"
	n "github.com/Magic-Kot/Store-notification-service/pkg/nats"
)

type App struct {
	name     string
	version  string
	cfg      config.Config
	deferred []func()

	postgresClient *sqlx.DB
	natsClient     *nats.Conn

	dbUser      persistence.DBUser
	userService *service.User
}

func New(name, version string, cfg config.Config) *App {
	return &App{ //nolint:exhaustruct
		name:    name,
		version: version,
		cfg:     cfg,
	}
}

func (app *App) Run() error {
	ctx := context.Background()

	logger, err := logging.NewLogger(&logging.LoggerDeps{Level: app.cfg.Logger.Level})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init logger")
	}

	ctx = logger.WithContext(ctx)

	natsClient, err := n.NewClient(ctx, &n.Client{Url: app.cfg.Nats.URL})
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}

	app.natsClient = natsClient

	// notification
	//

	return nil
}
