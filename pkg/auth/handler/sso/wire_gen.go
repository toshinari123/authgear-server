// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package sso

import (
	"github.com/gorilla/mux"
	"github.com/skygeario/skygear-server/pkg/auth"
	auth2 "github.com/skygeario/skygear-server/pkg/auth/dependency/auth"
	redis3 "github.com/skygeario/skygear-server/pkg/auth/dependency/auth/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/bearertoken"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/oob"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/password"
	provider2 "github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/provider"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/recoverycode"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/totp"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/challenge"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/hook"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/identity/anonymous"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/identity/loginid"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/identity/oauth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/identity/provider"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/interaction"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/interaction/flows"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/interaction/redis"
	oauth2 "github.com/skygeario/skygear-server/pkg/auth/dependency/oauth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/handler"
	pq2 "github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/pq"
	redis2 "github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oidc"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/session"
	redis4 "github.com/skygeario/skygear-server/pkg/auth/dependency/session/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/sso"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/urlprefix"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/userprofile"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/welcomemessage"
	"github.com/skygeario/skygear-server/pkg/core/async"
	"github.com/skygeario/skygear-server/pkg/core/auth/authinfo/pq"
	"github.com/skygeario/skygear-server/pkg/core/config"
	"github.com/skygeario/skygear-server/pkg/core/db"
	handler2 "github.com/skygeario/skygear-server/pkg/core/handler"
	"github.com/skygeario/skygear-server/pkg/core/logging"
	"github.com/skygeario/skygear-server/pkg/core/time"
	"github.com/skygeario/skygear-server/pkg/core/validation"
	"net/http"
)

// Injectors from wire.go:

func newAuthHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	urlprefixProvider := urlprefix.NewProvider(r)
	authHandlerHTMLProvider := sso.ProvideAuthHandlerHTMLProvider(urlprefixProvider)
	ssoProvider := sso.ProvideSSOProvider(context, tenantConfiguration)
	timeProvider := time.NewProvider()
	normalizerFactory := loginid.ProvideNormalizerFactory(tenantConfiguration)
	redirectURLFunc := ProvideRedirectURIForAPIFunc()
	oAuthProviderFactory := sso.ProvideOAuthProviderFactory(tenantConfiguration, urlprefixProvider, timeProvider, normalizerFactory, redirectURLFunc)
	oAuthProvider := provideOAuthProviderFromRequestVars(r, oAuthProviderFactory)
	store := redis.ProvideStore(context, tenantConfiguration, timeProvider)
	requestID := auth.ProvideLoggingRequestID(r)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	typeCheckerFactory := loginid.ProvideTypeCheckerFactory(tenantConfiguration, reservedNameChecker)
	checker := loginid.ProvideChecker(tenantConfiguration, typeCheckerFactory)
	loginidProvider := loginid.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration, checker, normalizerFactory)
	oauthProvider := oauth.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider)
	anonymousProvider := anonymous.ProvideProvider(sqlBuilder, sqlExecutor)
	providerProvider := provider.ProvideProvider(tenantConfiguration, loginidProvider, oauthProvider, anonymousProvider)
	historyStoreImpl := password.ProvideHistoryStore(timeProvider, sqlBuilder, sqlExecutor)
	passwordChecker := password.ProvideChecker(tenantConfiguration, historyStoreImpl)
	passwordProvider := password.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, factory, historyStoreImpl, passwordChecker, tenantConfiguration)
	totpProvider := totp.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, txContext, requestID, tenantConfiguration, executor)
	oobProvider := oob.ProvideProvider(tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider, engine, urlprefixProvider, queue)
	bearertokenProvider := bearertoken.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	recoverycodeProvider := recoverycode.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	provider3 := &provider2.Provider{
		Password:     passwordProvider,
		TOTP:         totpProvider,
		OOBOTP:       oobProvider,
		BearerToken:  bearertokenProvider,
		RecoveryCode: recoverycodeProvider,
	}
	authinfoStore := pq.ProvideStore(sqlBuilderFactory, sqlExecutor)
	userprofileStore := userprofile.ProvideStore(timeProvider, sqlBuilder, sqlExecutor)
	hookProvider := hook.ProvideHookProvider(context, sqlBuilder, sqlExecutor, requestID, tenantConfiguration, txContext, timeProvider, authinfoStore, userprofileStore, loginidProvider, factory)
	welcomemessageProvider := welcomemessage.ProvideProvider(tenantConfiguration, engine, queue)
	welcomeMessageProvider := interaction.ProvideWelcomeMessageProvider(welcomemessageProvider)
	userProvider := interaction.ProvideUserProvider(authinfoStore, userprofileStore, timeProvider, hookProvider, urlprefixProvider, queue, tenantConfiguration, welcomeMessageProvider)
	interactionProvider := interaction.ProvideProvider(store, timeProvider, factory, providerProvider, provider3, userProvider, oobProvider, tenantConfiguration, hookProvider)
	authorizationStore := &pq2.AuthorizationStore{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	grantStore := redis2.ProvideGrantStore(context, factory, tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider)
	eventStore := redis3.ProvideEventStore(context, tenantConfiguration)
	accessEventProvider := auth2.AccessEventProvider{
		Store: eventStore,
	}
	sessionStore := redis4.ProvideStore(context, tenantConfiguration, timeProvider, factory)
	authAccessEventProvider := &auth2.AccessEventProvider{
		Store: eventStore,
	}
	sessionProvider := session.ProvideSessionProvider(r, sessionStore, authAccessEventProvider, tenantConfiguration)
	challengeProvider := challenge.ProvideProvider(context, timeProvider, tenantConfiguration)
	anonymousFlow := &flows.AnonymousFlow{
		Interactions: interactionProvider,
		Anonymous:    anonymousProvider,
		Challenges:   challengeProvider,
	}
	idTokenIssuer := oidc.ProvideIDTokenIssuer(tenantConfiguration, urlprefixProvider, authinfoStore, userprofileStore, timeProvider)
	tokenGenerator := _wireTokenGeneratorValue
	tokenHandler := handler.ProvideTokenHandler(r, tenantConfiguration, factory, authorizationStore, grantStore, grantStore, grantStore, accessEventProvider, sessionProvider, anonymousFlow, idTokenIssuer, tokenGenerator, timeProvider)
	insecureCookieConfig := auth.ProvideSessionInsecureCookieConfig(m)
	cookieConfiguration := session.ProvideSessionCookieConfiguration(r, insecureCookieConfig, tenantConfiguration)
	userController := flows.ProvideUserController(authinfoStore, userprofileStore, tokenHandler, cookieConfiguration, sessionProvider, hookProvider, timeProvider, tenantConfiguration)
	authAPIFlow := &flows.AuthAPIFlow{
		Interactions:   interactionProvider,
		UserController: userController,
	}
	httpHandler := provideAuthHandler(txContext, tenantConfiguration, authHandlerHTMLProvider, ssoProvider, oAuthProvider, authAPIFlow)
	return httpHandler
}

