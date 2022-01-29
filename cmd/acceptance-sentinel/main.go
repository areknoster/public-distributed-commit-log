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
	ipfsstorage "github.com/areknoster/public-distributed-commit-log/storage/message/ipfs"
	"github.com/areknoster/public-distributed-commit-log/storage/pbcodec"
	memoryhead "github.com/areknoster/public-distributed-commit-log/thead/memory"
)

type Config struct {
	DaemonStorage ipfsstorage.Config
	Validator     validator.Config
	GRPC          grpc.ServerConfig
	Env           string `envconfig:"ENVIRONMENT" default:"LOCAL"`
	Key           struct {
		Path string `envconfig:"PRIV_KEY_PATH"`
		GCP  GCPConfig
	}
	Commiter commiter.MaxBufferCommiterConfig
}

type GCPConfig struct {
	ProjectID     string `envconfig:"PROJECT_ID"`
	SecretName    string `envconfig:"IPNS_KEY_SECRET_NAME"`
	SecretVersion string `envconfig:"IPNS_KEY_SECRET_VERSION"`
}

func main() {
	log.Info().Msg("initializing sentinel")
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal().Err(err).Msg("can't process environment variables for config")
	}

	codec := pbcodec.Json{}

	ipfsShell := ipfsstorage.NewShell(config.DaemonStorage)
	storage := ipfsstorage.NewStorage(ipfsShell, codec)
	messageValidator, err := validator.New(storage, codec, config.Validator)
	if err != nil {
		log.Fatal().Err(err).Msg("initialize message validator")
	}
	memPinner := pinner.NewMemoryPinner()
	headManager := memoryhead.NewHeadManager(cid.Undef)
	ipnsManager, err := setupIPNSManager(config, ipfsShell)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't set up ipns manager")
	}
	instantCommiter := commiter.NewMaxBufferCommitter(
		headManager,
		storage,
		memPinner,
		ipnsManager,
		config.Commiter)
	sentinel := service.New(messageValidator, memPinner, instantCommiter, headManager, ipnsManager)

	rateLimiter := ratelimiting.NewAlwaysAllowLimiter()
	grpcserver, err := grpc.NewServer(config.GRPC, rateLimiter)
	if err != nil {
		log.Panic().Err(err).Msg("initialize sentinel GRPC server")
	}
	sentinelpb.RegisterSentinelServer(grpcserver, sentinel)
	log.Fatal().Err(grpcserver.ListenAndServe()).Msg("error running grpc server")
}

func setupIPNSManager(config Config, shell *shell.Shell) (ipns.Manager, error) {
	return ipns.NewIPNSManager(shell)
}
