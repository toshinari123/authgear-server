package config

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"

	"github.com/authgear/authgear-server/pkg/util/validation"
)

type ServerConfig struct {
	// PublicListenAddr sets the listen address of the public server
	PublicListenAddr string `envconfig:"PUBLIC_LISTEN_ADDR" default:"0.0.0.0:3000"`
	// InternalListenAddr sets the listen address of the internal server
	InternalListenAddr string `envconfig:"INTERNAL_LISTEN_ADDR" default:"0.0.0.0:3001"`
	// TrustProxy sets whether HTTP headers from proxy are to be trusted
	TrustProxy bool `envconfig:"TRUST_PROXY" default:"false"`
	// DevMode sets whether the server would be run under development mode
	DevMode bool `envconfig:"DEV_MODE" default:"false"`
	// TLSCertFilePath sets the file path of TLS certificate.
	// It is required when development mode is enabled.
	// It is only used when development mode is enabled.
	TLSCertFilePath string `envconfig:"TLS_CERT_FILE_PATH" default:"tls-cert.pem"`
	// TLSKeyFilePath sets the file path of TLS private key.
	// It is required when development mode is enabled.
	// It is only used when development mode is enabled.
	TLSKeyFilePath string `envconfig:"TLS_KEY_FILE_PATH" default:"tls-key.pem"`
	// LogLevel sets the global log level
	LogLevel string `envconfig:"LOG_LEVEL" default:"warn"`
	// ConfigSource configures the source of app configurations
	ConfigSource ConfigurationSourceConfig `envconfig:"CONFIG_SOURCE"`

	// DefaultTemplateDirectory sets the directory for default template files
	DefaultTemplateDirectory string `envconfig:"DEFAULT_TEMPLATE_DIRECTORY" default:"templates"`
	// ReservedNameFilePath sets the file path for reserved name list
	ReservedNameFilePath string `envconfig:"RESERVED_NAME_FILE_PATH" default:"reserved_name.txt"`
	// StaticAsset configures serving static asset
	StaticAsset ServerStaticAssetConfig `envconfig:"STATIC_ASSET"`

	// SentryDSN sets the sentry DSN.
	SentryDSN string `envconfig:"SENTRY_DSN"`
}

func LoadServerConfigFromEnv() (*ServerConfig, error) {
	config := &ServerConfig{}

	err := envconfig.Process("", config)
	if err != nil {
		return nil, fmt.Errorf("cannot load server config: %w", err)
	}

	err = config.Validate()
	if err != nil {
		return nil, fmt.Errorf("invalid server config: %w", err)
	}

	return config, nil
}

func (c *ServerConfig) Validate() error {
	ctx := &validation.Context{}

	if c.StaticAsset.ServingEnabled && c.StaticAsset.URLPrefix == "" {
		ctx.Child("STATIC_ASSET_URL_PREFIX").EmitErrorMessage(
			"static asset URL prefix must be set when static assets are not served",
		)
	}

	switch c.ConfigSource.Type {
	case SourceTypeLocalFile:
		break
	default:
		sourceTypes := make([]string, len(SourceTypes))
		for i, t := range SourceTypes {
			sourceTypes[i] = string(t)
		}
		ctx.Child("CONFIG_SOURCE_TYPE").EmitErrorMessage(
			"invalid configuration source type; available: " + strings.Join(sourceTypes, ", "),
		)
	}

	return ctx.Error("invalid server configuration")
}

type ServerStaticAssetConfig struct {
	// ServingEnabled sets whether serving static assets is enabled
	ServingEnabled bool `envconfig:"SERVING_ENABLED" default:"true"`
	// Dir sets the local directory of static assets
	Dir string `envconfig:"DIR" default:"./static"`
	// URLPrefix sets the URL prefix for static assets
	URLPrefix string `envconfig:"URL_PREFIX" default:"/static"`
}

type SourceType string

const (
	SourceTypeLocalFile SourceType = "local_file"
)

var SourceTypes = []SourceType{
	SourceTypeLocalFile,
}

type ConfigurationSourceConfig struct {
	// Type sets the type of configuration source
	Type SourceType `envconfig:"TYPE" default:"local_file"`

	// AppConfigPath sets the path to app configuration YAML file for local file source
	AppConfigPath string `envconfig:"APP_CONFIG_PATH" default:"authgear.yaml"`
	// SecretConfigPath sets the path to secret configuration YAML file for local file source
	SecretConfigPath string `envconfig:"SECRET_CONFIG_PATH" default:"authgear.secrets.yaml"`
}