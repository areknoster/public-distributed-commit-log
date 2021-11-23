package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/areknoster/public-distributed-commit-log/cmd/openpollution/pb"
	"github.com/areknoster/public-distributed-commit-log/consumer"
	"github.com/areknoster/public-distributed-commit-log/head/memory"
	"github.com/areknoster/public-distributed-commit-log/head/sentinel_reader"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/storage"
	"github.com/areknoster/public-distributed-commit-log/storage/localfs"
)

type Config struct {
	Host string `envconfig:"SENTINEL_SERVICE_HOST" default:"localhost"`
	Port string `envconfig:"SENTINEL_SERVICE_PORT" default:"8000"`
}

func main() {
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal().Err(err).Msg("can't process environment variables for config")
	}

	conn, err := grpc.Dial(
		net.JoinHostPort(config.Host, config.Port),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("can't connect to sentinel")
	}
	sentinelClient := sentinelpb.NewSentinelClient(conn)
	sentinelHeadReader := sentinel_reader.NewSentinelHeadReader(sentinelClient)
	consumerOffsetManager := memory.NewHeadManager(cid.Undef)
	fsStorage, err := localfs.NewStorage("./storage")
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize storage")
	}
	messageStorage := storage.NewProtoMessageStorage(fsStorage)

	firstToLastConsumer := consumer.NewFirstToLastConsumer(
		sentinelHeadReader,
		consumerOffsetManager,
		messageStorage,
		consumer.FirstToLastConsumerConfig{
			PollInterval: 10 * time.Second,
			PollTimeout:  100 * time.Second,
		})

	c := make(chan os.Signal, 1)
	globalCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			cancel()
		}
	}()
	err = firstToLastConsumer.Consume(globalCtx, consumer.MessageFandlerFunc(
		func(ctx context.Context, unmarshallable storage.ProtoUnmarshallable) error {
			message := &pb.Message{}
			if err := unmarshallable.Unmarshall(message); err != nil {
				return fmt.Errorf("unmarshall message: %w", err)
			}
			jsonMessage, err := protojson.Marshal(message)
			if err != nil {
				return fmt.Errorf("marshall message: %w", err)
			}
			log.Info().RawJSON("message", jsonMessage).Msg("message received")
			return nil
		}))
	log.Fatal().Err(err).Msg("consume failed")
}
