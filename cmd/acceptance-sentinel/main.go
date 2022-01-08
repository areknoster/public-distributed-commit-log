package main

import (
	"github.com/ipfs/go-cid"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"

	"github.com/areknoster/public-distributed-commit-log/cmd/acceptance-sentinel/internal/validator"
	"github.com/areknoster/public-distributed-commit-log/grpc"
	commiter "github.com/areknoster/public-distributed-commit-log/sentinel/commiter"
	"github.com/areknoster/public-distributed-commit-log/sentinel/pinner"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/sentinel/service"
	daemonstorage "github.com/areknoster/public-distributed-commit-log/storage/ipfs/daemon"
	memoryhead "github.com/areknoster/public-distributed-commit-log/thead/memory"
)

type Config struct {
	DaemonStorage daemonstorage.Config
	Validator     validator.Config
	GRPC          grpc.ServerConfig
}

func main() {
	log.Info().Msg("initializing sentinel")
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal().Err(err).Msg("can't process environment variables for config")
	}

	storage := daemonstorage.NewStorage(daemonstorage.NewShell(config.DaemonStorage))
	messageValidator, err := validator.New(storage, config.Validator)
	if err != nil {
		log.Panic().Err(err).Msg("initialize message validator")
	}
	memPinner := pinner.NewMemoryPinner()
	headManager := memoryhead.NewHeadManager(cid.Undef)
	instantCommiter := commiter.NewInstant(headManager, storage, memPinner)
	sentinel := service.New(messageValidator, memPinner, instantCommiter, headManager)

	grpcserver, err := grpc.NewServer(config.GRPC)
	if err != nil {
		log.Panic().Err(err).Msg("initialize sentinel GRPC server")
	}
	sentinelpb.RegisterSentinelServer(grpcserver, sentinel)
	log.Fatal().Err(grpcserver.ListenAndServe()).Msg("error running grpc server")
}
