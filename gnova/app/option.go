package app

import (
	"Advanced_Shop/gnova/registry"
	"Advanced_Shop/gnova/server/rpcserver"
	"net/url"
	"os"
	"time"
)

// Option 特定签名的函数   这样返回多个options  就可以 送入 option...
type Option func(o *options)

type options struct {
	id        string
	endpoints []*url.URL
	name      string

	sigs []os.Signal

	//允许用户传入自己的实现
	registrar        registry.Registrar
	registrarTimeout time.Duration

	//stop超时时间
	stopTimeout time.Duration

	restServer *restserver.Server
	rpcServer  *rpcserver.Server
}

func WithRegistrar(registrar registry.Registrar) Option {
	return func(o *options) {
		o.registrar = registrar
	}
}

func WithEndpoints(endpoints []*url.URL) Option {
	return func(o *options) {
		o.endpoints = endpoints
	}
}

func WithRPCServer(server *rpcserver.Server) Option {
	return func(o *options) {
		o.rpcServer = server
	}
}

func WithRestServer(server *restserver.Server) Option {
	return func(o *options) {
		o.restServer = server
	}
}

func WithID(id string) Option {
	return func(o *options) {
		o.id = id
	}
}

func WithName(name string) Option {
	return func(o *options) {
		o.name = name
	}
}

func WithSigs(sigs []os.Signal) Option {
	return func(o *options) {
		o.sigs = sigs
	}
}
