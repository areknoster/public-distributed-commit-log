package main

import (
	"context"
	"crypto"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/google/uuid"
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

func readEd25519(privKeyPath string) (crypto.Signer, error) {
	pemContent, err := os.ReadFile(privKeyPath)
	if err != nil {
		return nil, fmt.Errorf("read %s file content: %w", privKeyPath, err)
	}

	block, _ := pem.Decode(pemContent)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}
	priv, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key: %w", err)
	}
	key, ok := priv.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not ed25519.PrivateKey but %T", priv)
	}
	return key, nil
}

func main() {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		log.Fatal().Err(err).Msg("load PDCL config")
	}

	writer := daemonstorage.NewStorage(shell.NewShell("localhost:5001"))

	signer, err := readEd25519(cfg.PrivKeyPath)
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
		Uuid:          uuid.NewString(),
		Created:       timestamppb.Now(),
	}
	log.Info().Interface("msg", msg).Msg("writing message")
	if err := prod.Produce(context.Background(), msg); err != nil {
		log.Fatal().Err(err).Msg("message production failed")
	}
}
