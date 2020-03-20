// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package oauth

import (
	"github.com/skygeario/skygear-server/pkg/auth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oauth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/handler"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/pq"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oauth/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/oidc"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/session"
	redis2 "github.com/skygeario/skygear-server/pkg/auth/dependency/session/redis"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/urlprefix"
	"github.com/skygeario/skygear-server/pkg/core/db"
	"github.com/skygeario/skygear-server/pkg/core/logging"
	"github.com/skygeario/skygear-server/pkg/core/time"
	"net/http"
)

// Injectors from wire.go:

func newAuthorizeHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	requestID := auth.ProvideLoggingRequestID(r)
	tenantConfiguration := auth.ProvideTenantConfig(context)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	authorizationStore := &pq.AuthorizationStore{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	provider := time.NewProvider()
	grantStore := redis.ProvideGrantStore(context, tenantConfiguration, sqlBuilder, sqlExecutor, provider)
	urlprefixProvider := urlprefix.NewProvider(r)
	endpointsProvider := &auth.EndpointsProvider{
		PrefixProvider: urlprefixProvider,
	}
	scopesValidator := _wireScopesValidatorValue
	tokenGenerator := _wireTokenGeneratorValue
	authorizationHandler := handler.ProvideAuthorizationHandler(context, tenantConfiguration, factory, authorizationStore, grantStore, endpointsProvider, endpointsProvider, scopesValidator, tokenGenerator, provider)
	httpHandler := provideAuthorizeHandler(factory, txContext, authorizationHandler)
	return httpHandler
}

var (
	_wireScopesValidatorValue = handler.ScopesValidator(oidc.ValidateScopes)
	_wireTokenGeneratorValue  = handler.TokenGenerator(handler.GenerateToken)
)

func newTokenHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	context := auth.ProvideContext(r)
	requestID := auth.ProvideLoggingRequestID(r)
	tenantConfiguration := auth.ProvideTenantConfig(context)
	factory := logging.ProvideLoggerFactory(context, requestID, tenantConfiguration)
	txContext := db.ProvideTxContext(context, tenantConfiguration)
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(context, tenantConfiguration)
	authorizationStore := &pq.AuthorizationStore{
		SQLBuilder:  sqlBuilder,
		SQLExecutor: sqlExecutor,
	}
	provider := time.NewProvider()
	grantStore := redis.ProvideGrantStore(context, tenantConfiguration, sqlBuilder, sqlExecutor, provider)
	store := redis2.ProvideStore(context, tenantConfiguration, provider, factory)
	eventStore := redis2.ProvideEventStore(context, tenantConfiguration)
	sessionProvider := session.ProvideSessionProvider(r, store, eventStore, tenantConfiguration)
	urlprefixProvider := urlprefix.NewProvider(r)
	idTokenIssuer := oidc.ProvideIDTokenIssuer(tenantConfiguration, urlprefixProvider, provider)
	tokenGenerator := _wireTokenGeneratorValue
	tokenHandler := handler.ProvideTokenHandler(context, tenantConfiguration, factory, authorizationStore, grantStore, sessionProvider, idTokenIssuer, tokenGenerator, provider)
	httpHandler := provideTokenHandler(factory, txContext, tokenHandler)
	return httpHandler
}

func newMetadataHandler(r *http.Request, m auth.DependencyMap) http.Handler {
	provider := urlprefix.NewProvider(r)
	endpointsProvider := &auth.EndpointsProvider{
		PrefixProvider: provider,
	}
	metadataProvider := &oauth.MetadataProvider{
		AuthorizeEndpoint:    endpointsProvider,
		TokenEndpoint:        endpointsProvider,
		AuthenticateEndpoint: endpointsProvider,
	}
	oidcMetadataProvider := &oidc.MetadataProvider{}
	httpHandler := provideMetadataHandler(metadataProvider, oidcMetadataProvider)
	return httpHandler
}

// wire.go:

func provideAuthorizeHandler(lf logging.Factory, tx db.TxContext, ah oauthAuthorizeHandler) http.Handler {
	h := &AuthorizeHandler{
		logger:       lf.NewLogger("oauth-authz-handler"),
		txContext:    tx,
		authzHandler: ah,
	}
	return h
}

func provideTokenHandler(lf logging.Factory, tx db.TxContext, th oauthTokenHandler) http.Handler {
	h := &TokenHandler{
		logger:       lf.NewLogger("oauth-token-handler"),
		txContext:    tx,
		tokenHandler: th,
	}
	return h
}

func provideMetadataHandler(oauth2 *oauth.MetadataProvider, oidc2 *oidc.MetadataProvider) http.Handler {
	h := &MetadataHandler{
		metaProviders: []oauthMetadataProvider{oauth2, oidc2},
	}
	return h
}
