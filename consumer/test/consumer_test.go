package test

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/areknoster/public-distributed-commit-log/consumer"
	"github.com/areknoster/public-distributed-commit-log/ipns"
	"github.com/areknoster/public-distributed-commit-log/pdclpb"
	"github.com/areknoster/public-distributed-commit-log/storage"
	memorystorage "github.com/areknoster/public-distributed-commit-log/storage/content/memory"
	messagestorage "github.com/areknoster/public-distributed-commit-log/storage/message"
	"github.com/areknoster/public-distributed-commit-log/storage/pbcodec"
	"github.com/areknoster/public-distributed-commit-log/test/testpb"
	"github.com/areknoster/public-distributed-commit-log/test/testutil"
	"github.com/areknoster/public-distributed-commit-log/thead"
	"github.com/areknoster/public-distributed-commit-log/thead/memory"
)

var consumerCreators = []createConsumer{
	createFirstToLastConsumer,
}

var ipnsMgrResolver = ipns.NewTestManager()

func TestEdgeCases(t *testing.T) {
	for _, creator := range consumerCreators {
		t.Run("for empty topic consumer should wait and finish with timeout error when context is done", func(t *testing.T) {
			mockReader := newMockMessageReader(t)
			c := creator(cid.Undef, mockReader, memory.NewHeadManager(cid.Undef))
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			err := c.Consume(ctx, consumer.MessageHandlerFunc(func(ctx context.Context, message storage.ProtoDecodable) error {
				t.Fatal("no message should be handled")
				return nil
			}))
			assert.ErrorIs(t, err, consumer.ErrContextDone)
		})

		t.Run("for topic with single empty commit consumer should  wait and finish with timeout error when context is done", func(t *testing.T) {
			mockReader := newMockMessageReader(t)
			head := mockReader.Commit(cid.Undef)
			c := creator(cid.Undef, mockReader, memory.NewHeadManager(head))
			require.NoError(t, ipnsMgrResolver.UpdateIPNSEntry(head.String()))
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			err := c.Consume(ctx, consumer.MessageHandlerFunc(func(ctx context.Context, message storage.ProtoDecodable) error {
				t.Fatal("no message should be handled")
				return nil
			}))
			assert.ErrorIs(t, err, consumer.ErrContextDone)
		})

		t.Run("single commit with single message and commit added later should be both consumed once despite non-existing offset", func(t *testing.T) {
			mockReader := newMockMessageReader(t)
			th := newTestHandler(t)
			mockReader.RegisterMessage(th.AddNext())
			head := mockReader.Commit(cid.Undef)
			headManager := memory.NewHeadManager(head)
			require.NoError(t, ipnsMgrResolver.UpdateIPNSEntry(head.String()))
			c := creator(testutil.RandomCID(t), mockReader, headManager)
			ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
			defer cancel()
			go func() {
				err := c.Consume(ctx, th)
				assert.ErrorIs(t, err, consumer.ErrContextDone)
			}()
			time.Sleep(10 * time.Millisecond)
			mockReader.RegisterMessage(th.AddNext())
			head = mockReader.Commit(head)
			require.NoError(t, headManager.SetHead(ctx, head))
			require.NoError(t, ipnsMgrResolver.UpdateIPNSEntry(head.String()))
			<-ctx.Done()
			th.AssertAllHandledOnce()
		})
	}
}

