package rrmcp

import "sync"

type rpc struct {
	p  *Plugin
	mu *sync.Mutex
}

func (r *rpc) Register() error {
	return nil
}

// do we need to explicitly start/stop the server?
func (r *rpc) StartServer() error {
	return nil
}

// do we need to explicitly start/stop the server?
func (r *rpc) StopServer() error {
	return nil
}
