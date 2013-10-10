// Copyright 2012 Jesse Allen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

import (
	"errors"
)

// subscription is our implementation of the fundamental
// pubsub communication
type subscription struct {
	pub  <-chan Message
	subs []chan<- Message
	stop chan<- bool
}

func (s *subscription) init(p Publisher, filter string) (err error) {
	s.pub, s.stop, err = p(filter)
	if err != nil {
		return err
	}
	return nil
}

func (s *subscription) start() {
	// TODO [jesse@jessecarl.com][2013-10-10]: actually start listening
	return
}

/*
  // TODO: move this to a start method
	go func() {
		for message := range s.pub {
			for _, sub := range s.subs {
				sub <- message
			}
		}
		for k, sub := range s.subs {
			close(sub)
			delete(s.subs, k)
		}
	}()
	return nil
}
*/

func (s *subscription) addSubscriber(sub Subscriber) error {
	m, stop, err := sub()
	if err != nil {
		return err
	}
	for _, existing := range s.subs {
		if existing == m {
			return errors.New(errAlreadySubscribed)
		}
	}
	s.subs = append(s.subs, m)
	go func() {
		<-stop
		s.removeSubscriber(m)
	}()
	return nil
}

// removeSubscriber will remove the indicated subscriber and
// if there are none left, send a signal to stop incoming messages
// and return true to indicate that this subscription is dead.
func (s *subscription) removeSubscriber(sub chan<- Message) bool {
	for i, existing := range s.subs {
		if existing == sub {
			s.subs = append(s.subs[:i], s.subs[i+1:]...)
			break
		}
	}
	close(sub)
	if len(s.subs) == 0 {
		// no more subscribers, signal the publisher to stop
		close(s.stop)
		return true
	}
	return false
}
