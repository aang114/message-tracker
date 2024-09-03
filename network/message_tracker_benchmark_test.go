package network_test

import (
	"github.com/ChainSafe/gossamer-go-interview/network"
	"math/rand"
	"testing"
)

func generateMessages(length int) []*network.Message {
	messages := make([]*network.Message, length)
	for i := 0; i < length; i++ {
		messages[i] = generateMessage(i)
	}
	return messages
}

func BenchmarkMessageTracker_Add(b *testing.B) {
	mt := network.NewMessageTracker(b.N)
	messages := generateMessages(b.N)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mt.Add(messages[i])
	}
}

// Calls `MessageTracker.Delete()` on every added message in a random order
//
// Note: `MessageTracker.Delete()` is only called on messages that exist
func BenchmarkMessageTracker_Delete(b *testing.B) {
	mt := network.NewMessageTracker(b.N)
	messages := generateMessages(b.N)

	for i := 0; i < b.N; i++ {
		_ = mt.Add(messages[i])
	}

	ids := make([]string, b.N)
	for _, index := range rand.Perm(b.N) {
		ids[index] = messages[index].ID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mt.Delete(ids[i])
	}
}

// Calls `MessageTracker.Message()` on every added message in a random order
//
// Note: `MessageTracker.Message()` is only called on messages that exist
func BenchmarkMessageTracker_Message(b *testing.B) {
	mt := network.NewMessageTracker(b.N)
	messages := generateMessages(b.N)

	for i := 0; i < b.N; i++ {
		_ = mt.Add(messages[i])
	}

	ids := make([]string, b.N)
	for _, index := range rand.Perm(b.N) {
		ids[index] = messages[index].ID
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = mt.Message(ids[i])
	}
}

func BenchmarkMessageTracker_Messages(b *testing.B) {
	mt := network.NewMessageTracker(1000)
	messages := generateMessages(1000)

	for i := 0; i < 1000; i++ {
		_ = mt.Add(messages[i])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mt.Messages()
	}
}
