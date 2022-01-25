package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ipfs/go-cid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/areknoster/public-distributed-commit-log/consumer"
	pdclcrypto "github.com/areknoster/public-distributed-commit-log/crypto"
	"github.com/areknoster/public-distributed-commit-log/grpc"
	"github.com/areknoster/public-distributed-commit-log/ipns"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/storage"
	ipfsstorage "github.com/areknoster/public-distributed-commit-log/storage/message/ipfs"
	"github.com/areknoster/public-distributed-commit-log/storage/pbcodec"
	"github.com/areknoster/public-distributed-commit-log/test/testpb"
	"github.com/areknoster/public-distributed-commit-log/thead/memory"
)

type Config struct {
	SentinelConn grpc.ConnConfig
}

func main() {
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal().Err(err).Msg("can't process environment variables for config")
	}
	setupPDCL(context.Background(), config)
}

func waitForShutdown() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	<-signalCh
	log.Debug().Msg("interruption signal received")
}

func setupPDCL(ctx context.Context, config Config) {
	conn, err := grpc.Dial(config.SentinelConn)
	if err != nil {
		log.Fatal().Err(err).Msg("can't connect to sentinel")
	}
	sentinelClient := sentinelpb.NewSentinelClient(conn)
	resp, err := sentinelClient.GetHeadIPNS(context.Background(), &sentinelpb.GetHeadIPNSRequest{})
	if err != nil {
		log.Fatal().Err(err).Msg("getting ipns address")
	}
	log.Info().Msgf("IPNS head address is %s", resp.IpnsAddr)
	if resp.IpnsAddr == "" {
		log.Fatal().Msg("could not get valid IPNS address from sentinel")
	}
	consumerOffsetManager := memory.NewHeadManager(cid.Undef)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize storage")
	}
	shell := shell.NewShell("localhost:5001")
	ipnsResolver := ipns.NewIPNSResolver(shell)
	reader := ipfsstorage.NewStorage(shell, pbcodec.Json{})

	firstToLastConsumer := consumer.NewFirstToLastConsumer(
		consumerOffsetManager,
		reader,
		pdclcrypto.NewSignedMessageUnwrapper(reader, pbcodec.Json{}),
		consumer.FirstToLastConsumerConfig{
			PollInterval: 20 * time.Second,
			PollTimeout:  20 * time.Second,
			IPNSAddr:     resp.IpnsAddr,
		},
		ipnsResolver,
	)

	err = firstToLastConsumer.Consume(ctx, consumer.MessageHandlerFunc(
		func(ctx context.Context, decodable storage.ProtoDecodable) error {
			pbCommit := &testpb.Message{}
			if err := decodable.Decode(pbCommit); err != nil {
				return fmt.Errorf("decode to commit proto: %w", err)
			}

			fmt.Println(protojson.Format(pbCommit))
			return nil
		}))
	waitForShutdown()
	log.Debug().Msg("server stopped consuming messages")
	if err != nil {
		log.Fatal().Err(err).Msg("consuming messages")
	}
}
