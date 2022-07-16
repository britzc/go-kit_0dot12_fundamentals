package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getGreeting(t *testing.T) {
	expected := "Welcome to Go kit 0.12 Fundamentals!"

	actual := getGreeting()

	assert.True(t, expected == actual, "~2|Test expected message: %s, not message %s~", expected, actual)
}
