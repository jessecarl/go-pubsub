// Copyright 2012 Jesse Allen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

import (
	"hash/fnv"
	"io"
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

// PubSub manages the Subscriber/Publisher/Filter relationship.
type PubSub struct {
	subscriptions map[string](*subscription)
}

func New() *PubSub {
	ps := new(PubSub)
	return ps
}

// Register adds new Publisher/Filter to Subscriber relationships.
// These relationships are oriented around Publisher/Filter instead
// of Subscriber because the Publisher is initiating the communication
// of Messages.
//
// Subscribers cannot be added more than once. This is enforced silently
// rather than by returning an error value as the caller should not know
// or care about the current subscription state.
//
// If there is an error with the Publisher, nothing will be added.
// If there is an error with a Subscriber, that Subscriber and any
// listed after will not be added, but subscribers already added
// will remain.
func (ps *PubSub) Register(p Publisher, f Filter, subs ...Subscriber) (err error) {
	k := generateKey(p, f)
	s := ps.subscriptions[k]
	// publishers only need to be added once
	if s == nil {
		s= new(subscription)
		err = s.init(p, f)
		if err != nil {
			return err
		}
	}
	for _, sub := range subs {
		err = s.addSubscriber(sub)
		if err != nil {
			return err
		}
	}
	return nil
}

// a simple way of combining publisher identifiers with filter
// identifiers
func generateKey(p Publisher, f Filter) string {
	h := fnv.New64a()
	io.WriteString(h, p.Identify())
	io.WriteString(h, f.Identify())
	return string(h.Sum(nil))
}

// UnRegister removes existing Publisher/Filter to Subscriber(s)
// relationships.
func (ps *PubSub) UnRegister(p Publisher, f Filter, subs ...Subscriber) {
	k := generateKey(p, f)
	s := ps.subscriptions[k]
	if s == nil {
		// no subscribers to unregister
		return
	}

	for _, sub := range subs {
		if dead := s.removeSubscriber(sub); dead {
			delete(ps.subscriptions, k)
		}
	}
}
