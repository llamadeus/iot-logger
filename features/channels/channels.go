// This package contains the code implementing the queries,
// mutations and subscriptions defined in logger.graphql.
package channels

import (
	"context"
	"github.com/llamadeus/iot-logger/graph/types"
	"github.com/llamadeus/iot-logger/internal/tailbuf"
	"github.com/llamadeus/iot-logger/internal/utils"
	"sync"
	"time"
)

// Defines a channel which a client can push data to.
type Channel struct {
	*sync.Mutex
	Name      string
	History   *tailbuf.TailBuf
	LastUsed  time.Time
	Listeners map[string]*Listener
}

// Defines a listener within a channel.
type Listener struct {
	Id      string
	Channel chan *types.Message
}

const historyLength = 100

var (
	mutex    sync.Mutex
	channels = make(map[string]*Channel)
)

func init() {
	ticker := time.NewTicker(1 * time.Minute)
	quit := make(chan struct{})

	go func() {
		for {
			select {
			case <-ticker.C:
				cleanChannels()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

// Get the message history for the given channel.
func History(ctx context.Context, channelName string) ([]*types.Message, error) {
	channel := getChannel(channelName)
	messages := messagesFromHistory(channel.History)

	return messages, nil
}

// Add a message to the given channel.
func AddMessage(ctx context.Context, channelName string, message string) (bool, error) {
	channel := getChannel(channelName)
	msg := types.Message{
		ID:        utils.NewID(),
		Text:      message,
		Timestamp: time.Now(),
	}

	channel.History.Write(&msg)

	received := false
	for _, listener := range channel.Listeners {
		select {
		case listener.Channel <- &msg:
			received = true
		}
	}

	return received, nil
}

// Subscription handler for added messages.
func MessageAdded(ctx context.Context, channelName string) (<-chan *types.Message, error) {
	channel := getChannel(channelName)
	listener := enterChannel(channel)

	go func() {
		<-ctx.Done()
		leaveChannel(channel, listener)
	}()

	return listener.Channel, nil
}

// Get the channel by its name.
// A new channel is created if the channel did not exist.
func getChannel(channelName string) *Channel {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := channels[channelName]; !ok {
		channels[channelName] = &Channel{
			Mutex:     new(sync.Mutex),
			Name:      channelName,
			History:   tailbuf.New(historyLength),
			LastUsed:  time.Now(),
			Listeners: make(map[string]*Listener),
		}
	}

	channels[channelName].LastUsed = time.Now()

	return channels[channelName]
}

// Clean unused channels.
func cleanChannels() {
	mutex.Lock()
	defer mutex.Unlock()

	for key, channel := range channels {
		if len(channel.Listeners) > 0 {
			continue
		}

		elapsed := time.Now().Sub(channel.LastUsed)

		if elapsed.Minutes() >= 1 {
			delete(channels, key)
		}
	}
}

// Enters a channel, returning a Listener.
func enterChannel(channel *Channel) *Listener {
	channel.Lock()
	defer channel.Unlock()

	var id string

	for {
		id = utils.NewID()

		if _, ok := channel.Listeners[id]; !ok {
			break
		}
	}

	channel.Listeners[id] = &Listener{
		Id:      id,
		Channel: make(chan *types.Message),
	}

	return channel.Listeners[id]
}

// Removes a listener from the given channel.
func leaveChannel(channel *Channel, listener *Listener) {
	channel.Lock()
	defer channel.Unlock()

	delete(channel.Listeners, listener.Id)
}

// Transform messages within the history buffer to a messages array.
func messagesFromHistory(history *tailbuf.TailBuf) []*types.Message {
	data := history.Read()
	result := make([]*types.Message, len(data))

	for i, value := range data {
		result[i] = value.(*types.Message)
	}

	return result
}
