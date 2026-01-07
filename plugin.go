package rrmcp

import (
	"context"
	"sync"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/roadrunner-server/endure/v2/dep"
	"github.com/roadrunner-server/errors"
	"github.com/roadrunner-server/mcp/v5/api"
	"go.uber.org/zap"
)

// worker pool - pass a name and description to the worker
// logger
// configurer - what configuration options do we need?
const (
	pluginName string = "mcp"
	RrMode     string = "RR_MODE"
)

type Plugin struct {
	server   *mcp.Server
	config   *config
	log      *zap.Logger
	mu       sync.RWMutex
	mcpTools map[string]api.MCPTool
}

func (p *Plugin) Init(cfg api.Configurer, log api.Logger, server api.Server) error {
	const op = errors.Op("mcp_plugin_init")

	if !cfg.Has(pluginName) {
		return errors.E(errors.Disabled)
	}

	err := cfg.UnmarshalKey(pluginName, &p.config)
	if err != nil {
		return errors.E(op, err)
	}

	err = p.config.InitDefaults()
	if err != nil {
		return errors.E(op, err)
	}

	p.log = log.NamedLogger(pluginName)
	p.server = mcp.NewServer(&mcp.Implementation{Name: "test", Version: "v1.0.0"}, nil)

	return nil
}

func (p *Plugin) Serve() chan error {
	errch := make(chan error, 1)

	p.mu.Lock()
	defer p.mu.Unlock()

	// add user defined tools
	for _, tool := range p.mcpTools {
		mcp.AddTool(p.server, &mcp.Tool{Name: tool.Name(), Description: tool.Description()}, tool.Tool)
	}

	// Create a server with a single tool.
	mcp.AddTool(p.server, &mcp.Tool{Name: "test", Description: "say hi"}, sendToWorker)
	// Run the server over stdin/stdout, until the client disconnects.
	if err := p.server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		p.log.Error("mcp server stopped with error", zap.Error(err))
		errch <- err
		return errch
	}

	return errch
}

func (p *Plugin) Stop() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return nil
}

func (p *Plugin) Name() string {
	return pluginName
}

func (p *Plugin) Collects() []*dep.In {
	return []*dep.In{
		dep.Fits(func(pp any) {
			p.mu.Lock()
			t := pp.(api.MCPTool)
			p.mcpTools[t.Name()] = t
			p.mu.Unlock()
		}, (*api.MCPTool)(nil)),
		dep.Fits(func(pp any) {
			p.mu.Lock()
			t := pp.(api.MCPTools)
			for _, tool := range t.Tools() {
				if _, exists := p.mcpTools[tool.Name()]; exists {
					p.log.Warn("mcp tool already added, skipping", zap.String("tool", tool.Name()))
					continue
				}
				p.mcpTools[tool.Name()] = tool
			}
			p.mu.Unlock()
		}, (*api.MCPTools)(nil)),
	}
}

func (p *Plugin) RPC() any {
	return &rpc{
		p:  p,
		mu: &sync.Mutex{},
	}
}
