package network

import (
	"container/list"
	"sync"
)

// Struct that implements the interface `MessageTracker`
type tracker struct {
	// Read-Write Mutex to prevent data races
	mu sync.RWMutex
	// Maximum number of messages that can be stored in the message tracker
	maxLength int
	// Doubly Linked list of messages where they are stored in the order they are received
	linkedList *list.List
	// Map that maps a message's ID to the message's linked list node
	idToElementMap map[string]*list.Element
}

// Creates a new `tracker` instance with maximum length `maxLength`
func newTracker(maxLength int) *tracker {
	return &tracker{
		maxLength:      maxLength,
		linkedList:     list.New(),
		idToElementMap: make(map[string]*list.Element),
	}
}

func (t *tracker) Add(message *Message) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// sanity check of the message
	if message == nil || message.ID == "" {
		return ErrInvalidMessage
	}

	// return early if message ID already exists
	if _, ok := t.idToElementMap[message.ID]; ok {
		return nil
	}

	// pop earliest element if adding new element will exceed max length
	if t.linkedList.Len() >= t.maxLength {
		firstElement := t.linkedList.Front()
		t.linkedList.Remove(firstElement)
		delete(t.idToElementMap, firstElement.Value.(*Message).ID)
	}
	// add message
	newElement := t.linkedList.PushBack(message)
	t.idToElementMap[message.ID] = newElement

	return nil
}

func (t *tracker) Delete(id string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	element, ok := t.idToElementMap[id]
	if !ok {
		return ErrMessageNotFound
	}
	t.linkedList.Remove(element)
	delete(t.idToElementMap, id)

	return nil
}

func (t *tracker) Message(id string) (*Message, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	element, ok := t.idToElementMap[id]
	if !ok {
		return nil, ErrMessageNotFound
	}
	return element.Value.(*Message), nil
}

func (t *tracker) Messages() []*Message {
	t.mu.RLock()
	defer t.mu.RUnlock()

	messages := make([]*Message, t.linkedList.Len())

	i := 0
	for element := t.linkedList.Front(); element != nil; element = element.Next() {
		messages[i] = element.Value.(*Message)
		i++
	}

	return messages
}
