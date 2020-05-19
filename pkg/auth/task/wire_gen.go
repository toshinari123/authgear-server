// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package task

import (
	"context"
	"github.com/google/wire"
	"github.com/skygeario/skygear-server/pkg/auth"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/authenticator/password"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/identity/loginid"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/userprofile"
	"github.com/skygeario/skygear-server/pkg/auth/dependency/userverify"
	"github.com/skygeario/skygear-server/pkg/core/async"
	"github.com/skygeario/skygear-server/pkg/core/auth/authinfo/pq"
	"github.com/skygeario/skygear-server/pkg/core/db"
	"github.com/skygeario/skygear-server/pkg/core/logging"
	"github.com/skygeario/skygear-server/pkg/core/mail"
	"github.com/skygeario/skygear-server/pkg/core/sms"
	"github.com/skygeario/skygear-server/pkg/core/time"
)

// Injectors from wire.go:

func newVerifyCodeSendTask(ctx context.Context, m auth.DependencyMap) async.Task {
	tenantConfiguration := auth.ProvideTenantConfig(ctx, m)
	engine := auth.ProvideTemplateEngine(tenantConfiguration, m)
	sender := mail.ProvideMailSender(ctx, tenantConfiguration)
	client := sms.ProvideSMSClient(ctx, tenantConfiguration)
	codeSenderFactory := userverify.NewDefaultUserVerifyCodeSenderFactory(tenantConfiguration, engine, sender, client)
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlExecutor := db.ProvideSQLExecutor(ctx, tenantConfiguration)
	store := pq.ProvideStore(sqlBuilderFactory, sqlExecutor)
	provider := time.NewProvider()
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	userprofileStore := userprofile.ProvideStore(provider, sqlBuilder, sqlExecutor)
	userverifyProvider := userverify.ProvideProvider(tenantConfiguration, provider, sqlBuilder, sqlExecutor)
	reservedNameChecker := auth.ProvideReservedNameChecker(m)
	typeCheckerFactory := loginid.ProvideTypeCheckerFactory(tenantConfiguration, reservedNameChecker)
	checker := loginid.ProvideChecker(tenantConfiguration, typeCheckerFactory)
	normalizerFactory := loginid.ProvideNormalizerFactory(tenantConfiguration)
	loginidProvider := loginid.ProvideProvider(sqlBuilder, sqlExecutor, provider, tenantConfiguration, checker, normalizerFactory)
	txContext := db.ProvideTxContext(ctx, tenantConfiguration)
	requestID := ProvideLoggingRequestID(ctx)
	factory := logging.ProvideLoggerFactory(ctx, requestID, tenantConfiguration)
	verifyCodeSendTask := &VerifyCodeSendTask{
		CodeSenderFactory:        codeSenderFactory,
		AuthInfoStore:            store,
		UserProfileStore:         userprofileStore,
		UserVerificationProvider: userverifyProvider,
		LoginIDProvider:          loginidProvider,
		TxContext:                txContext,
		LoggerFactory:            factory,
	}
	return verifyCodeSendTask
}

func newPwHouseKeeperTask(ctx context.Context, m auth.DependencyMap) async.Task {
	tenantConfiguration := auth.ProvideTenantConfig(ctx, m)
	txContext := db.ProvideTxContext(ctx, tenantConfiguration)
	requestID := ProvideLoggingRequestID(ctx)
	factory := logging.ProvideLoggerFactory(ctx, requestID, tenantConfiguration)
	provider := time.NewProvider()
	sqlBuilderFactory := db.ProvideSQLBuilderFactory(tenantConfiguration)
	sqlBuilder := auth.ProvideAuthSQLBuilder(sqlBuilderFactory)
	sqlExecutor := db.ProvideSQLExecutor(ctx, tenantConfiguration)
	historyStoreImpl := password.ProvideHistoryStore(provider, sqlBuilder, sqlExecutor)
	housekeeper := password.ProvideHousekeeper(historyStoreImpl, factory, tenantConfiguration)
	pwHousekeeperTask := &PwHousekeeperTask{
		TxContext:     txContext,
		LoggerFactory: factory,
		PwHousekeeper: housekeeper,
	}
	return pwHousekeeperTask
}

func newSendMessagesTask(ctx context.Context, m auth.DependencyMap) async.Task {
	tenantConfiguration := auth.ProvideTenantConfig(ctx, m)
	sender := mail.ProvideMailSender(ctx, tenantConfiguration)
	client := sms.ProvideSMSClient(ctx, tenantConfiguration)
	requestID := ProvideLoggingRequestID(ctx)
	factory := logging.ProvideLoggerFactory(ctx, requestID, tenantConfiguration)
	sendMessagesTask := &SendMessagesTask{
		EmailSender:   sender,
		SMSClient:     client,
		LoggerFactory: factory,
	}
	return sendMessagesTask
}

// wire.go:

var DependencySet = wire.NewSet(
	ProvideLoggingRequestID,
)
