package main

import (
	"context"
	"flag"
	"flukis/invokiss/app/http/controller"
	"flukis/invokiss/database/querier"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {
	var configFileName string
	flag.StringVar(
		&configFileName,
		"c",
		"config.yml",
		"Configuration file name in *.yml format",
	)
	flag.Parse()

	cfg := loadConfig(configFileName)
	ctx := context.Background()

	pool, err := pgxpool.New(
		ctx,
		cfg.DBCfg.ConnStr(),
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to database")
	}

	writeProduct := querier.NewProductWriteModel(pool)
	readProduct := querier.NewProductReadModel(pool)
	writeCategory := querier.NewCategoryWriteModel(pool)
	readCategory := querier.NewCategoryReadModel(pool)

	productController := controller.NewProductController(
		writeProduct,
		readProduct,
	)

	categoryController := controller.NewCategoryController(
		writeCategory,
		readCategory,
	)

	r := chi.NewRouter()

	r.Mount("/api/product", productController.Routes())
	r.Mount("/api/category", categoryController.Routes())

	log.Info().Msg(fmt.Sprintf("starting up server on: %s", cfg.Listen.Addr()))
	server := &http.Server{
		Handler:      r,
		Addr:         cfg.Listen.Addr(),
		ReadTimeout:  time.Second * time.Duration(cfg.Listen.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(cfg.Listen.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(cfg.Listen.IdleTimeout),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("failed to start the server")
		return
	}
	log.Info().Msg("server stop")
}
