package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var mH myHandler
	h := NoSurf(&mH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("type is not a Handler but it is of type %s", v)
	}
}

func TestSessionLoad(t *testing.T) {
	var mH myHandler
	h := sessionLoad(&mH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("type is not a Handler but it is of type %s", v)
	}
}
