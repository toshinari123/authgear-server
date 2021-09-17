// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package portal

import (
	"github.com/authgear/authgear-server/pkg/lib/admin/authz"
	"github.com/authgear/authgear-server/pkg/lib/analytic"
	"github.com/authgear/authgear-server/pkg/lib/infra/db/auditdb"
	"github.com/authgear/authgear-server/pkg/lib/infra/db/globaldb"
	"github.com/authgear/authgear-server/pkg/lib/infra/mail"
	"github.com/authgear/authgear-server/pkg/lib/infra/middleware"
	"github.com/authgear/authgear-server/pkg/portal/appresource/factory"
	"github.com/authgear/authgear-server/pkg/portal/deps"
	"github.com/authgear/authgear-server/pkg/portal/endpoint"
	"github.com/authgear/authgear-server/pkg/portal/graphql"
	"github.com/authgear/authgear-server/pkg/portal/lib/plan"
	"github.com/authgear/authgear-server/pkg/portal/loader"
	"github.com/authgear/authgear-server/pkg/portal/service"
	"github.com/authgear/authgear-server/pkg/portal/session"
	"github.com/authgear/authgear-server/pkg/portal/smtp"
	"github.com/authgear/authgear-server/pkg/portal/task"
	"github.com/authgear/authgear-server/pkg/portal/task/tasks"
	"github.com/authgear/authgear-server/pkg/portal/transport"
	"github.com/authgear/authgear-server/pkg/util/clock"
	"github.com/authgear/authgear-server/pkg/util/httproute"
	"github.com/authgear/authgear-server/pkg/util/intl"
	"github.com/authgear/authgear-server/pkg/util/template"
	"net/http"
)

import (
	_ "github.com/authgear/authgear-server/pkg/auth"
)

// Injectors from wire.go:

func newPanicEndMiddleware(p *deps.RequestProvider) httproute.Middleware {
	panicEndMiddleware := &middleware.PanicEndMiddleware{}
	return panicEndMiddleware
}

func newPanicLogMiddleware(p *deps.RequestProvider) httproute.Middleware {
	rootProvider := p.RootProvider
	factory := rootProvider.LoggerFactory
	panicLogMiddlewareLogger := middleware.NewPanicLogMiddlewareLogger(factory)
	panicLogMiddleware := &middleware.PanicLogMiddleware{
		Logger: panicLogMiddlewareLogger,
	}
	return panicLogMiddleware
}

func newPanicWriteEmptyResponseMiddleware(p *deps.RequestProvider) httproute.Middleware {
	panicWriteEmptyResponseMiddleware := &middleware.PanicWriteEmptyResponseMiddleware{}
	return panicWriteEmptyResponseMiddleware
}

func newBodyLimitMiddleware(p *deps.RequestProvider) httproute.Middleware {
	bodyLimitMiddleware := &middleware.BodyLimitMiddleware{}
	return bodyLimitMiddleware
}

func newSentryMiddleware(p *deps.RequestProvider) httproute.Middleware {
	rootProvider := p.RootProvider
	hub := rootProvider.SentryHub
	environmentConfig := rootProvider.EnvironmentConfig
	trustProxy := environmentConfig.TrustProxy
	sentryMiddleware := &middleware.SentryMiddleware{
		SentryHub:  hub,
		TrustProxy: trustProxy,
	}
	return sentryMiddleware
}

func newSessionInfoMiddleware(p *deps.RequestProvider) httproute.Middleware {
	sessionInfoMiddleware := &session.SessionInfoMiddleware{}
	return sessionInfoMiddleware
}

func newSessionRequiredMiddleware(p *deps.RequestProvider) httproute.Middleware {
	sessionRequiredMiddleware := &session.SessionRequiredMiddleware{}
	return sessionRequiredMiddleware
}

