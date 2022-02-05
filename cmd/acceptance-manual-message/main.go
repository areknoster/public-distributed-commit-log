package main

import (
	"context"
	"crypto"
	"time"

	"github.com/google/uuid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	pdclcrypto "github.com/areknoster/public-distributed-commit-log/crypto"
	"github.com/areknoster/public-distributed-commit-log/grpc"
	"github.com/areknoster/public-distributed-commit-log/producer"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	ipfsstorage "github.com/areknoster/public-distributed-commit-log/storage/message/ipfs"
	"github.com/areknoster/public-distributed-commit-log/storage/pbcodec"
	"github.com/areknoster/public-distributed-commit-log/test/testpb"
)

type Config struct {
	SentinelConn             grpc.ConnConfig
	ConcurrentProducerConfig producer.BasicConcurrentProducerConfig
	SignerID                 string `envconfig:"SIGNER_ID" required:"true"`
	PrivKeyPath              string `envconfig:"PRODUCER_KEY_PATH" required:"true"`
}

func main() {
	config := &Config{}
	if err := envconfig.Process("", config); err != nil {
		log.Fatal().Err(err).Msg("load PDCL config")
	}

	codec := pbcodec.Json{}

	writer := ipfsstorage.NewStorage(shell.NewShell("localhost:5001"), codec)

	privKey, err := pdclcrypto.LoadFromPKCSFromPEMFile(config.PrivKeyPath)
	if err != nil {
		log.Fatal().Err(err).Msg("get privKey")
	}

	signer, ok := privKey.(crypto.Signer)
	if !ok {
		log.Fatal().Msgf("key is not private crypto.Signer type but %T", privKey)
	}

	signedWriter := pdclcrypto.NewSignedMessageWriter(writer, codec, config.SignerID, signer)

	sentinelConn, err := grpc.Dial(config.SentinelConn)
	if err != nil {
		log.Fatal().Err(err).Msg("dial sentinel")
	}
	sentinelClient := sentinelpb.NewSentinelClient(sentinelConn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	blockingProducer := producer.NewBlockingProducer(signedWriter, sentinelClient)
	concurrentProducer := producer.StartBasicConcurrentProducer(ctx, blockingProducer, config.ConcurrentProducerConfig)

	go queueMessages(ctx, concurrentProducer.Messages())
	handleErrors(concurrentProducer.Errors())
}

func handleErrors(errors <-chan producer.Error) {
	for err := range errors {
		log.Error().
			RawJSON("message", []byte(protojson.Format(err.Message))).
			Err(err.Err).
			Msg("message production failed")
	}
}

func queueMessages(ctx context.Context, messages chan<- proto.Message) {
	defer close(messages)
	for i := int64(0); i < 100; i++ {
		select {
		case <-ctx.Done():
			log.Error().Int64("message_index", i).Msg("context done, production stopped before all messages were sent")
			return
		case messages <- &testpb.Message{
			IdIncremental: i,
			Uuid:          uuid.NewString(),
			Created:       timestamppb.Now(),
		}: // pass
		}
	}
	log.Info().Msg("messages queued")
}
