package main

import (
	"context"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/areknoster/public-distributed-commit-log/cmd/openpollution/pb"
	"github.com/areknoster/public-distributed-commit-log/producer"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/storage"
	"github.com/areknoster/public-distributed-commit-log/storage/localfs"
)

type Config struct {
	Host     string        `envconfig:"SENTINEL_SERVICE_HOST" default:"localhost"`
	Port     string        `envconfig:"SENTINEL_SERVICE_PORT" default:"8000"`
	Interval time.Duration `envconfig:"INTERVAL" default:"1s"`
}

func main() {
	config := Config{}
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal().Err(err).Msg("can't process environment variables for config")
	}

	// TODO: make this configurable
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("reading user home directory")
	}
	contentStorage, err := localfs.NewStorage(dirname + "/.local/share/pdcl/storage")
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize contentStorage")
	}
	messageStorage := storage.NewProtoMessageStorage(contentStorage)

	conn, err := grpc.Dial(
		net.JoinHostPort(config.Host, config.Port),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("can't connect to sentinel")
	}
	sentinelClient := sentinelpb.NewSentinelClient(conn)
	messageProducer := producer.NewMessageProducer(messageStorage, sentinelClient)
	r := randomOPMessageProducer{producer: messageProducer, interval: config.Interval}
	r.run()
}

type randomOPMessageProducer struct {
	producer producer.Producer
	interval time.Duration
}

func (r *randomOPMessageProducer) run() {
	for {
		time.Sleep(r.interval)
		message := &pb.Message{
			MeasureTime: timestamppb.Now(),
			Location: &pb.Location{
				Latitude:   rand.NormFloat64() * 90,
				Longtitude: rand.NormFloat64() * 180,
			},
			PollutionLevel: rand.NormFloat64() * 100,
		}
		log.Info().Time("measure_time", message.MeasureTime.AsTime()).Msg("produced message")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		if err := r.producer.Produce(ctx, message); err != nil {
			log.Fatal().Err(err).Msg("error producing message")
		}
		cancel()
	}
}
