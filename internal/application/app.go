package application

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	_ "github.com/jackc/pgx/v5/stdlib" // postgres driver
	"github.com/jmoiron/sqlx"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/Magic-Kot/Store-notification-service/internal/config"
	"github.com/Magic-Kot/Store-notification-service/internal/domain/subjects"
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
}

func New(name, version string, cfg config.Config) *App {
	return &App{ //nolint:exhaustruct
		name:    name,
		version: version,
		cfg:     cfg,
	}
}

func (app *App) Run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	defer stop()

	logger, err := logging.NewLogger(&logging.LoggerDeps{Level: app.cfg.Logger.Level})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init logger")
	}

	ctx = logger.WithContext(ctx)

	g, ctx := errgroup.WithContext(ctx)

	natsClient, err := n.NewClient(ctx, &n.Client{Url: app.cfg.Nats.URL})
	if err != nil {
		return fmt.Errorf("NewClient: %w", err)
	}

	app.natsClient = natsClient

	err = subjects.SubscribeNotification(app.natsClient)
	if err != nil {
		return fmt.Errorf("subjects.SubscribeNotification: %w", err)
	}

	app.runServer(ctx, g)

	if err = g.Wait(); err != nil {
		return fmt.Errorf("g.Wait: %w", err)
	}

	zerolog.Ctx(ctx).Info().Msg("server stopped")

	return nil
}

func (app *App) runServer(ctx context.Context, g *errgroup.Group) {
	g.Go(func() error {
		go func() {
			<-ctx.Done()

			ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), app.cfg.Server.ShutdownTimeout) //nolint:govet
			defer cancel()

			if err := app.natsClient.Drain(); err != nil {
				zerolog.Ctx(ctx).Error().Err(err).Msg("drain failed")
			}
			fmt.Printf("\n\ndrain success\n\n")
		}()

		zerolog.Ctx(ctx).Info().Msg("server started")

		<-ctx.Done()

		return nil
	})
}
