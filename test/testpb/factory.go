package testpb

import (
	"math/rand"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func MakeCurrentRandomTestMessage() *Message {
	return &Message{
		IdIncremental: rand.Int63(),
		Uuid:          uuid.NewString(),
		Created:       timestamppb.Now(),
	}
}
