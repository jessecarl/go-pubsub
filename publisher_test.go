// Copyright 2012 Jesse Allen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub_test

import (
	. "./"
)

type testPublisher struct {
	Filters map[string]struct {
		Message chan Message
		Stop    chan bool
		Err     error
	}
}

func (p *testPublisher) Publish(filter string) (<-chan Message, chan<- bool, error) {
	r := p.Filters[filter]
	// even returning the "zero-value" here is perfectly acceptable
	return r.Message, r.Stop, r.Err
}

type testSubscriber struct {
	Messages chan Message
	Stop     chan bool
	Err      error
}

func (s *testSubscriber) Subscribe() (<-chan Message, chan<- bool, error) {
	return s.Messages, s.Stop, s.Err
}
