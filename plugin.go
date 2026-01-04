package rrmcp

import (
	"context"
	"log"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// worker pool - pass a name and description to the worker
// logger
// configurer - what configuration options do we need?

const pluginName = "mcp"

type Plugin struct {
	server *mcp.Server
}

func (p *Plugin) Init() error {
	p.server = mcp.NewServer(&mcp.Implementation{Name: "greeter", Version: "v1.0.0"}, nil)
	return nil
}

func (p *Plugin) Serve() chan error {
	// Create a server with a single tool.
	mcp.AddTool(p.server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
	// Run the server over stdin/stdout, until the client disconnects.
	if err := p.server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (p *Plugin) Stop() error {
	return nil
}

func (p *Plugin) Name() string {
	return pluginName
}

func (p *Plugin) RPC() any {
	return &rpc{
		p:  p,
		mu: &sync.Mutex{},
	}
}