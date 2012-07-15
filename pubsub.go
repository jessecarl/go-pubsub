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
