package main

import (
	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"

	"github.com/areknoster/public-distributed-commit-log/cmd/acceptance-sentinel/internal/validator"
	"github.com/areknoster/public-distributed-commit-log/grpc"
	"github.com/areknoster/public-distributed-commit-log/ipns"
	"github.com/areknoster/public-distributed-commit-log/ratelimiting"
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
	PrivKeyPath   string `envconfig:"PRIV_KEY_PATH"`
}

func main() {
	log.Info().Msg("initializing sentinel")
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal().Err(err).Msg("can't process environment variables for config")
	}

	shell := daemonstorage.NewShell(config.DaemonStorage)
	storage := daemonstorage.NewStorage(shell)
	messageValidator, err := validator.New(storage, config.Validator)
	if err != nil {
		log.Fatal().Err(err).Msg("initialize message validator")
	}
	memPinner := pinner.NewMemoryPinner()
	headManager := memoryhead.NewHeadManager(cid.Undef)
	ipnsManager, err := setupIPNSManager(config.PrivKeyPath, shell)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't set up ipns manager")
	}
	instantCommiter := commiter.NewInstant(headManager, storage, memPinner, ipnsManager)
	sentinel := service.New(messageValidator, memPinner, instantCommiter, headManager)

	rateLimiter := ratelimiting.NewAlwaysAllowLimiter()
	grpcserver, err := grpc.NewServer(config.GRPC, rateLimiter)
	if err != nil {
		log.Panic().Err(err).Msg("initialize sentinel GRPC server")
	}
	sentinelpb.RegisterSentinelServer(grpcserver, sentinel)
	log.Fatal().Err(grpcserver.ListenAndServe()).Msg("error running grpc server")
}

func setupIPNSManager(keyPath string, shell *shell.Shell) (ipns.Manager, error) {
	if keyPath == "" {
		return ipns.NewNopManager(), nil
	}
	privKey, pubKey, err := ipns.ReadECKeys(keyPath)
	if err != nil {
		return nil, err
	}
	return ipns.NewIPNSManager(privKey, pubKey, shell), nil
}