func newGraphQLHandler(p *deps.RequestProvider) http.Handler {
	rootProvider := p.RootProvider
	environmentConfig := rootProvider.EnvironmentConfig
	devMode := environmentConfig.DevMode
	logFactory := rootProvider.LoggerFactory
	logger := graphql.NewLogger(logFactory)
	authgearConfig := rootProvider.AuthgearConfig
	adminAPIConfig := rootProvider.AdminAPIConfig
	controller := rootProvider.ConfigSourceController
	configSource := deps.ProvideConfigSource(controller)
	clock := _wireSystemClockValue
	adder := &authz.Adder{
		Clock: clock,
	}
	adminAPIService := &service.AdminAPIService{
		AuthgearConfig: authgearConfig,
		AdminAPIConfig: adminAPIConfig,
		ConfigSource:   configSource,
		AuthzAdder:     adder,
	}
	userLoader := loader.NewUserLoader(adminAPIService)
	appServiceLogger := service.NewAppServiceLogger(logFactory)
	databaseEnvironmentConfig := rootProvider.DatabaseConfig
	sqlBuilder := globaldb.NewSQLBuilder(databaseEnvironmentConfig)
	request := p.Request
	context := deps.ProvideRequestContext(request)
	pool := rootProvider.Database
	handle := globaldb.NewHandle(context, pool, databaseEnvironmentConfig, logFactory)
	sqlExecutor := globaldb.NewSQLExecutor(context, handle)
	appConfig := rootProvider.AppConfig
	configServiceLogger := service.NewConfigServiceLogger(logFactory)
	domainImplementationType := rootProvider.DomainImplementation
	kubernetesConfig := rootProvider.KubernetesConfig
	kubernetesLogger := service.NewKubernetesLogger(logFactory)
	kubernetes := &service.Kubernetes{
		KubernetesConfig: kubernetesConfig,
		AppConfig:        appConfig,
		Logger:           kubernetesLogger,
	}
	configService := &service.ConfigService{
		Context:              context,
		Logger:               configServiceLogger,
		AppConfig:            appConfig,
		Controller:           controller,
		ConfigSource:         configSource,
		DomainImplementation: domainImplementationType,
		Kubernetes:           kubernetes,
	}
	mailConfig := rootProvider.MailConfig
	inProcessExecutorLogger := task.NewInProcessExecutorLogger(logFactory)
	mailLogger := mail.NewLogger(logFactory)
	smtpConfig := rootProvider.SMTPConfig
	smtpServerCredentials := deps.ProvideSMTPServerCredentials(smtpConfig)
	dialer := mail.NewGomailDialer(smtpServerCredentials)
	sender := &mail.Sender{
		Logger:       mailLogger,
		DevMode:      devMode,
		GomailDialer: dialer,
	}
	sendMessagesLogger := tasks.NewSendMessagesLogger(logFactory)
	sendMessagesTask := &tasks.SendMessagesTask{
		EmailSender: sender,
		Logger:      sendMessagesLogger,
	}
	inProcessExecutor := task.NewExecutor(inProcessExecutorLogger, sendMessagesTask)
	inProcessQueue := &task.InProcessQueue{
		Executor: inProcessExecutor,
	}
	trustProxy := environmentConfig.TrustProxy
	requestOriginProvider := &endpoint.RequestOriginProvider{
		Request:    request,
		TrustProxy: trustProxy,
	}
	endpointsProvider := &endpoint.EndpointsProvider{
		OriginProvider: requestOriginProvider,
	}
	manager := rootProvider.Resources
	defaultLanguageTag := _wireDefaultLanguageTagValue
	supportedLanguageTags := _wireSupportedLanguageTagsValue
	resolver := &template.Resolver{
		Resources:             manager,
		DefaultLanguageTag:    defaultLanguageTag,
		SupportedLanguageTags: supportedLanguageTags,
	}
	engine := &template.Engine{
		Resolver: resolver,
	}
	collaboratorService := &service.CollaboratorService{
		Context:        context,
		Clock:          clock,
		SQLBuilder:     sqlBuilder,
		SQLExecutor:    sqlExecutor,
		MailConfig:     mailConfig,
		TaskQueue:      inProcessQueue,
		Endpoints:      endpointsProvider,
		TemplateEngine: engine,
		AdminAPI:       adminAPIService,
	}
	authzService := &service.AuthzService{
		Context:       context,
		Configs:       configService,
		Collaborators: collaboratorService,
	}
	domainService := &service.DomainService{
		Context:      context,
		Clock:        clock,
		DomainConfig: configService,
		SQLBuilder:   sqlBuilder,
		SQLExecutor:  sqlExecutor,
	}
	appBaseResources := deps.ProvideAppBaseResources(rootProvider)
	managerFactory := &factory.ManagerFactory{
		AppBaseResources: appBaseResources,
	}
	store := &plan.Store{
		Clock:       clock,
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	planService := &plan.Service{
		PlanStore: store,
		AppConfig: appConfig,
	}
	appService := &service.AppService{
		Logger:           appServiceLogger,
		SQLBuilder:       sqlBuilder,
		SQLExecutor:      sqlExecutor,
		AppConfig:        appConfig,
		AppConfigs:       configService,
		AppAuthz:         authzService,
		AppAdminAPI:      adminAPIService,
		AppDomains:       domainService,
		Resources:        manager,
		AppResMgrFactory: managerFactory,
		Plan:             planService,
		Clock:            clock,
	}
	appLoader := loader.NewAppLoader(appService, authzService)
	domainLoader := loader.NewDomainLoader(domainService, authzService)
	collaboratorLoader := loader.NewCollaboratorLoader(collaboratorService, authzService)
	collaboratorInvitationLoader := loader.NewCollaboratorInvitationLoader(collaboratorService, authzService)
	smtpService := &smtp.Service{
		Context: context,
	}
	auditDatabaseCredentials := deps.ProvideAuditDatabaseCredentials(environmentConfig)
	auditdbSQLBuilder := auditdb.NewSQLBuilder(auditDatabaseCredentials)
	databaseConfig := deps.ProvideDatabaseConfig(databaseEnvironmentConfig)
	readHandle := auditdb.NewReadHandle(context, pool, databaseConfig, auditDatabaseCredentials, logFactory)
	readSQLExecutor := auditdb.NewReadSQLExecutor(context, readHandle)
	auditDBReadStore := &analytic.AuditDBReadStore{
		SQLBuilder:  auditdbSQLBuilder,
		SQLExecutor: readSQLExecutor,
	}
	chartService := &analytic.ChartService{
		AuditStore: auditDBReadStore,
	}
	graphqlContext := &graphql.Context{
		GQLLogger:               logger,
		Users:                   userLoader,
		Apps:                    appLoader,
		Domains:                 domainLoader,
		Collaborators:           collaboratorLoader,
		CollaboratorInvitations: collaboratorInvitationLoader,
		AuthzService:            authzService,
		AppService:              appService,
		DomainService:           domainService,
		CollaboratorService:     collaboratorService,
		SMTPService:             smtpService,
		AppResMgrFactory:        managerFactory,
		AnalyticChartService:    chartService,
	}
	graphQLHandler := &transport.GraphQLHandler{
		DevMode:        devMode,
		GraphQLContext: graphqlContext,
		Database:       handle,
		AuditDatabase:  readHandle,
	}
	return graphQLHandler
}

var (
	_wireSystemClockValue           = clock.NewSystemClock()
	_wireDefaultLanguageTagValue    = template.DefaultLanguageTag(intl.DefaultLanguage)
	_wireSupportedLanguageTagsValue = template.SupportedLanguageTags([]string{intl.DefaultLanguage})
)

func newSystemConfigHandler(p *deps.RequestProvider) http.Handler {
	rootProvider := p.RootProvider
	authgearConfig := rootProvider.AuthgearConfig
	appConfig := rootProvider.AppConfig
	searchConfig := rootProvider.SearchConfig
	auditLogConfig := rootProvider.AuditLogConfig
	manager := rootProvider.Resources
	systemConfigProvider := &service.SystemConfigProvider{
		AuthgearConfig: authgearConfig,
		AppConfig:      appConfig,
		SearchConfig:   searchConfig,
		AuditLogConfig: auditLogConfig,
		Resources:      manager,
	}
	systemConfigHandler := &transport.SystemConfigHandler{
		SystemConfig: systemConfigProvider,
	}
	return systemConfigHandler
}

func newAdminAPIHandler(p *deps.RequestProvider) http.Handler {
	request := p.Request
	context := deps.ProvideRequestContext(request)
	rootProvider := p.RootProvider
	pool := rootProvider.Database
	databaseEnvironmentConfig := rootProvider.DatabaseConfig
	logFactory := rootProvider.LoggerFactory
	handle := globaldb.NewHandle(context, pool, databaseEnvironmentConfig, logFactory)
	configServiceLogger := service.NewConfigServiceLogger(logFactory)
	appConfig := rootProvider.AppConfig
	controller := rootProvider.ConfigSourceController
	configSource := deps.ProvideConfigSource(controller)
	domainImplementationType := rootProvider.DomainImplementation
	kubernetesConfig := rootProvider.KubernetesConfig
	kubernetesLogger := service.NewKubernetesLogger(logFactory)
	kubernetes := &service.Kubernetes{
		KubernetesConfig: kubernetesConfig,
		AppConfig:        appConfig,
		Logger:           kubernetesLogger,
	}
	configService := &service.ConfigService{
		Context:              context,
		Logger:               configServiceLogger,
		AppConfig:            appConfig,
		Controller:           controller,
		ConfigSource:         configSource,
		DomainImplementation: domainImplementationType,
		Kubernetes:           kubernetes,
	}
	clockClock := _wireSystemClockValue
	sqlBuilder := globaldb.NewSQLBuilder(databaseEnvironmentConfig)
	sqlExecutor := globaldb.NewSQLExecutor(context, handle)
	mailConfig := rootProvider.MailConfig
	inProcessExecutorLogger := task.NewInProcessExecutorLogger(logFactory)
	logger := mail.NewLogger(logFactory)
	environmentConfig := rootProvider.EnvironmentConfig
	devMode := environmentConfig.DevMode
	smtpConfig := rootProvider.SMTPConfig
	smtpServerCredentials := deps.ProvideSMTPServerCredentials(smtpConfig)
	dialer := mail.NewGomailDialer(smtpServerCredentials)
	sender := &mail.Sender{
		Logger:       logger,
		DevMode:      devMode,
		GomailDialer: dialer,
	}
	sendMessagesLogger := tasks.NewSendMessagesLogger(logFactory)
	sendMessagesTask := &tasks.SendMessagesTask{
		EmailSender: sender,
		Logger:      sendMessagesLogger,
	}
	inProcessExecutor := task.NewExecutor(inProcessExecutorLogger, sendMessagesTask)
	inProcessQueue := &task.InProcessQueue{
		Executor: inProcessExecutor,
	}
	trustProxy := environmentConfig.TrustProxy
	requestOriginProvider := &endpoint.RequestOriginProvider{
		Request:    request,
		TrustProxy: trustProxy,
	}
	endpointsProvider := &endpoint.EndpointsProvider{
		OriginProvider: requestOriginProvider,
	}
	manager := rootProvider.Resources
	defaultLanguageTag := _wireDefaultLanguageTagValue
	supportedLanguageTags := _wireSupportedLanguageTagsValue
	resolver := &template.Resolver{
		Resources:             manager,
		DefaultLanguageTag:    defaultLanguageTag,
		SupportedLanguageTags: supportedLanguageTags,
	}
	engine := &template.Engine{
		Resolver: resolver,
	}
	authgearConfig := rootProvider.AuthgearConfig
	adminAPIConfig := rootProvider.AdminAPIConfig
	adder := &authz.Adder{
		Clock: clockClock,
	}
	adminAPIService := &service.AdminAPIService{
		AuthgearConfig: authgearConfig,
		AdminAPIConfig: adminAPIConfig,
		ConfigSource:   configSource,
		AuthzAdder:     adder,
	}
	collaboratorService := &service.CollaboratorService{
		Context:        context,
		Clock:          clockClock,
		SQLBuilder:     sqlBuilder,
		SQLExecutor:    sqlExecutor,
		MailConfig:     mailConfig,
		TaskQueue:      inProcessQueue,
		Endpoints:      endpointsProvider,
		TemplateEngine: engine,
		AdminAPI:       adminAPIService,
	}
	authzService := &service.AuthzService{
		Context:       context,
		Configs:       configService,
		Collaborators: collaboratorService,
	}
	adminAPILogger := transport.NewAdminAPILogger(logFactory)
	adminAPIHandler := &transport.AdminAPIHandler{
		Database: handle,
		Authz:    authzService,
		AdminAPI: adminAPIService,
		Logger:   adminAPILogger,
	}
	return adminAPIHandler
}

func newStaticAssetsHandler(p *deps.RequestProvider) http.Handler {
	rootProvider := p.RootProvider
	manager := rootProvider.Resources
	staticAssetsHandler := &transport.StaticAssetsHandler{
		Resources: manager,
	}
	return staticAssetsHandler
}
