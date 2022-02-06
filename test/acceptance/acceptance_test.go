package acceptance

import (
	"context"
	"crypto"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/areknoster/public-distributed-commit-log/consumer"
	pdclcrypto "github.com/areknoster/public-distributed-commit-log/crypto"
	"github.com/areknoster/public-distributed-commit-log/grpc"
	"github.com/areknoster/public-distributed-commit-log/ipns"
	"github.com/areknoster/public-distributed-commit-log/pdcl"
	"github.com/areknoster/public-distributed-commit-log/producer"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/storage"
	ipfsstorage "github.com/areknoster/public-distributed-commit-log/storage/message/ipfs"
	"github.com/areknoster/public-distributed-commit-log/storage/pbcodec"
	"github.com/areknoster/public-distributed-commit-log/test/testpb"
	memoryhead "github.com/areknoster/public-distributed-commit-log/thead/memory"
)

type Config struct {
	SentinelConn           grpc.ConnConfig
	Daemon                 ipfsstorage.Config
	SignerID               string `envconfig:"SIGNER_ID" required:"true"`
	ProducerPrivateKey     string `envconfig:"PRODUCER_KEY"`
	ProducerPrivateKeyPath string `envconfig:"PRODUCER_KEY_PATH"`
}

// TestAcceptance checks for acceptance requirements.
// due to the nature of PDCL, they are run with real test deployment
// which is deployed in the cloud.
// This gives us an idea, if the code really serves its purpose
func BenchmarkAcceptance(b *testing.B) {
	logsFile, err := os.OpenFile("bechmark.log", os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0o777)
	require.NoError(b, err)
	log.Logger = log.Output(logsFile).Level(zerolog.InfoLevel)

	globalCtx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	globalCtx, cancel = context.WithTimeout(globalCtx, time.Hour)
	defer cancel()

	deps := &dependencies{}
	deps.init(b, globalCtx)
	go func() {
		<-globalCtx.Done()
		close(deps.concurrentProducer.Messages())
		deps.globalWg.Wait()
	}()

	for _, messagesNumber := range []int{100, 500, 2500} {
		if b.Failed() {
			return
		}
		b.Run(fmt.Sprintf("Acceptance benchmark with %d messages", messagesNumber), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				produceConsumeN(globalCtx, b, deps, messagesNumber)
			}
		})
	}
}

func produceConsumeN(globalCtx context.Context, b *testing.B, deps *dependencies, messagesNumber int) {
	produceCtx, cancelProduce := context.WithCancel(globalCtx)
	defer cancelProduce()

	expectedToReceive := make(map[string]struct{})
	messages := deps.concurrentProducer.Messages()
	for j := 0; j < messagesNumber; j++ {
		if produceCtx.Err() != nil {
			return
		}
		messageUUID := uuid.NewString()
		expectedToReceive[messageUUID] = struct{}{}
		messages <- &testpb.Message{
			IdIncremental: int64(j),
			Uuid:          messageUUID,
			Created:       timestamppb.Now(),
		}
	}

	consumeCtx, cancelConsume := context.WithCancel(globalCtx)
	receivedUUIDsChan := make(chan string, 20)
	go func() {
		for gotUUID := range receivedUUIDsChan {
			delete(expectedToReceive, gotUUID)
			if len(expectedToReceive) == 0 {
				log.Info().Int("sent_received_messages", messagesNumber).Msg("received all messages")
				cancelConsume()
			}
		}
	}()

	err := deps.consumer.Consume(consumeCtx,
		consumer.MessageHandlerFunc(func(ctx context.Context, message storage.ProtoDecodable) error {
			tm := &testpb.Message{}
			if err := message.Decode(tm); err != nil {
				return err
			}
			log.Info().Int64("id_incremental", tm.IdIncremental).Str("uuid", tm.Uuid).Msg("message received")
			receivedUUIDsChan <- tm.Uuid
			return nil
		}))
	close(receivedUUIDsChan)
	require.ErrorIs(b, err, consumer.ErrContextDone)
}

