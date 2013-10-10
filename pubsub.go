// Copyright 2012 Jesse Allen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

import (
	"errors"
	"time"
)

const (
	errPublisherExists   = "Error, Published channel already exists in PubSub"
	errAlreadySubscribed = "Error, Subscriber channel already subscribed to Published channel"
)

// Message is the container for whatever data is being published.
type Message struct {
	OneLine   [70]byte
	TimeStamp time.Time
	FullText  string
}

// Publishers are sources of new information.
// We only care that they provide a channel for receiving
// Messages to send to Subscribers and a channel to send
// a signal to stop sending Messages.
type Publisher func(filter string) (m <-chan Message, stop chan<- bool, err error)

// Subscribers are recipients of new information.
// We only care that they provide a channel for sending
// Messages from Publishers and a channel to receive a signal
// to stop sending Messages.
type Subscriber func() (m chan<- Message, stop <-chan bool, err error)

// PubSub manages the Subscriber/Publisher/Filter relationship.
type PubSub struct {
	subscriptions [](*subscription)
}

// TODO: Currently no bootstrapping here, but there may be soon enough
func New() *PubSub {
	ps := new(PubSub)
	return ps
}

// Register adds new Publisher/Filter to Subscriber relationships.
// These relationships are oriented around Publisher/Filter instead
// of Subscriber because the Publisher is initiating the communication
// of Messages.
//
// Subscribers will not be added more than once. It is not currently
// and error state.
//
// TODO: Consider returning error for already subscribed Subscriber.
func (ps *PubSub) Register(p Publisher, filter string, sub Subscriber) error {
	s, err := ps.newSubscription(p, filter)
	if err != nil && err.Error() != errPublisherExists {
		return err
	} else if err.Error() == errPublisherExists {
		// TODO [jesse@jessecarl.com][2013-10-10]: add logging for these conditions
	}
	err = s.addSubscriber(sub)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PubSub) newSubscription(p Publisher, filter string) (*subscription, error) {
	s := new(subscription)
	if err := s.init(p, filter); err != nil {
		return nil, err
	}
	// prevent publisher from being added more than once
	for _, existingSub := range ps.subscriptions {
		if existingSub.pub == s.pub {
			return existingSub, errors.New(errPublisherExists)
		}
	}
	s.start()
	ps.subscriptions = append(ps.subscriptions, s)
	return s, nil
}
