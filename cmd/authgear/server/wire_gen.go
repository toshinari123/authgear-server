// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package server

import (
	"github.com/authgear/authgear-server/pkg/lib/config/source"
	"github.com/authgear/authgear-server/pkg/lib/deps"
	"github.com/authgear/authgear-server/pkg/lib/infra/task/executor"
	"github.com/authgear/authgear-server/pkg/lib/infra/task/queue"
)

// Injectors from wire.go:

func newConfigSource(p *deps.RootProvider) source.Source {
	serverConfig := p.ServerConfig
	factory := p.LoggerFactory
	localFileLogger := source.NewLocalFileLogger(factory)
	localFile := &source.LocalFile{
		Logger:       localFileLogger,
		ServerConfig: serverConfig,
	}
	sourceSource := source.NewSource(serverConfig, localFile)
	return sourceSource
}

func newInProcessQueue(p *deps.AppProvider, e *executor.InProcessExecutor) *queue.InProcessQueue {
	handle := p.Database
	config := p.Config
	captureTaskContext := deps.ProvideCaptureTaskContext(config)
	inProcessQueue := &queue.InProcessQueue{
		Database:       handle,
		CaptureContext: captureTaskContext,
		Executor:       e,
	}
	return inProcessQueue
}