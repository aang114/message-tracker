# Go Interview - Gossamer

## Task Description

Implement a peer-to-peer (p2p) message tracker. There is a `Message` type that is found in `network/message.go`.  
```go 
// Message is received from peers in a p2p network.
type Message struct {
	ID     string
	PeerID string
	Data   []byte
}
```
Each message is uniquely identified by the `Message.ID`. Messages with the same ID may be received by multiple peers.  Peers are uniquely identified by their own ID stored in `Message.PeerID`. 

The interface for the message tracker is defined in `network/message_tracker.go`.  
```go 
// MessageTracker tracks a configurable fixed amount of messages.
// Messages are stored first-in-first-out.  Duplicate messages should not be stored in the queue.
type MessageTracker interface {
	// Add will add a message to the tracker
	Add(message *Message) (err error)
	// Delete will delete message from tracker
	Delete(id string) (err error)
	// Get returns a message for a given ID.  Message is retained in tracker
	Message(id string) (message *Message, err error)
	// All returns messages in the order in which they were received
	Messages() (messages []*Message)
}
```

There is an exported constructor `network.NewMessageTracker(length int)` which accepts a length parameter.  This parameter should be used to constrain the number of messages in your implementation.

There are some tests within the `network_test` package found in `network/message_tracker_test.go` which call the `network.NewMessageTracker` and test the functionality from outside the `network` package.

There are a few key points to take into account when implementing this tracker:

- The tracker is meant to be a hot path in our program so performance is critical.
- Duplicate messages based on `Message.ID` should only be returned by `MessageTracker.All()` once.
- The tracker should only hold a configurable maximum amount of messages so it does not grow in size indefinitely.

## Submission Criteria
- Implement the `MessageTracker` interface, and ensure tests in `network/message_tracker_test.go` pass.
- Write unit tests for your `MessageTracker` implementation and obtain 70%+ code coverage.
- BONUS: Write benchmarks for your tracker implementation.
- BONUS: Write a design document that describes your implementation and the technical choices that you made.

## Submission

You must use `git` to track your changes.

You can either submit us:

- a URL to your Git repository
- a zip file containing your Git repository