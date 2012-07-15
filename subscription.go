// Copyright 2012 Jesse Allen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pubsub

// subplex is our implementation of the fundamental
// pubsub communication
type subscription struct {
	pub  <-chan Message
	subs map[string]chan<- Message
	stop chan<- bool
}
