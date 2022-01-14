package main

import (
	"context"
	"fmt"
	"time"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/kelseyhightower/envconfig"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/rs/zerolog/log"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"

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
	Env           string `envconfig:"ENVIRONMENT" default:"LOCAL"`
	Key           struct {
		Path string `envconfig:"PRIV_KEY_PATH"`
		GCP  GCPConfig
	}
}

type GCPConfig struct {
	ProjectID     string `envconfig:"PROJECT_ID"`
	SecretName    string `envconfig:"IPNS_KEY_SECRET_NAME"`
	SecretVersion string `envconfig:"IPNS_KEY_SECRET_VERSION"`
}

const (
	EnvLocal = "LOCAL"
	EnvGCP   = "GCP"
)

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
	ipnsManager, err := setupIPNSManager(config, shell)
	if err != nil {
		log.Fatal().Err(err).Msg("couldn't set up ipns manager")
	}
	instantCommiter := commiter.NewInstant(headManager, storage, memPinner, ipnsManager)
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
	privKey, pubKey, err := getIPNSKeyPair(config)
	if err != nil {
		return nil, err
	}
	return ipns.NewIPNSManager(privKey, pubKey, shell), nil
}

func getIPNSKeyPair(config Config) (crypto.PrivKey, crypto.PubKey, error) {
	switch config.Env {
	case EnvLocal:
		return ipns.ReadKeyPair(config.Key.Path)
	case EnvGCP:
		return getKeyPairFromSecretManager(config.Key.GCP)
	default:
		return nil, nil, fmt.Errorf("unsupported environment: %s", config.Env)
	}
}

func getKeyPairFromSecretManager(config GCPConfig) (privKey crypto.PrivKey, pubKey crypto.PubKey, err error) {
	// Create the client.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("create secret manager client: %w", err)
	}
	defer client.Close()

	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: config.SecretVersion,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		return nil, nil, fmt.Errorf("access secret version: %w", err)
	}

	return ipns.ParseKeyPair(result.Payload.Data)
}