type dependencies struct {
	globalWg sync.WaitGroup

	config              Config
	sh                  *shell.Shell
	codec               storage.Codec
	sentinelClient      sentinelpb.SentinelClient
	messageStorage      *ipfsstorage.DaemonStorage
	signedMessageWriter *pdclcrypto.SignedMessageWriter
	concurrentProducer  *producer.BasicConcurrentProducer
	consumer            *consumer.FirstToLastConsumer
	ipnsResolver        *ipns.IPNSResolver
	headIPNS            string
}

func (d *dependencies) init(tb testing.TB, globalCtx context.Context) {
	d.initConfig(tb)
	d.sh = ipfsstorage.NewShell(d.config.Daemon)
	d.codec = pbcodec.Json{}
	d.messageStorage = ipfsstorage.NewStorage(d.sh, d.codec)
	d.initSentinelClient(tb)
	d.initSignedMessageWriter(tb)
	d.initProducer(globalCtx)
	d.ipnsResolver = ipns.NewIPNSResolver(d.sh)
	d.initIPNSAddr(tb, globalCtx)
	d.initConsumerWithLatestHead(tb)

	d.globalWg.Add(1)
	go func() {
		for err := range d.concurrentProducer.Errors() {
			tb.Fatalf("got an error: %s", err.Err.Error())
		}
		d.globalWg.Done()
	}()
}

func (d *dependencies) initConsumerWithLatestHead(t testing.TB) {
	resolvedHead, err := d.ipnsResolver.ResolveIPNS(d.headIPNS)
	require.NoError(t, err)

	headCID, err := pdcl.ParseCID(strings.TrimPrefix(resolvedHead, "/ipfs/"))
	require.NoError(t, err)
	memoryTheadManager := memoryhead.NewHeadManager(headCID)
	signedMessageReadUnwrapper := pdclcrypto.NewSignedMessageUnwrapper(d.messageStorage, pbcodec.Json{})
	d.consumer = consumer.NewFirstToLastConsumer(
		memoryTheadManager,
		d.messageStorage,
		signedMessageReadUnwrapper,
		consumer.FirstToLastConsumerConfig{
			PollInterval: 20 * time.Second,
			PollTimeout:  20 * time.Second,
			IPNSAddr:     d.headIPNS,
		},
		d.ipnsResolver,
	)
}

func (d *dependencies) initIPNSAddr(t testing.TB, globalCtx context.Context) {
	headIPNSResp, err := d.sentinelClient.GetHeadIPNS(globalCtx, &sentinelpb.GetHeadIPNSRequest{})
	require.NoError(t, err)
	d.headIPNS = headIPNSResp.IpnsAddr
}

func (d *dependencies) initProducer(globalCtx context.Context) {
	blockingProducer := producer.NewBlockingProducer(d.signedMessageWriter, d.sentinelClient)
	d.concurrentProducer = producer.StartBasicConcurrentProducer(
		globalCtx,
		blockingProducer,
		producer.BasicConcurrentProducerConfig{
			JobsNumber:     200,
			ProduceTimeout: 3 * time.Minute,
			ErrBuf:         20,
			MessageBuf:     1000,
		},
	)
}

func (d *dependencies) initSignedMessageWriter(t testing.TB) {
	var key crypto.PrivateKey
	var err error

	switch {
	case d.config.ProducerPrivateKey != "":
		key, err = pdclcrypto.ParsePKCSKeyFromPEM([]byte(d.config.ProducerPrivateKey))
		require.NoError(t, err)
	case d.config.ProducerPrivateKeyPath != "":
		key, err = pdclcrypto.LoadFromPKCSFromPEMFile(d.config.ProducerPrivateKeyPath)
		require.NoError(t, err)
	default:
		t.Fatal("producer key is not set")
	}

	signer, ok := key.(crypto.Signer)
	require.True(t, ok)

	d.signedMessageWriter = pdclcrypto.NewSignedMessageWriter(
		d.messageStorage,
		pbcodec.Json{},
		d.config.SignerID,
		signer,
	)
}

func (d *dependencies) initSentinelClient(t testing.TB) {
	conn, err := grpc.Dial(d.config.SentinelConn)
	require.NoError(t, err)
	d.sentinelClient = sentinelpb.NewSentinelClient(conn)
}

func (d *dependencies) initConfig(t testing.TB) {
	cfg := Config{}
	require.NoError(t, envconfig.Process("", &cfg))
	d.config = cfg
}
