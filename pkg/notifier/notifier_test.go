package notifier

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNotifier(t *testing.T) {
	notifier := NewNotifier()
	defer notifier.Close()

	t.Run("AddSubscriber", func(t *testing.T) {
		ch := notifier.AddSubscriber("user1")
		assert.NotNil(t, ch)
		assert.Equal(t, 1, len(notifier.receiverChannels))
	})

	t.Run("RemoveSubscriber", func(t *testing.T) {
		notifier.AddSubscriber("user1")
		notifier.RemoveSubscriber("user1")
		assert.Equal(t, 0, len(notifier.receiverChannels))
	})

	t.Run("Broadcast", func(t *testing.T) {
		ch := notifier.AddSubscriber("user1")
		msg := ChangeData{OperationType: ChangeOperationInsert, UserId: "user1"}
		notifier.Broadcast(msg)

		select {
		case received := <-ch:
			assert.Equal(t, msg, received)
		case <-time.After(time.Second):
			t.Fatal("Expected to receive a message")
		}
	})

	t.Run("BroadcastNonBlocking", func(t *testing.T) {
		ch := notifier.AddSubscriber("user2")

		for i := 0; i < 10; i++ {
			notifier.Broadcast(ChangeData{OperationType: ChangeOperationInsert, UserId: "user2"})
		}

		notifier.Broadcast(ChangeData{OperationType: ChangeOperationInsert, UserId: "overflow"})

		count := 0
		for {
			select {
			case <-ch:
				count++
			case <-time.After(time.Millisecond * 100):
				assert.Equal(t, 10, count)
				return
			}
		}
	})

	t.Run("CloseNotifier", func(t *testing.T) {
		ch := notifier.AddSubscriber("user3")
		notifier.Close()

		select {
		case _, ok := <-ch:
			assert.False(t, ok, "Expected channel to be closed")
		default:
			t.Fatal("Expected channel to be closed")
		}
	})
}
