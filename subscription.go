// Copyright 2012 Jesse Allen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

// subscription is our implementation of the fundamental
// pubsub communication
type subscription struct {
	pub  <-chan Message
	subs map[string]chan<- Message
	stop chan<- bool
}

func (s *subscription) init(p Publisher, f Filter) (err error) {
	s.pub, s.stop, err = p.Publish(f)
	if err != nil {
		return err
	}
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

func (s *subscription) addSubscriber(sub Subscriber) (err error) {
	var stop <-chan bool
	k := sub.Identify()
	// subscribers only need to be added once
	sc := s.subs[k]
	if sc == nil {
		s.subs[k], stop, err = sub.Subscribe()
		if err != nil {
			delete(s.subs, k)
			return err
		}
		go func() {
			<-stop
			close(s.subs[k])
			delete(s.subs, k)
		}()
	}
	return nil
}

// removeSubscriber will remove the indicated subscriber and
// if there are none left, send a signal to stop incoming messages
// and return true to indicate that this subscription is dead.
func (s *subscription) removeSubscriber(sub Subscriber) bool {
	k := sub.Identify()
	sc := s.subs[k]
	if sc != nil {
		close(sc)
		delete(s.subs, k)
		if len(s.subs) == 0 {
			// no more subscribers
			s.stop <- true
			return true
		}
	}
	return false
}
