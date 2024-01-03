package common

import "net/http"

type Router struct {
	Method string
	Pattern string
	Handler http.HandlerFunc
}