// Copyright 2012 Jesse Allen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

import (
	"time"
)

// Message is the container for whatever data is being published.
type Message struct {
	OneLine   [70]byte
	TimeStamp time.Time
	FullText  string
}

// Filters are an identifiable interface used to specify
// what results a Publisher should send over a channel
type Filter interface {
	Identify() string
}

// Publishers are identifiable sources of new information.
// We only care that they provide a channel for receiving
// Messages to send to Subscribers and a channel to send
// a signal to stop sending Messages.
//
// The value sent over the boolean channel should not
// matter, though it will be true. The act of sending
// over the channel is a request to stop sending Messages.
type Publisher interface {
	Publish(Filter) (<-chan Message, chan<- bool, error)
	Identify() string
}

// Subscribers are identifiable recipients of new information.
// We only care that they provide a channel for sending
// Messages from Publishers and a channel to receive a signal
// to stop sending Messages.
//
// The value sent over the boolean channel does not matter.
// The act of sending a value over the channel is a request
// to stop sending Messages.
type Subscriber interface {
	Subscribe() (chan<- Message, <-chan bool, error)
	Identify() string
}
