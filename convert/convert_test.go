// +build interactive

package main

import (
	"testing"
)

// having these as tests makes it easy to run them from editors/IDE's
// the build tag is so it will be skipped by go test ./...

func TestConvert(*testing.T) {
	main()
}

func TestRead(*testing.T) {
	Read()
}