var (
	_wireTokenGeneratorValue = handler.TokenGenerator(oauth2.GenerateToken)
)

func newAuthResultHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	requestID := auth.ProvideLoggingRequestID(r)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	requireAuthz := handler2.NewRequireAuthzFactory(factory)
	validator := auth.ProvideValidator(m)
	ssoProvider := sso.ProvideSSOProvider(context, tenantConfiguration)
	timeProvider := time.NewProvider()
	store := redis.ProvideStore(context, tenantConfiguration, timeProvider)
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	typeCheckerFactory := loginid.ProvideTypeCheckerFactory(tenantConfiguration, reservedNameChecker)
	checker := loginid.ProvideChecker(tenantConfiguration, typeCheckerFactory)
	normalizerFactory := loginid.ProvideNormalizerFactory(tenantConfiguration)
	loginidProvider := loginid.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration, checker, normalizerFactory)
	oauthProvider := oauth.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider)
	anonymousProvider := anonymous.ProvideProvider(sqlBuilder, sqlExecutor)
	providerProvider := provider.ProvideProvider(tenantConfiguration, loginidProvider, oauthProvider, anonymousProvider)
	historyStoreImpl := password.ProvideHistoryStore(timeProvider, sqlBuilder, sqlExecutor)
	passwordChecker := password.ProvideChecker(tenantConfiguration, historyStoreImpl)
	passwordProvider := password.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, factory, historyStoreImpl, passwordChecker, tenantConfiguration)
	totpProvider := totp.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	urlprefixProvider := urlprefix.NewProvider(r)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, txContext, requestID, tenantConfiguration, executor)
	oobProvider := oob.ProvideProvider(tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider, engine, urlprefixProvider, queue)
	bearertokenProvider := bearertoken.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	recoverycodeProvider := recoverycode.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	provider3 := &provider2.Provider{
		Password:     passwordProvider,
		TOTP:         totpProvider,
		OOBOTP:       oobProvider,
		BearerToken:  bearertokenProvider,
		RecoveryCode: recoverycodeProvider,
	}
	authinfoStore := pq.ProvideStore(sqlBuilderFactory, sqlExecutor)
	userprofileStore := userprofile.ProvideStore(timeProvider, sqlBuilder, sqlExecutor)
	hookProvider := hook.ProvideHookProvider(context, sqlBuilder, sqlExecutor, requestID, tenantConfiguration, txContext, timeProvider, authinfoStore, userprofileStore, loginidProvider, factory)
	welcomemessageProvider := welcomemessage.ProvideProvider(tenantConfiguration, engine, queue)
	welcomeMessageProvider := interaction.ProvideWelcomeMessageProvider(welcomemessageProvider)
	userProvider := interaction.ProvideUserProvider(authinfoStore, userprofileStore, timeProvider, hookProvider, urlprefixProvider, queue, tenantConfiguration, welcomeMessageProvider)
	interactionProvider := interaction.ProvideProvider(store, timeProvider, factory, providerProvider, provider3, userProvider, oobProvider, tenantConfiguration, hookProvider)
	authorizationStore := &pq2.AuthorizationStore{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	grantStore := redis2.ProvideGrantStore(context, factory, tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider)
	eventStore := redis3.ProvideEventStore(context, tenantConfiguration)
	accessEventProvider := auth2.AccessEventProvider{
		Store: eventStore,
	}
	sessionStore := redis4.ProvideStore(context, tenantConfiguration, timeProvider, factory)
	authAccessEventProvider := &auth2.AccessEventProvider{
		Store: eventStore,
	}
	sessionProvider := session.ProvideSessionProvider(r, sessionStore, authAccessEventProvider, tenantConfiguration)
	challengeProvider := challenge.ProvideProvider(context, timeProvider, tenantConfiguration)
	anonymousFlow := &flows.AnonymousFlow{
		Interactions: interactionProvider,
		Anonymous:    anonymousProvider,
		Challenges:   challengeProvider,
	}
	idTokenIssuer := oidc.ProvideIDTokenIssuer(tenantConfiguration, urlprefixProvider, authinfoStore, userprofileStore, timeProvider)
	tokenGenerator := _wireTokenGeneratorValue
	tokenHandler := handler.ProvideTokenHandler(r, tenantConfiguration, factory, authorizationStore, grantStore, grantStore, grantStore, accessEventProvider, sessionProvider, anonymousFlow, idTokenIssuer, tokenGenerator, timeProvider)
	insecureCookieConfig := auth.ProvideSessionInsecureCookieConfig(m)
	cookieConfiguration := session.ProvideSessionCookieConfiguration(r, insecureCookieConfig, tenantConfiguration)
	userController := flows.ProvideUserController(authinfoStore, userprofileStore, tokenHandler, cookieConfiguration, sessionProvider, hookProvider, timeProvider, tenantConfiguration)
	authAPIFlow := &flows.AuthAPIFlow{
		Interactions:   interactionProvider,
		UserController: userController,
	}
	httpHandler := provideAuthResultHandler(txContext, requireAuthz, validator, ssoProvider, authAPIFlow)
	return httpHandler
}

func newLinkHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	requestID := auth.ProvideLoggingRequestID(r)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	requireAuthz := handler2.NewRequireAuthzFactory(factory)
	validator := auth.ProvideValidator(m)
	ssoProvider := sso.ProvideSSOProvider(context, tenantConfiguration)
	urlprefixProvider := urlprefix.NewProvider(r)
	timeProvider := time.NewProvider()
	normalizerFactory := loginid.ProvideNormalizerFactory(tenantConfiguration)
	redirectURLFunc := ProvideRedirectURIForAPIFunc()
	oAuthProviderFactory := sso.ProvideOAuthProviderFactory(tenantConfiguration, urlprefixProvider, timeProvider, normalizerFactory, redirectURLFunc)
	oAuthProvider := provideOAuthProviderFromRequestVars(r, oAuthProviderFactory)
	store := redis.ProvideStore(context, tenantConfiguration, timeProvider)
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	typeCheckerFactory := loginid.ProvideTypeCheckerFactory(tenantConfiguration, reservedNameChecker)
	checker := loginid.ProvideChecker(tenantConfiguration, typeCheckerFactory)
	loginidProvider := loginid.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration, checker, normalizerFactory)
	oauthProvider := oauth.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider)
	anonymousProvider := anonymous.ProvideProvider(sqlBuilder, sqlExecutor)
	providerProvider := provider.ProvideProvider(tenantConfiguration, loginidProvider, oauthProvider, anonymousProvider)
	historyStoreImpl := password.ProvideHistoryStore(timeProvider, sqlBuilder, sqlExecutor)
	passwordChecker := password.ProvideChecker(tenantConfiguration, historyStoreImpl)
	passwordProvider := password.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, factory, historyStoreImpl, passwordChecker, tenantConfiguration)
	totpProvider := totp.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, txContext, requestID, tenantConfiguration, executor)
	oobProvider := oob.ProvideProvider(tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider, engine, urlprefixProvider, queue)
	bearertokenProvider := bearertoken.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	recoverycodeProvider := recoverycode.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	provider3 := &provider2.Provider{
		Password:     passwordProvider,
		TOTP:         totpProvider,
		OOBOTP:       oobProvider,
		BearerToken:  bearertokenProvider,
		RecoveryCode: recoverycodeProvider,
	}
	authinfoStore := pq.ProvideStore(sqlBuilderFactory, sqlExecutor)
	userprofileStore := userprofile.ProvideStore(timeProvider, sqlBuilder, sqlExecutor)
	hookProvider := hook.ProvideHookProvider(context, sqlBuilder, sqlExecutor, requestID, tenantConfiguration, txContext, timeProvider, authinfoStore, userprofileStore, loginidProvider, factory)
	welcomemessageProvider := welcomemessage.ProvideProvider(tenantConfiguration, engine, queue)
	welcomeMessageProvider := interaction.ProvideWelcomeMessageProvider(welcomemessageProvider)
	userProvider := interaction.ProvideUserProvider(authinfoStore, userprofileStore, timeProvider, hookProvider, urlprefixProvider, queue, tenantConfiguration, welcomeMessageProvider)
	interactionProvider := interaction.ProvideProvider(store, timeProvider, factory, providerProvider, provider3, userProvider, oobProvider, tenantConfiguration, hookProvider)
	authorizationStore := &pq2.AuthorizationStore{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	grantStore := redis2.ProvideGrantStore(context, factory, tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider)
	eventStore := redis3.ProvideEventStore(context, tenantConfiguration)
	accessEventProvider := auth2.AccessEventProvider{
		Store: eventStore,
	}
	sessionStore := redis4.ProvideStore(context, tenantConfiguration, timeProvider, factory)
	authAccessEventProvider := &auth2.AccessEventProvider{
		Store: eventStore,
	}
	sessionProvider := session.ProvideSessionProvider(r, sessionStore, authAccessEventProvider, tenantConfiguration)
	challengeProvider := challenge.ProvideProvider(context, timeProvider, tenantConfiguration)
	anonymousFlow := &flows.AnonymousFlow{
		Interactions: interactionProvider,
		Anonymous:    anonymousProvider,
		Challenges:   challengeProvider,
	}
	idTokenIssuer := oidc.ProvideIDTokenIssuer(tenantConfiguration, urlprefixProvider, authinfoStore, userprofileStore, timeProvider)
	tokenGenerator := _wireTokenGeneratorValue
	tokenHandler := handler.ProvideTokenHandler(r, tenantConfiguration, factory, authorizationStore, grantStore, grantStore, grantStore, accessEventProvider, sessionProvider, anonymousFlow, idTokenIssuer, tokenGenerator, timeProvider)
	insecureCookieConfig := auth.ProvideSessionInsecureCookieConfig(m)
	cookieConfiguration := session.ProvideSessionCookieConfiguration(r, insecureCookieConfig, tenantConfiguration)
	userController := flows.ProvideUserController(authinfoStore, userprofileStore, tokenHandler, cookieConfiguration, sessionProvider, hookProvider, timeProvider, tenantConfiguration)
	authAPIFlow := &flows.AuthAPIFlow{
		Interactions:   interactionProvider,
		UserController: userController,
	}
	httpHandler := provideLinkHandler(txContext, requireAuthz, validator, ssoProvider, oAuthProvider, authAPIFlow)
	return httpHandler
}

func newLoginHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	requestID := auth.ProvideLoggingRequestID(r)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	requireAuthz := handler2.NewRequireAuthzFactory(factory)
	validator := auth.ProvideValidator(m)
	ssoProvider := sso.ProvideSSOProvider(context, tenantConfiguration)
	timeProvider := time.NewProvider()
	store := redis.ProvideStore(context, tenantConfiguration, timeProvider)
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	typeCheckerFactory := loginid.ProvideTypeCheckerFactory(tenantConfiguration, reservedNameChecker)
	checker := loginid.ProvideChecker(tenantConfiguration, typeCheckerFactory)
	normalizerFactory := loginid.ProvideNormalizerFactory(tenantConfiguration)
	loginidProvider := loginid.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration, checker, normalizerFactory)
	oauthProvider := oauth.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider)
	anonymousProvider := anonymous.ProvideProvider(sqlBuilder, sqlExecutor)
	providerProvider := provider.ProvideProvider(tenantConfiguration, loginidProvider, oauthProvider, anonymousProvider)
	historyStoreImpl := password.ProvideHistoryStore(timeProvider, sqlBuilder, sqlExecutor)
	passwordChecker := password.ProvideChecker(tenantConfiguration, historyStoreImpl)
	passwordProvider := password.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, factory, historyStoreImpl, passwordChecker, tenantConfiguration)
	totpProvider := totp.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	urlprefixProvider := urlprefix.NewProvider(r)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, txContext, requestID, tenantConfiguration, executor)
	oobProvider := oob.ProvideProvider(tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider, engine, urlprefixProvider, queue)
	bearertokenProvider := bearertoken.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	recoverycodeProvider := recoverycode.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	provider3 := &provider2.Provider{
		Password:     passwordProvider,
		TOTP:         totpProvider,
		OOBOTP:       oobProvider,
		BearerToken:  bearertokenProvider,
		RecoveryCode: recoverycodeProvider,
	}
	authinfoStore := pq.ProvideStore(sqlBuilderFactory, sqlExecutor)
	userprofileStore := userprofile.ProvideStore(timeProvider, sqlBuilder, sqlExecutor)
	hookProvider := hook.ProvideHookProvider(context, sqlBuilder, sqlExecutor, requestID, tenantConfiguration, txContext, timeProvider, authinfoStore, userprofileStore, loginidProvider, factory)
	welcomemessageProvider := welcomemessage.ProvideProvider(tenantConfiguration, engine, queue)
	welcomeMessageProvider := interaction.ProvideWelcomeMessageProvider(welcomemessageProvider)
	userProvider := interaction.ProvideUserProvider(authinfoStore, userprofileStore, timeProvider, hookProvider, urlprefixProvider, queue, tenantConfiguration, welcomeMessageProvider)
	interactionProvider := interaction.ProvideProvider(store, timeProvider, factory, providerProvider, provider3, userProvider, oobProvider, tenantConfiguration, hookProvider)
	authorizationStore := &pq2.AuthorizationStore{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	grantStore := redis2.ProvideGrantStore(context, factory, tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider)
	eventStore := redis3.ProvideEventStore(context, tenantConfiguration)
	accessEventProvider := auth2.AccessEventProvider{
		Store: eventStore,
	}
	sessionStore := redis4.ProvideStore(context, tenantConfiguration, timeProvider, factory)
	authAccessEventProvider := &auth2.AccessEventProvider{
		Store: eventStore,
	}
	sessionProvider := session.ProvideSessionProvider(r, sessionStore, authAccessEventProvider, tenantConfiguration)
	challengeProvider := challenge.ProvideProvider(context, timeProvider, tenantConfiguration)
	anonymousFlow := &flows.AnonymousFlow{
		Interactions: interactionProvider,
		Anonymous:    anonymousProvider,
		Challenges:   challengeProvider,
	}
	idTokenIssuer := oidc.ProvideIDTokenIssuer(tenantConfiguration, urlprefixProvider, authinfoStore, userprofileStore, timeProvider)
	tokenGenerator := _wireTokenGeneratorValue
	tokenHandler := handler.ProvideTokenHandler(r, tenantConfiguration, factory, authorizationStore, grantStore, grantStore, grantStore, accessEventProvider, sessionProvider, anonymousFlow, idTokenIssuer, tokenGenerator, timeProvider)
	insecureCookieConfig := auth.ProvideSessionInsecureCookieConfig(m)
	cookieConfiguration := session.ProvideSessionCookieConfiguration(r, insecureCookieConfig, tenantConfiguration)
	userController := flows.ProvideUserController(authinfoStore, userprofileStore, tokenHandler, cookieConfiguration, sessionProvider, hookProvider, timeProvider, tenantConfiguration)
	authAPIFlow := &flows.AuthAPIFlow{
		Interactions:   interactionProvider,
		UserController: userController,
	}
	redirectURLFunc := ProvideRedirectURIForAPIFunc()
	oAuthProviderFactory := sso.ProvideOAuthProviderFactory(tenantConfiguration, urlprefixProvider, timeProvider, normalizerFactory, redirectURLFunc)
	oAuthProvider := provideOAuthProviderFromRequestVars(r, oAuthProviderFactory)
	httpHandler := provideLoginHandler(txContext, requireAuthz, validator, ssoProvider, authAPIFlow, oAuthProvider)
	return httpHandler
}

func newAuthRedirectHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	ssoProvider := sso.ProvideSSOProvider(context, tenantConfiguration)
	urlprefixProvider := urlprefix.NewProvider(r)
	timeProvider := time.NewProvider()
	normalizerFactory := loginid.ProvideNormalizerFactory(tenantConfiguration)
	redirectURLFunc := ProvideRedirectURIForAPIFunc()
	oAuthProviderFactory := sso.ProvideOAuthProviderFactory(tenantConfiguration, urlprefixProvider, timeProvider, normalizerFactory, redirectURLFunc)
	oAuthProvider := provideOAuthProviderFromRequestVars(r, oAuthProviderFactory)
	httpHandler := provideAuthRedirectHandler(ssoProvider, oAuthProvider)
	return httpHandler
}

func newLoginAuthURLHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	requestID := auth.ProvideLoggingRequestID(r)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	requireAuthz := handler2.NewRequireAuthzFactory(factory)
	validator := auth.ProvideValidator(m)
	ssoProvider := sso.ProvideSSOProvider(context, tenantConfiguration)
	urlprefixProvider := urlprefix.NewProvider(r)
	timeProvider := time.NewProvider()
	normalizerFactory := loginid.ProvideNormalizerFactory(tenantConfiguration)
	redirectURLFunc := ProvideRedirectURIForAPIFunc()
	oAuthProviderFactory := sso.ProvideOAuthProviderFactory(tenantConfiguration, urlprefixProvider, timeProvider, normalizerFactory, redirectURLFunc)
	oAuthProvider := provideOAuthProviderFromRequestVars(r, oAuthProviderFactory)
	ssoSsoAction := providerLoginSSOAction()
	httpHandler := provideAuthURLHandler(txContext, requireAuthz, validator, ssoProvider, tenantConfiguration, oAuthProvider, ssoSsoAction)
	return httpHandler
}

func newLinkAuthURLHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	requestID := auth.ProvideLoggingRequestID(r)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	requireAuthz := handler2.NewRequireAuthzFactory(factory)
	validator := auth.ProvideValidator(m)
	ssoProvider := sso.ProvideSSOProvider(context, tenantConfiguration)
	urlprefixProvider := urlprefix.NewProvider(r)
	timeProvider := time.NewProvider()
	normalizerFactory := loginid.ProvideNormalizerFactory(tenantConfiguration)
	redirectURLFunc := ProvideRedirectURIForAPIFunc()
	oAuthProviderFactory := sso.ProvideOAuthProviderFactory(tenantConfiguration, urlprefixProvider, timeProvider, normalizerFactory, redirectURLFunc)
	oAuthProvider := provideOAuthProviderFromRequestVars(r, oAuthProviderFactory)
	ssoSsoAction := providerLinkSSOAction()
	httpHandler := provideAuthURLHandler(txContext, requireAuthz, validator, ssoProvider, tenantConfiguration, oAuthProvider, ssoSsoAction)
	return httpHandler
}

func newUnlinkHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	tenantConfiguration := auth.ProvideTenantConfig(context, m)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	requestID := auth.ProvideLoggingRequestID(r)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	requireAuthz := handler2.NewRequireAuthzFactory(factory)
	urlprefixProvider := urlprefix.NewProvider(r)
	timeProvider := time.NewProvider()
	normalizerFactory := loginid.ProvideNormalizerFactory(tenantConfiguration)
	redirectURLFunc := ProvideRedirectURIForAPIFunc()
	oAuthProviderFactory := sso.ProvideOAuthProviderFactory(tenantConfiguration, urlprefixProvider, timeProvider, normalizerFactory, redirectURLFunc)
	store := redis.ProvideStore(context, tenantConfiguration, timeProvider)
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	typeCheckerFactory := loginid.ProvideTypeCheckerFactory(tenantConfiguration, reservedNameChecker)
	checker := loginid.ProvideChecker(tenantConfiguration, typeCheckerFactory)
	loginidProvider := loginid.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration, checker, normalizerFactory)
	oauthProvider := oauth.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider)
	anonymousProvider := anonymous.ProvideProvider(sqlBuilder, sqlExecutor)
	providerProvider := provider.ProvideProvider(tenantConfiguration, loginidProvider, oauthProvider, anonymousProvider)
	historyStoreImpl := password.ProvideHistoryStore(timeProvider, sqlBuilder, sqlExecutor)
	passwordChecker := password.ProvideChecker(tenantConfiguration, historyStoreImpl)
	passwordProvider := password.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, factory, historyStoreImpl, passwordChecker, tenantConfiguration)
	totpProvider := totp.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	executor := auth.ProvideTaskExecutor(m)
	queue := async.ProvideTaskQueue(context, txContext, requestID, tenantConfiguration, executor)
	oobProvider := oob.ProvideProvider(tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider, engine, urlprefixProvider, queue)
	bearertokenProvider := bearertoken.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	recoverycodeProvider := recoverycode.ProvideProvider(sqlBuilder, sqlExecutor, timeProvider, tenantConfiguration)
	provider3 := &provider2.Provider{
		Password:     passwordProvider,
		TOTP:         totpProvider,
		OOBOTP:       oobProvider,
		BearerToken:  bearertokenProvider,
		RecoveryCode: recoverycodeProvider,
	}
	authinfoStore := pq.ProvideStore(sqlBuilderFactory, sqlExecutor)
	userprofileStore := userprofile.ProvideStore(timeProvider, sqlBuilder, sqlExecutor)
	hookProvider := hook.ProvideHookProvider(context, sqlBuilder, sqlExecutor, requestID, tenantConfiguration, txContext, timeProvider, authinfoStore, userprofileStore, loginidProvider, factory)
	welcomemessageProvider := welcomemessage.ProvideProvider(tenantConfiguration, engine, queue)
	welcomeMessageProvider := interaction.ProvideWelcomeMessageProvider(welcomemessageProvider)
	userProvider := interaction.ProvideUserProvider(authinfoStore, userprofileStore, timeProvider, hookProvider, urlprefixProvider, queue, tenantConfiguration, welcomeMessageProvider)
	interactionProvider := interaction.ProvideProvider(store, timeProvider, factory, providerProvider, provider3, userProvider, oobProvider, tenantConfiguration, hookProvider)
	authorizationStore := &pq2.AuthorizationStore{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	grantStore := redis2.ProvideGrantStore(context, factory, tenantConfiguration, sqlBuilder, sqlExecutor, timeProvider)
	eventStore := redis3.ProvideEventStore(context, tenantConfiguration)
	accessEventProvider := auth2.AccessEventProvider{
		Store: eventStore,
	}
	sessionStore := redis4.ProvideStore(context, tenantConfiguration, timeProvider, factory)
	authAccessEventProvider := &auth2.AccessEventProvider{
		Store: eventStore,
	}
	sessionProvider := session.ProvideSessionProvider(r, sessionStore, authAccessEventProvider, tenantConfiguration)
	challengeProvider := challenge.ProvideProvider(context, timeProvider, tenantConfiguration)
	anonymousFlow := &flows.AnonymousFlow{
		Interactions: interactionProvider,
		Anonymous:    anonymousProvider,
		Challenges:   challengeProvider,
	}
	idTokenIssuer := oidc.ProvideIDTokenIssuer(tenantConfiguration, urlprefixProvider, authinfoStore, userprofileStore, timeProvider)
	tokenGenerator := _wireTokenGeneratorValue
	tokenHandler := handler.ProvideTokenHandler(r, tenantConfiguration, factory, authorizationStore, grantStore, grantStore, grantStore, accessEventProvider, sessionProvider, anonymousFlow, idTokenIssuer, tokenGenerator, timeProvider)
	insecureCookieConfig := auth.ProvideSessionInsecureCookieConfig(m)
	cookieConfiguration := session.ProvideSessionCookieConfiguration(r, insecureCookieConfig, tenantConfiguration)
	userController := flows.ProvideUserController(authinfoStore, userprofileStore, tokenHandler, cookieConfiguration, sessionProvider, hookProvider, timeProvider, tenantConfiguration)
	authAPIFlow := &flows.AuthAPIFlow{
		Interactions:   interactionProvider,
		UserController: userController,
	}
	httpHandler := providerUnlinkHandler(txContext, requireAuthz, oAuthProviderFactory, authAPIFlow)
	return httpHandler
}