func TestLinearIncremental(t *testing.T) {
	for _, creator := range consumerCreators {
		t.Run("for empty initial consumer offset and incrementally moved topic head, every message should be consumed", func(t *testing.T) {
			mockReader := newMockMessageReader(t)
			th := newTestHandler(t)
			headManager := memory.NewHeadManager(cid.Undef)
			c := creator(testutil.RandomCID(t), mockReader, headManager)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				err := c.Consume(ctx, th)
				assert.ErrorIs(t, err, consumer.ErrContextDone)
			}()
			head := cid.Undef
			for i := 0; i < 100; i++ {
				mockReader.RegisterMessage(th.AddNext())
				head = mockReader.Commit(head)
				require.NoError(t, headManager.SetHead(ctx, head))
				require.NoError(t, ipnsMgrResolver.UpdateIPNSEntry(head.String()))
				time.Sleep(time.Millisecond)
			}

			time.Sleep(100 * time.Millisecond)
			cancel()
			th.AssertAllHandledOnce()
		})

		t.Run("message from offset should not be read", func(t *testing.T) {
			mockReader := newMockMessageReader(t)
			mockReader.RegisterMessage(&testpb.Message{IdIncremental: math.MaxInt64})
			head := mockReader.Commit(cid.Undef)

			th := newTestHandler(t)
			headManager := memory.NewHeadManager(head)
			c := creator(head, mockReader, headManager) // start from current head offset
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				err := c.Consume(ctx, th)
				assert.ErrorIs(t, err, consumer.ErrContextDone)
			}()
			mockReader.RegisterMessage(th.AddNext())
			head = mockReader.Commit(head)
			require.NoError(t, headManager.SetHead(ctx, head))
			require.NoError(t, ipnsMgrResolver.UpdateIPNSEntry(head.String()))

			time.Sleep(100 * time.Millisecond)
			cancel()
			th.AssertAllHandledOnce() // would fail if the initial message is read
		})
	}
}

func TestRandomSizedCommits(t *testing.T) {
	for _, creator := range consumerCreators {
		t.Run("for empty initial consumer offset and incrementally moved topic head, every message should be consumed", func(t *testing.T) {
			mockReader := newMockMessageReader(t)
			th := newTestHandler(t)
			headManager := memory.NewHeadManager(cid.Undef)
			c := creator(testutil.RandomCID(t), mockReader, headManager)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			go func() {
				err := c.Consume(ctx, th)
				assert.ErrorIs(t, err, consumer.ErrContextDone)
			}()
			head := cid.Undef
			for i := 0; i < 100; i++ {
				messagesInCommit := rand.Intn(100)
				for j := 0; j < messagesInCommit; j++ {
					mockReader.RegisterMessage(th.AddNext())
				}
				head = mockReader.Commit(head)
				require.NoError(t, headManager.SetHead(ctx, head))
				require.NoError(t, ipnsMgrResolver.UpdateIPNSEntry(head.String()))
				time.Sleep(5 * time.Millisecond)
			}

			time.Sleep(100 * time.Millisecond)
			cancel()
			th.AssertAllHandledOnce()
		})
	}
}

func newMockMessageReader(t *testing.T) *mockMessageReader {
	messageStorage := messagestorage.NewContentStorageWrapper(memorystorage.NewStorage(), pbcodec.ProtoBuf{})
	return &mockMessageReader{
		messageStorage: messageStorage,
		t:              t,
	}
}

type accessor func(ctx context.Context) (storage.ProtoDecodable, error)

// mockMessageReader registers accessors for messages that can simulate multiple storage events, e.g. timeout, errors, correct access
type mockMessageReader struct {
	uncommitedMessages []cid.Cid
	messageStorage     storage.MessageStorage
	accessors          sync.Map
	t                  *testing.T
}

func (m *mockMessageReader) Read(ctx context.Context, cid cid.Cid) (storage.ProtoDecodable, error) {
	v, found := m.accessors.Load(cid)
	require.True(m.t, found, fmt.Sprintf("attempt to acess message %s which was not registered", cid))
	accessor := v.(accessor)
	return accessor(ctx)
}

func (m *mockMessageReader) register(c cid.Cid, f func(ctx context.Context) (storage.ProtoDecodable, error)) {
	m.accessors.Store(c, accessor(f))
}

// there will be an error when message is accessed
func (m *mockMessageReader) RegisterMessageWithError(message proto.Message, err error) {
	messageCid, writeErr := m.messageStorage.Write(context.TODO(), message)
	require.NoError(m.t, writeErr)
	m.uncommitedMessages = append(m.uncommitedMessages, messageCid)
	m.register(messageCid, func(ctx context.Context) (storage.ProtoDecodable, error) {
		return nil, err
	})
}

