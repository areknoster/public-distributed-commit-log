package itest

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
	"github.com/areknoster/public-distributed-commit-log/ipns"
	"github.com/areknoster/public-distributed-commit-log/producer"
	"github.com/areknoster/public-distributed-commit-log/ratelimiting"
	"github.com/areknoster/public-distributed-commit-log/sentinel/commiter"
	"github.com/areknoster/public-distributed-commit-log/sentinel/pinner"
	"github.com/areknoster/public-distributed-commit-log/sentinel/sentinelpb"
	"github.com/areknoster/public-distributed-commit-log/sentinel/service"
	"github.com/areknoster/public-distributed-commit-log/storage"
	memorystorage "github.com/areknoster/public-distributed-commit-log/storage/memory"
	"github.com/areknoster/public-distributed-commit-log/test/testpb"
	"github.com/areknoster/public-distributed-commit-log/thead/memory"
	"github.com/areknoster/public-distributed-commit-log/thead/sentinelhead"
)

// ProduceConsumeTestSuite can be used to create tests with multiple DI setups, based on setup methods overrides
type ProduceConsumeTestSuite struct {
	suite.Suite

	messageStorage storage.MessageStorage
	grpcServer     *grpc.Server
	sentinelClient sentinelpb.SentinelClient
	producer       *producer.MessageProducer
	globalCtx      context.Context
}

func (s *ProduceConsumeTestSuite) SetupSuite() {
	s.setupMessageStorage()
	s.setupSentinel()
	s.setupSentinelClient()
	s.setupProducer()
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
	validator := &messageValidator{
		messageStorage: s.messageStorage,
		t:              s.T(),
	}
	memoryPinner := pinner.NewMemoryPinner()
	headManager := memory.NewHeadManager(cid.Undef)
	instantCommiter := commiter.NewInstant(headManager, s.messageStorage, memoryPinner, ipns.NewNopManager())
	sentinelService := service.New(validator, memoryPinner, instantCommiter, headManager, ipns.NewNopManager())
	grpcServer, err := grpc.NewServer(grpc.ServerConfig{
		Host: "localhost",
		Port: "8000",
	}, ratelimiting.NewAlwaysAllowLimiter())
	s.grpcServer = grpcServer
	s.Require().NoError(err)
	sentinelpb.RegisterSentinelServer(grpcServer, sentinelService)
}

func (s *ProduceConsumeTestSuite) setupMessageStorage() {
	if s.messageStorage != nil {
		return
	}
	contentStorage := &memorystorage.Storage{}
	s.messageStorage = storage.NewProtoMessageStorage(contentStorage)
}

func (s *ProduceConsumeTestSuite) setupProducer() {
	s.producer = producer.NewMessageProducer(s.messageStorage, s.sentinelClient)
}

func (s *ProduceConsumeTestSuite) setupSentinelClient() {
	conn, err := grpc.Dial(grpc.ConnConfig{
		Host: "localhost",
		Port: "8000",
	})
	s.Require().NoError(err)
	s.sentinelClient = sentinelpb.NewSentinelClient(conn)
}

func (s *ProduceConsumeTestSuite) newConsumer(offset cid.Cid) consumer.Consumer {
	return consumer.NewFirstToLastConsumer(
		sentinelhead.New(s.sentinelClient),
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
		messages[i] = &testpb.Message{IdIncremental: i}
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
		err := cons.Consume(consumeCtx, consumer.MessageHandlerFunc(func(ctx context.Context, message storage.ProtoUnmarshallable) error {
			testMessage := &testpb.Message{}
			if err := message.Unmarshall(testMessage); err != nil {
				close(idsChan)
				return err
			}
			idsChan <- testMessage.IdIncremental
			return nil
		}))
		s.Assert().ErrorIs(err, consumer.ErrContextDone)
		close(idsChan)
	}()

	for id := range idsChan {
		delete(idSet, id)
		if len(idSet) == 0 {
			s.T().Log(len(idSet))
			cancel()
		}
	}
	s.Assert().Len(idSet, 0)
}

type messageValidator struct {
	messageStorage storage.MessageStorage
	t              *testing.T
}

func (m *messageValidator) Validate(ctx context.Context, cid cid.Cid) error {
	unmarshallable, err := m.messageStorage.Read(ctx, cid)
	require.NoError(m.t, err)
	message := &testpb.Message{}
	err = unmarshallable.Unmarshall(message)
	require.NoError(m.t, err)
	return nil
}
