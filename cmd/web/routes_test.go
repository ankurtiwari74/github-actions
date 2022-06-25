package main

import (
	"bookings/internals/config"
	"net/http"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app *config.AppConfig
	h := routes(app)

	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("TYpe is not a handler,  but its a %s", v)
	}
}