// message will be accessed without issues
func (m *mockMessageReader) RegisterMessage(message *testpb.Message) {
	messageCid, writeErr := m.messageStorage.Write(context.TODO(), message)
	require.NoError(m.t, writeErr)
	m.uncommitedMessages = append(m.uncommitedMessages, messageCid)
	m.register(messageCid, func(ctx context.Context) (storage.ProtoDecodable, error) {
		message, err := m.messageStorage.Read(context.TODO(), messageCid)
		require.NoError(m.t, err)
		return message, nil
	})
}

func (m *mockMessageReader) writeCommit(previous cid.Cid) cid.Cid {
	stringCIDs := make([]string, len(m.uncommitedMessages))
	for i, message := range m.uncommitedMessages {
		stringCIDs[i] = message.String()
	}
	m.uncommitedMessages = m.uncommitedMessages[:0]
	messageCid, writeErr := m.messageStorage.Write(context.TODO(), &pdclpb.Commit{
		Created:           timestamppb.Now(),
		PreviousCommitCid: previous.String(),
		MessagesCids:      stringCIDs,
	})
	require.NoError(m.t, writeErr)
	return messageCid
}

func (m *mockMessageReader) Commit(previous cid.Cid) cid.Cid {
	messageCid := m.writeCommit(previous)
	m.register(messageCid, func(ctx context.Context) (storage.ProtoDecodable, error) {
		message, err := m.messageStorage.Read(context.TODO(), messageCid)
		require.NoError(m.t, err)
		return message, nil
	})
	return messageCid
}

func (m *mockMessageReader) CommitWithError(previous cid.Cid) cid.Cid {
	messageCid := m.writeCommit(previous)
	m.register(messageCid, func(ctx context.Context) (storage.ProtoDecodable, error) {
		message, err := m.messageStorage.Read(context.TODO(), messageCid)
		require.NoError(m.t, err)
		return message, nil
	})
	return messageCid
}

type createConsumer func(initialOffset cid.Cid, messagesTree storage.MessageReader, headReader thead.Reader) consumer.Consumer

var createFirstToLastConsumer createConsumer = func(initialOffset cid.Cid, messagesTree storage.MessageReader, headReader thead.Reader) consumer.Consumer {
	ipnsMgrResolver.UpdateIPNSEntry(initialOffset.String())
	return consumer.NewFirstToLastConsumer(memory.NewHeadManager(initialOffset), messagesTree, consumer.FirstToLastConsumerConfig{
		PollInterval: 50 * time.Millisecond,
		PollTimeout:  25 * time.Millisecond,
	}, ipnsMgrResolver, "")
}

type testHandler struct {
	index    int64
	messages sync.Map
	t        *testing.T
}

func newTestHandler(t *testing.T) *testHandler {
	return &testHandler{
		t: t,
	}
}

func (th *testHandler) Handle(ctx context.Context, message storage.ProtoDecodable) error {
	testMessage := &testpb.Message{}
	require.NoError(th.t, message.Decode(testMessage))

	currentValue, found := th.messages.Load(testMessage.IdIncremental)
	require.True(th.t, found, "handle on unknown message")
	incr := currentValue.(int) + 1
	th.messages.Store(testMessage.IdIncremental, incr)
	return nil
}

func (th *testHandler) AssertAllHandledOnce() {
	th.messages.Range(func(index, handledTimes interface{}) bool {
		assert.EqualValues(th.t, 1, handledTimes, index)
		return true
	})
}

func (th *testHandler) Add(m *testpb.Message) {
	v := int(0)
	th.messages.Store(m.IdIncremental, v)
}

func (th *testHandler) Next() *testpb.Message {
	th.index++
	return &testpb.Message{IdIncremental: th.index}
}

func (th *testHandler) AddNext() *testpb.Message {
	next := th.Next()
	th.Add(next)
	return next
}
