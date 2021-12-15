/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */

package nrjmx

import "bytes"

var maxBufferSize = 1024 * 1024

// LimitedBuffer will ensure that the buffer does not exceed the maxCap.
// When maxCap is reached old data will be truncated.
type LimitedBuffer struct {
	maxCap int
	buff   bytes.Buffer
}

// NewDefaultLimitedBuffer returns a LimitedBuffer with a maximum capacity of maxBufferSize.
func NewDefaultLimitedBuffer() *LimitedBuffer {
	return NewLimitedBuffer(maxBufferSize)
}

// NewLimitedBuffer returns a LimitedBuffer with a maximum capacity of maxCap.
func NewLimitedBuffer(maxCap int) *LimitedBuffer {
	return &LimitedBuffer{
		maxCap: maxCap,
	}
}

// Write appends data to the buffer. If the the maxCap is exceeded old data will be truncated.
func (lb *LimitedBuffer) Write(p []byte) (int, error) {
	if len(p) > lb.maxCap {
		p = p[len(p)-lb.maxCap:]
	}
	if len(p)+lb.buff.Len() > lb.maxCap {
		data := lb.buff.String()
		data = data[(len(p)+lb.buff.Len())-lb.maxCap:]
		lb.buff.Reset()
		_, err := lb.buff.Write([]byte(data))
		if err != nil {
			return 0, err
		}
	}
	return lb.buff.Write(p)
}

// WriteString writes to the buffer.
func (lb *LimitedBuffer) WriteString(p string) (int, error) {
	return lb.Write([]byte(p))
}

// String returns the value from the buffer.
func (lb *LimitedBuffer) String() string {
	return lb.buff.String()
}