// wire.go:

func provideOAuthProviderFromRequestVars(r *http.Request, spf *sso.OAuthProviderFactory) sso.OAuthProvider {
	vars := mux.Vars(r)
	return spf.NewOAuthProvider(vars["provider"])
}

func ProvideRedirectURIForAPIFunc() sso.RedirectURLFunc {
	return RedirectURIForAPI
}

func provideAuthHandler(
	tx db.TxContext,
	cfg *config.TenantConfiguration,
	hp sso.AuthHandlerHTMLProvider,
	sp sso.Provider,
	op sso.OAuthProvider,
	f OAuthHandlerInteractionFlow,
) http.Handler {
	h := &AuthHandler{
		TxContext:               tx,
		TenantConfiguration:     cfg,
		AuthHandlerHTMLProvider: hp,
		SSOProvider:             sp,
		OAuthProvider:           op,
		Interactions:            f,
	}
	return h
}

func provideAuthResultHandler(
	tx db.TxContext,
	requireAuthz handler2.RequireAuthz,
	v *validation.Validator,
	sp sso.Provider,
	f OAuthResultInteractionFlow,
) http.Handler {
	h := &AuthResultHandler{
		TxContext:    tx,
		Validator:    v,
		SSOProvider:  sp,
		Interactions: f,
	}
	return requireAuthz(h, h)
}

