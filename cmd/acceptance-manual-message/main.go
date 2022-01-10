package main

import (
	"context"
	"time"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/areknoster/public-distributed-commit-log/grpc"
	"github.com/areknoster/public-distributed-commit-log/producer"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/signing"
	daemonstorage "github.com/areknoster/public-distributed-commit-log/storage/ipfs/daemon"
	"github.com/areknoster/public-distributed-commit-log/test/testpb"
)

type Config struct {
	SentinelConn grpc.ConnConfig
	SignerID     string `envconfig:"SIGNER_ID" required:"true"`
	PrivKeyPath  string `envconfig:"PRIV_KEY_PATH" required:"true"`
}

func main() {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatal().Err(err).Msg("load PDCL config")
	}

	writer := daemonstorage.NewStorage(shell.NewShell("localhost:5001"))

	signer, err := signing.ReadEd25519(cfg.PrivKeyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("get signer")
	}

	signedWriter := signing.NewSignedMessageWriter(writer, cfg.SignerID, signer)

	sentinelConn, err := grpc.Dial(cfg.SentinelConn)
	if err != nil {
		log.Fatal().Err(err).Msg("dial sentinel")
	}
	sentinelClient := sentinelpb.NewSentinelClient(sentinelConn)

	prod := producer.NewMessageProducer(signedWriter, sentinelClient)

	msg := &testpb.Message{
		IdIncremental: 333,
		Uuid:          "abc",
		Created:       timestamppb.New(time.Time{}),
	}
	log.Info().Interface("msg", msg).Msg("writing message")
	if err := prod.Produce(context.Background(), msg); err != nil {
		log.Fatal().Err(err).Msg("message production failed")
	}
}
