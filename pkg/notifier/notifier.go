package notifier

import (
	"log"
	"sync"
)

const (
	ChangeOperationInsert string = "insert"
	ChangeOperationUpdate string = "update"
	ChangeOperationDelete string = "delete"
)

type ChangeData struct {
	OperationType string `json:"operationType"`
	UserId        string `json:"id"`
}

type Notifier interface {
	AddSubscriber(id string) <-chan ChangeData
	RemoveSubscriber(id string)
	Broadcast(msg ChangeData)
	Close()
}

type NotifierImpl struct {
	receiverChannels map[string]chan ChangeData
	mu               sync.Mutex
}

func NewNotifier() NotifierImpl {
	return NotifierImpl{receiverChannels: map[string]chan ChangeData{}}
}

func (n *NotifierImpl) AddSubscriber(id string) <-chan ChangeData {
	log.Print("Adding a new subscriber, ", id)

	n.mu.Lock()
	defer n.mu.Unlock()

	ch := make(chan ChangeData, 10)
	n.receiverChannels[id] = ch

	return ch
}

func (n *NotifierImpl) RemoveSubscriber(id string) {
	n.mu.Lock()
	defer n.mu.Unlock()

	_, exists := n.receiverChannels[id]
	if !exists {
		return
	}

	log.Print("Removing a subscriber, ", id)

	close(n.receiverChannels[id])

	delete(n.receiverChannels, id)
}

func (n *NotifierImpl) Broadcast(msg ChangeData) {
	log.Printf("Broadcasting a message to subscribers, (%s from %s)", msg.OperationType, msg.UserId)

	n.mu.Lock()
	defer n.mu.Unlock()

	go func() {
		for _, ch := range n.receiverChannels {
			ch <- msg
		}
	}()
}

func (n *NotifierImpl) Close() {
	log.Print("Closing the notifier")

	n.mu.Lock()
	defer n.mu.Unlock()

	for _, ch := range n.receiverChannels {
		close(ch)
	}
}