func provideLinkHandler(
	tx db.TxContext,
	requireAuthz handler2.RequireAuthz,
	v *validation.Validator,
	sp sso.Provider,
	op sso.OAuthProvider,
	f OAuthLinkInteractionFlow,
) http.Handler {
	h := &LinkHandler{
		TxContext:     tx,
		Validator:     v,
		SSOProvider:   sp,
		OAuthProvider: op,
		Interactions:  f,
	}
	return requireAuthz(h, h)
}

func provideLoginHandler(
	tx db.TxContext,
	requireAuthz handler2.RequireAuthz,
	v *validation.Validator,
	sp sso.Provider,
	f OAuthLoginInteractionFlow,
	op sso.OAuthProvider,
) http.Handler {
	h := &LoginHandler{
		TxContext:     tx,
		Validator:     v,
		SSOProvider:   sp,
		OAuthProvider: op,
		Interactions:  f,
	}
	return requireAuthz(h, h)
}

func provideAuthRedirectHandler(
	sp sso.Provider,
	op sso.OAuthProvider,
) http.Handler {
	h := &AuthRedirectHandler{
		SSOProvider:   sp,
		OAuthProvider: op,
	}
	return h
}

func provideAuthURLHandler(
	tx db.TxContext,
	requireAuthz handler2.RequireAuthz,
	v *validation.Validator,
	sp sso.Provider,
	cfg *config.TenantConfiguration,
	op sso.OAuthProvider,
	action ssoAction,
) http.Handler {
	h := &AuthURLHandler{
		TxContext:                  tx,
		Validator:                  v,
		SSOProvider:                sp,
		OAuthConflictConfiguration: cfg.AppConfig.AuthAPI.OnIdentityConflict.OAuth,
		OAuthProvider:              op,
		Action:                     action,
	}
	return requireAuthz(h, h)
}

func providerLoginSSOAction() ssoAction {
	return ssoActionLogin
}

func providerLinkSSOAction() ssoAction {
	return ssoActionLink
}

func providerUnlinkHandler(
	tx db.TxContext,
	requireAuthz handler2.RequireAuthz,
	spf *sso.OAuthProviderFactory,
	f OAuthUnlinkInteractionFlow,
) http.Handler {
	h := &UnlinkHandler{
		TxContext:       tx,
		ProviderFactory: spf,
		Interactions:    f,
	}
	return requireAuthz(h, h)
}
