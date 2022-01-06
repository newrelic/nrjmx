/*
 * Copyright 2021 New Relic Corporation. All rights reserved.
 * SPDX-License-Identifier: Apache-2.0
 */
package nrjmx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWriteString_LengthGreaterThanMaxCap(t *testing.T) {
	// GIVEN a LimitedBuffer with maximum capacity
	maxCap := 4
	buff := NewLimitedBuffer(maxCap)

	// WHEN a value that exceeded the capacity is added
	n, err := buff.WriteString("12345")

	// THEN no error is returned and expected value is stored
	assert.NoError(t, err)
	assert.Equal(t, 4, n)

	n, err = buff.WriteString("67")
	assert.NoError(t, err)
	assert.Equal(t, 2, n)
	assert.Equal(t, "4567", buff.String())
}

func TestWriteString_TruncateWhenMaxCapIsExceeded(t *testing.T) {
	// GIVEN a LimitedBuffer with maximum capacity
	maxCap := 5
	buff := NewLimitedBuffer(maxCap)

	// WHEN adding data
	n, err := buff.WriteString("12")
	assert.NoError(t, err)
	assert.Equal(t, 2, n)

	// THEN buffer is correctly truncated
	n, err = buff.WriteString("3456")
	assert.NoError(t, err)
	assert.Equal(t, "23456", buff.String())
}
