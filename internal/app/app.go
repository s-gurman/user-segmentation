package app

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/s-gurman/user-segmentation/config"
	httpapi "github.com/s-gurman/user-segmentation/internal/handler/http"
	segmentrepo "github.com/s-gurman/user-segmentation/internal/repository/segmentation/postgresql"
	segmentsvc "github.com/s-gurman/user-segmentation/internal/service/segmentation"
	httpserver "github.com/s-gurman/user-segmentation/pkg/http/server"
	"github.com/s-gurman/user-segmentation/pkg/logger"
	"github.com/s-gurman/user-segmentation/pkg/postgres"
)

const _defaultPGConnTimeout = 5 * time.Second

// @Title          User Segmentation Service API
// @Version        1.0
// @Description    This API provides dynamic user segmentation to conduct experiments.
// @Host           localhost:8081
// @BasePath       /api
func Run(cfg config.Config) {
	l := logger.New()
	defer l.Sync() // nolint:errcheck

	ctx, cancel := context.WithTimeout(context.Background(), _defaultPGConnTimeout)
	defer cancel()

	l.Infof("app Run - connecting to postresql at '%s'", cfg.PG.Address)
	pg, err := postgres.New(ctx, cfg.PG)
	if err != nil {
		l.Panicf("app Run - new postgres: %s", err)
	}
	defer pg.Close()

	repo := segmentrepo.NewPostgreSQL(pg)
	service := segmentsvc.New(repo.Segment, repo.Experiment)

	router := httpapi.NewRouter(service, l)
	server := httpserver.New(cfg.HTTP, router)

	ctx, cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	go func() {
		l.Infof("app Run - starting http server at 'localhost:%s'", cfg.HTTP.Port)

		if err = server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Errorf("app Run - server start err: %s", err)
		}
		cancel()
	}()

	<-ctx.Done()
	l.Info("app Run - gracefully shutdown ...")

	if err = server.Stop(); err != nil {
		l.Errorf("app Run - server stop err: %s", err)
	}
}
