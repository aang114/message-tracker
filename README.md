# Message Tracker

This repository contains an implementation of the message tracker interface `MessageTracker` defined in `network/message_tracker.go`:

```go 
// MessageTracker tracks a **configurable fixed amount of messages**.
// Messages are stored **first-in-first-out**.  **Duplicate messages** should **not** be stored in the **queue**.
type MessageTracker interface {
	// Add will add a message to the tracker
	Add(message *Message) (err error)
	// Delete will delete message from tracker
	Delete(id string) (err error)
	// Get returns a message for a given ID.  Message is retained in tracker
	Message(id string) (message *Message, err error)
	// All returns messages **in the order** in which they were received
	Messages() (messages []*Message)
}
```

The `Message` type that is found in `network/message.go`.

```go 
// Message is received from peers in a p2p network.
type Message struct {
	ID     string
	PeerID string
	Data   []byte
}
```

Each message is uniquely identified by the `Message.ID`. Messages with the same ID may be received by multiple peers. Peers are uniquely identified by their own ID stored in `Message.PeerID`.

## Implementation

A doubly linked list (where a node's value is `*Message`) and a hash map (that maps a message ID to the message's doubly linked list node) were used to store the messages. If the maximum length is about to be exceeded after an addition, the earliest element gets removed.

To prevent data races, a read-write mutex is used. A read-write mutex was used rather than a normal mutex to allow concurrent reading.

Since performance is critical, a doubly linked list was chosen instead of alternatives (such as Go slices) since it has a better time complexity for the addition, deletion and get operations. Furthermore, it has the same time complexity as alternatives such as Go slices for the retrieve-all operation (although it may be slower in practice). However, since I assumed that the first 3 operations would require faster performance and would be called more often, a doubly linked list was chosen.

### Addition, Deletion and Get

The combination of a doubly linked list and a hash map allows `MessageTracker.Add()`, `MessageTracker.Delete()` and `MessageTracker.Message()` to be performed in constant time complexity (i.e O(1)) - which is ideal for performance. However, using two data structures requires more memory usage.

If a Go Slice were used with a hash map instead, this would result in a linear time complexity of O(n) for `MessageTracker.Add()` and `MessageTracker.Delete()`. This is because `MessageTracker.Add()` may result in a reallocation if the capacity of the Go Slice is exceeded and `MessageTracker.Delete()` may result in other elements getting shifted.

### Retrieving all Messages

`MessageTracker.Messages()` is performed in linear complexity (i.e O(n)) since the doubly linked list is traversed. This is the same time complexity if a Go Slice was used instead.

In practice, a Go Slice may be faster since its elements are stored in a contiguous block of memory. Furthermore, since the return type is a slice, the doubly linked list implementation requires the elements to be copied to a newly created slice every time the function is called - which increase its time taken and requires more memory usage.
