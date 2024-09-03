package network_test

import (
	"fmt"
	"github.com/ChainSafe/gossamer-go-interview/network"
	"github.com/stretchr/testify/assert"
	"slices"
	"sync"
	"testing"
)

func generateMessage(n int) *network.Message {
	return &network.Message{
		ID:     fmt.Sprintf("someID%d", n),
		PeerID: fmt.Sprintf("somePeerID%d", n),
		Data:   []byte{0, 1, 1},
	}
}

func TestMessageTracker_Add(t *testing.T) {
	t.Run("add nil message", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		err := mt.Add(nil)
		assert.ErrorIs(t, err, network.ErrInvalidMessage)
	})

	t.Run("add message with empty ID string", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)
		msg := generateMessage(0)
		msg.ID = ""

		err := mt.Add(msg)
		assert.ErrorIs(t, err, network.ErrInvalidMessage)
	})

	t.Run("add, get, then all messages", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)

			msg, err := mt.Message(generateMessage(i).ID)
			assert.NoError(t, err)
			assert.NotNil(t, msg)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
			generateMessage(4),
		}, msgs)
	})

	t.Run("add, get, then all messages, delete some", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)

			msg, err := mt.Message(generateMessage(i).ID)
			assert.NoError(t, err)
			assert.NotNil(t, msg)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
			generateMessage(4),
		}, msgs)

		for i := 0; i < length-2; i++ {
			err := mt.Delete(generateMessage(i).ID)
			assert.NoError(t, err)
		}

		msgs = mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(3),
			generateMessage(4),
		}, msgs)

	})

	t.Run("not full, with duplicates", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length-1; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}
		for i := 0; i < length-1; i++ {
			err := mt.Add(generateMessage(length - 2))
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
		}, msgs)
	})

	t.Run("not full, with duplicates from other peers", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length-1; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}
		for i := 0; i < length-1; i++ {
			msg := generateMessage(length - 2)
			msg.PeerID = "somePeerID0"
			err := mt.Add(msg)
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(0),
			generateMessage(1),
			generateMessage(2),
			generateMessage(3),
		}, msgs)
	})

}

func TestMessageTracker_AddConcurrently(t *testing.T) {
	t.Run("concurrently add, get, then all messages", func(t *testing.T) {
		var wg sync.WaitGroup
		length := 1000
		mt := network.NewMessageTracker(length)

		for i := 0; i < length; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				msg := generateMessage(i)
				err := mt.Add(msg)
				assert.NoError(t, err)

				retrievedMessage, err := mt.Message(msg.ID)
				assert.NoError(t, err)
				assert.NotNil(t, msg, retrievedMessage)
			}(i)
		}
		wg.Wait()

		msgs := mt.Messages()
		assert.Len(t, msgs, length)
		for i := 0; i < length; i++ {
			msg := generateMessage(i)
			slices.Contains(msgs, msg)
		}
	})

	t.Run("concurrently add, delete, then all messages", func(t *testing.T) {
		var wg sync.WaitGroup
		length := 1000
		mt := network.NewMessageTracker(length)

		for i := 0; i < length; i++ {
			wg.Add(1)

			go func(i int) {
				defer wg.Done()

				msg := generateMessage(i)
				err := mt.Add(msg)
				assert.NoError(t, err)

				err = mt.Delete(msg.ID)
				assert.NoError(t, err)
			}(i)
		}
		wg.Wait()

		msgs := mt.Messages()
		assert.Empty(t, msgs)
	})
}

func TestMessageTracker_Cleanup(t *testing.T) {
	t.Run("overflow and cleanup", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length*2; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(5),
			generateMessage(6),
			generateMessage(7),
			generateMessage(8),
			generateMessage(9),
		}, msgs)
	})

	t.Run("overflow and cleanup with duplicate", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)

		for i := 0; i < length*2; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}

		for i := length; i < length*2; i++ {
			err := mt.Add(generateMessage(i))
			assert.NoError(t, err)
		}

		msgs := mt.Messages()
		assert.Equal(t, []*network.Message{
			generateMessage(5),
			generateMessage(6),
			generateMessage(7),
			generateMessage(8),
			generateMessage(9),
		}, msgs)
	})
}

func TestMessageTracker_Delete(t *testing.T) {
	t.Run("empty tracker", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)
		err := mt.Delete("bleh")
		assert.ErrorIs(t, err, network.ErrMessageNotFound)
	})
}

func TestMessageTracker_Message(t *testing.T) {
	t.Run("empty tracker", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)
		msg, err := mt.Message("bleh")
		assert.ErrorIs(t, err, network.ErrMessageNotFound)
		assert.Nil(t, msg)
	})
}

func TestMessageTracker_Messages(t *testing.T) {
	t.Run("empty tracker", func(t *testing.T) {
		length := 5
		mt := network.NewMessageTracker(length)
		messages := mt.Messages()
		assert.Empty(t, messages)
	})
}
