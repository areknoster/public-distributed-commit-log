package main

import (
	"context"
	"crypto"

	"github.com/google/uuid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/areknoster/public-distributed-commit-log/crypto"
	"github.com/areknoster/public-distributed-commit-log/grpc"
	"github.com/areknoster/public-distributed-commit-log/producer"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	ipfsstorage "github.com/areknoster/public-distributed-commit-log/storage/message/ipfs"
	"github.com/areknoster/public-distributed-commit-log/storage/pbcodec"
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

	codec := pbcodec.Json{}

	writer := ipfsstorage.NewStorage(shell.NewShell("localhost:5001"), codec)

	privKey, err := pdclcrypto.LoadFromPKCSFromPEMFile(cfg.PrivKeyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("get privKey")
	}

	signer, ok := privKey.(crypto.Signer)
	if !ok {
		log.Fatal().Msgf("key is not private crypto.Signer type but %T", privKey)
	}

	signedWriter := pdclcrypto.NewSignedMessageWriter(writer, codec, cfg.SignerID, signer)

	sentinelConn, err := grpc.Dial(cfg.SentinelConn)
	if err != nil {
		log.Fatal().Err(err).Msg("dial sentinel")
	}
	sentinelClient := sentinelpb.NewSentinelClient(sentinelConn)

	prod := producer.NewMessageProducer(signedWriter, sentinelClient)

	var g errgroup.Group
	for i := 0; i < 10; i++ {
		id := i
		g.Go(func() error {
			msg := &testpb.Message{
				IdIncremental: int64(id),
				Uuid:          uuid.NewString(),
				Created:       timestamppb.Now(),
			}
			log.Info().Interface("msg", msg).Msg("writing message")
			return prod.Produce(context.Background(), msg)
		})
	}
	if err := g.Wait(); err != nil {
		log.Fatal().Err(err).Msg("messags production failed")
	}
}
