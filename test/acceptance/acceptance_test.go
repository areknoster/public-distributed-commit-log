package acceptance

import (
	"testing"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/require"

	"github.com/areknoster/public-distributed-commit-log/grpc"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/storage"
	daemonstorage "github.com/areknoster/public-distributed-commit-log/storage/ipfs/daemon"
)

type Config struct {
	sentinelConn grpc.ConnConfig
	daemon       daemonstorage.Config
}

// TestAcceptance checks for acceptance requirements.
// due to the nature of PDCL, they are run with real test deployment
// which is separate from this test logic.
func TestAcceptance(t *testing.T) {
	config := initConfig(t)
	messageStorage := initStorage(t, config)
	sentinelClient := initSentinelClient(t, config)
	// signedMessageWriter := signing.NewSignedMessageWriter(storage, )
	//
	// prod := producer.NewMessageProducer(, sentinelClient)
	//
	// t.Run("Producer should be able to add at least 100 messages per second", func(t *testing.T) {
	//
	// })
	_, _ = messageStorage, sentinelClient
}

func initSentinelClient(t *testing.T, config Config) sentinelpb.SentinelClient {
	conn, err := grpc.Dial(config.sentinelConn)
	require.NoError(t, err)
	return sentinelpb.NewSentinelClient(conn)
}

func initStorage(t *testing.T, config Config) storage.MessageStorage {
	return daemonstorage.NewStorage(daemonstorage.NewShell(config.daemon))
}

func initConfig(t *testing.T) Config {
	cfg := Config{}
	require.NoError(t, envconfig.Process("", &cfg))
	return cfg
}
