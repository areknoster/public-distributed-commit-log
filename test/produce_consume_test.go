package test

import (
	"context"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"golang.org/x/sync/errgroup"

	consumer "github.com/areknoster/public-distributed-commit-log/consumer"
	"github.com/areknoster/public-distributed-commit-log/grpc"
	"github.com/areknoster/public-distributed-commit-log/head/memory"
	"github.com/areknoster/public-distributed-commit-log/head/sentinel_reader"
	"github.com/areknoster/public-distributed-commit-log/producer"
	"github.com/areknoster/public-distributed-commit-log/sentinel/commiter"
	"github.com/areknoster/public-distributed-commit-log/sentinel/pinner"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/sentinel/service"
	"github.com/areknoster/public-distributed-commit-log/storage"
	memorystorage "github.com/areknoster/public-distributed-commit-log/storage/memory"
	"github.com/areknoster/public-distributed-commit-log/test/testpb"
)

type ProduceConsumeTestSuite struct {
	suite.Suite
	messageStorage storage.MessageStorage
	grpcServer     *grpc.Server
	sentinelClient sentinelpb.SentinelClient
	producer       *producer.MessageProducer
	headReader     *sentinel_reader.SentinelHeadReader
	globalCtx      context.Context
}

func (s *ProduceConsumeTestSuite) SetupSuite() {
	s.setupSentinel()
	s.setupProducer()
	s.setupConsumerDependencies()
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	g, ctx := errgroup.WithContext(ctx)
	s.globalCtx = ctx
	g.Go(s.grpcServer.ListenAndServe)
	g.Go(func() error {
		<-ctx.Done()
		s.grpcServer.GracefulStop()
		return nil
	})
	s.T().Cleanup(func() {
		cancel()
		s.Require().NoError(g.Wait())
	})
}

func (s *ProduceConsumeTestSuite) setupSentinel() {
	contentStorage := memorystorage.Storage{}
	s.messageStorage = storage.NewProtoMessageStorage(contentStorage)
	validator := &mockValidator{
		messageStorage: s.messageStorage,
		t:              s.T(),
	}
	memoryPinner := pinner.NewMemoryPinner()
	headManager := memory.NewHeadManager(cid.Undef)
	instantCommiter := commiter.NewInstant(headManager, s.messageStorage, memoryPinner)
	sentinelService := service.New(validator, memoryPinner, instantCommiter, headManager)
	grpcServer, err := grpc.NewServer(grpc.ServerConfig{
		Host: "localhost",
		Port: "8000",
	})
	s.grpcServer = grpcServer
	s.Require().NoError(err)
	sentinelpb.RegisterSentinelServer(grpcServer, sentinelService)
}

func (s *ProduceConsumeTestSuite) setupProducer() {
	conn, err := grpc.Dial(grpc.ConnConfig{
		Host: "localhost",
		Port: "8000",
	})
	s.Require().NoError(err)
	s.sentinelClient = sentinelpb.NewSentinelClient(conn)
	s.producer = producer.NewMessageProducer(s.messageStorage, s.sentinelClient)
}

func (s *ProduceConsumeTestSuite) setupConsumerDependencies() {
	s.headReader = sentinel_reader.NewSentinelHeadReader(s.sentinelClient)
}

func (s *ProduceConsumeTestSuite) newConsumer(offset cid.Cid) *consumer.FirstToLastConsumer {
	return consumer.NewFirstToLastConsumer(
		s.headReader,
		memory.NewHeadManager(offset),
		s.messageStorage,
		consumer.FirstToLastConsumerConfig{
			PollInterval: time.Second,
			PollTimeout:  200 * time.Millisecond,
		})
}

func (s *ProduceConsumeTestSuite) TestProduceConsume() {
	ctx, cancel := context.WithTimeout(s.globalCtx, 10*time.Second)
	defer cancel()

	const messageNumber = 10
	messages := make([]*testpb.Message, messageNumber)
	idSet := make(map[int64]struct{}, messageNumber)
	for i := int64(0); i < messageNumber; i++ {
		messages[i] = &testpb.Message{Id: i}
		idSet[i] = struct{}{}
	}

	s.Run("when messsages are first produced, and then a consumer is started, all produced messages should be handled at the beginning", func() {
		s.consumeFromStart(ctx, messages, idSet)
	})
}

func (s *ProduceConsumeTestSuite) consumeFromStart(ctx context.Context, messages []*testpb.Message, idSet map[int64]struct{}) {
	for _, message := range messages {
		s.Require().NoError(s.producer.Produce(ctx, message))
	}
	idsChan := make(chan int64, 5)

	cons := s.newConsumer(cid.Undef)
	consumeCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		err := cons.Consume(consumeCtx, consumer.MessageFandlerFunc(func(ctx context.Context, message storage.ProtoUnmarshallable) error {
			testMessage := &testpb.Message{}
			if err := message.Unmarshall(testMessage); err != nil {
				close(idsChan)
				return err
			}
			idsChan <- testMessage.Id
			return nil
		}))
		s.Assert().NoError(err)
		close(idsChan)
	}()

	for id := range idsChan {
		delete(idSet, id)
		s.T().Logf("Received: %v", id)
		if len(idSet) == 0 {
			cancel()
		}
	}
	s.Assert().Len(idSet, 0)
}

type mockValidator struct {
	messageStorage storage.MessageStorage
	t              *testing.T
}

func (m *mockValidator) Validate(ctx context.Context, cid cid.Cid) error {
	unmarshallable, err := m.messageStorage.Read(ctx, cid)
	require.NoError(m.t, err)
	message := &testpb.Message{}
	err = unmarshallable.Unmarshall(message)
	require.NoError(m.t, err)
	return nil
}

func TestProduceConsumeTestSuite(t *testing.T) {
	ts := &ProduceConsumeTestSuite{}
	suite.Run(t, ts)
}
