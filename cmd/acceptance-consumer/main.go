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

	"github.com/areknoster/public-distributed-commit-log/consumer"
	"github.com/areknoster/public-distributed-commit-log/grpc"
	"github.com/areknoster/public-distributed-commit-log/ipns"
	"github.com/areknoster/public-distributed-commit-log/pdclpb"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/signing"
	"github.com/areknoster/public-distributed-commit-log/storage"
	daemonstorage "github.com/areknoster/public-distributed-commit-log/storage/ipfs/daemon"
	"github.com/areknoster/public-distributed-commit-log/thead/memory"
	"github.com/areknoster/public-distributed-commit-log/thead/sentinelhead"
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
	sentinelHeadReader := sentinelhead.New(sentinelClient)
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
	reader := daemonstorage.NewStorage(shell)
	ipnsMgr := ipns.NewIPNSManager(nil, nil, shell)

	firstToLastConsumer := consumer.NewFirstToLastConsumer(
		sentinelHeadReader,
		consumerOffsetManager,
		&signing.SignedMessageUnwrapper{Base: reader},
		consumer.FirstToLastConsumerConfig{
			PollInterval: 10 * time.Second,
			PollTimeout:  100 * time.Second,
		},
		ipnsMgr,
		resp.IpnsAddr)

	err = firstToLastConsumer.Consume(ctx, consumer.MessageHandlerFunc(
		func(ctx context.Context, unmarshallable storage.ProtoUnmarshallable) error {
			pbCommit := &pdclpb.Commit{}
			if err := unmarshallable.Unmarshall(pbCommit); err != nil {
				return fmt.Errorf("unmarshall to commit proto: %w", err)
			}

			fmt.Printf("DEBUG: %+v", pbCommit)
			return nil
		}))
	waitForShutdown()
	log.Debug().Msg("server stopped consuming messages")
	if err != nil {
		log.Fatal().Err(err).Msg("consuming messages")
	}
}
