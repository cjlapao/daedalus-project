package models

const (
	// UUIDs and IDs
	// UUIDs and IDs
	// Moved to internal/models/constants.go
	// UnknownTenantID, UnknownUserID, etc.

	SystemStoragePath         = ".locally"
	DefaultStoragePath        = ".locally" // Mapped SystemStoragePath
	DefaultPageSizeInt        = 20
	DefaultPageSize           = "20"
	DefaultRetentionDays      = 90
	DefaultJwtExpirationHours = 24

	// Headers
	APIKeyAuthorizationHeader = "X-API-KEY"
	TenantIDHeader            = "X-Tenant-ID"
	CorrelationIDHeader       = "X-Correlation-ID"
	RequestIDHeader           = "X-Request-ID"
	SecurityLevelHeader       = "X-Security-Level"
	UserAgentHeader           = "X-User-Agent"
	UserIPHeader              = "X-User-IP"
	UsernameHeader            = "X-Username"
	StartTimeHeader           = "X-Start-Time"

	// Certificates Defaults
	DefaultRootCertificateCommonName     = "Locally"
	DefaultCertificateSubDomain          = "Locally Root CA"
	DefaultLocallyDomain                 = "locally.local"
	DefaultCertificateFQDN               = "*.locally.local"
	DefaultCertificateExpiresInYears     = 10
	DefaultCertificateKeySize            = 4096
	DefaultCertificateSignatureAlgorithm = "SHA512"
	DefaultCertificateCountry            = "UK"
	DefaultCertificateState              = "London"
	DefaultCertificateCity               = "London"
	DefaultCertificateOrganization       = "Locally"
	DefaultCertificateOrganizationalUnit = "Locally IT"
	DefaultCertificateAdminEmailAddress  = "admin@locally.local"

	VariablesPrefix string = "${{"

	// Config
	ConfigFilePathKey  = "config.file_path"
	ConfigFilePathEnv  = "CONFIG_FILE_PATH"
	ConfigFilePathFlag = "config-path"

	// Keys
	DebugKey       = "debug"
	EnvironmentKey = "environment"

	// Logger
	LogLevelKey        = "logger.level"
	LogFormatKey       = "logger.format"
	LogEnableCallerKey = "logger.enable_caller"
	LogUseStdoutKey    = "logger.use_stdout"
	LogFilePathKey     = "logger.file_path"

	// Application
	ApplicationNameKey        = "application.name"
	ApplicationVersionKey     = "application.version"
	ApplicationEnvironmentKey = "application.environment"
	ApplicationGitCommitKey   = "application.git_commit"
	ApplicationBuildDateKey   = "application.build_date"
	ApplicationAuthorKey      = "application.author"

	// Server
	ServerAPIPortKey     = "server.api.port"
	ServerBindAddressKey = "server.api.bind_address"
	ServerBaseURLKey     = "server.base_url"
	ServerAPIPrefixKey   = "server.api_prefix"

	// Auth
	AuthRootPasswordKey            = "auth.root_password"
	JwtAuthSecretKey               = "auth.jwt.secret"
	JwtIssuerKey                   = "auth.jwt.issuer"
	JwtExpirationKey               = "auth.jwt.expiration"
	AuthValidationTokenDurationKey = "auth.validation_token_duration"

	// Encryption
	EncryptionMasterSecretKey = "encryption.master_secret"
	EncryptionGlobalSecretKey = "encryption.global_secret"

	// Database
	DatabaseTypeKey        = "database.type"
	DatabaseStoragePathKey = "database.storage_path"
	DatabaseHostKey        = "database.host"
	DatabasePortKey        = "database.port"
	DatabaseDatabaseKey    = "database.name"
	DatabaseUsernameKey    = "database.username"
	DatabasePasswordKey    = "database.password"
	DatabaseSSLModeKey     = "database.ssl_mode"
	DatabaseMigrateKey     = "database.migrate"

	// Activity
	ActivityRetentionDaysKey = "activity.retention_days"

	// Pagination
	PaginationDefaultPageSizeKey = "pagination.default_page_size"

	// CORS
	CorsAllowOriginsKey  = "cors.allow_origins"
	CorsAllowMethodsKey  = "cors.allow_methods"
	CorsAllowHeadersKey  = "cors.allow_headers"
	CorsExposeHeadersKey = "cors.expose_headers"

	// Context
	TenantIDContextKey = "tenant_id"

	// Security (New)
	SecurityPasswordMinLengthKey        = "security.password.min_length"
	SecurityPasswordRequireNumberKey    = "security.password.require.number"
	SecurityPasswordRequireSpecialKey   = "security.password.require.special"
	SecurityPasswordRequireUppercaseKey = "security.password.require.uppercase"

	// Security Headers
	SecurityHeadersEnabledKey               = "security.headers.enabled"
	SecurityHeadersHSTSMaxAgeKey            = "security.headers.hsts.max_age"
	SecurityHeadersHSTSIncludeSubdomainsKey = "security.headers.hsts.include_subdomains"
	SecurityHeadersCSPKey                   = "security.headers.csp"
	SecurityHeadersFrameOptionsKey          = "security.headers.frame_options"
	SecurityHeadersContentTypeOptionsKey    = "security.headers.content_type_options"
	SecurityHeadersReferrerPolicyKey        = "security.headers.referrer_policy"
	SecurityHeadersPermissionsPolicyKey     = "security.headers.permissions_policy"

	// Vaults
	VaultsConfigKey = "vaults"

	// Env Keys
	DebugEnvKey       = "DEBUG"
	EnvironmentEnvKey = "ENVIRONMENT"

	// Logger Env Keys
	LogLevelEnvKey        = "LOG_LEVEL"
	LogFormatEnvKey       = "LOG_FORMAT"
	LogEnableCallerEnvKey = "LOG_ENABLE_CALLER"
	LogUseStdoutEnvKey    = "LOG_USE_STDOUT"
	LogFilePathEnvKey     = "LOG_FILE_PATH"

	ServerAPIPortEnvKey             = "SERVER_PORT"
	ServerBindAddressEnvKey         = "SERVER_BIND_ADDRESS"
	ServerBaseURLEnvKey             = "SERVER_BASE_URL"
	ServerAPIPrefixEnvKey           = "SERVER_API_PREFIX"
	AuthRootPasswordEnvKey          = "AUTH_ROOT_PASSWORD"
	JwtAuthSecretEnvKey             = "JWT_SECRET"
	JwtIssuerEnvKey                 = "JWT_ISSUER"
	JwtExpirationEnvKey             = "JWT_EXPIRATION"
	EncryptionMasterSecretEnvKey    = "ENCRYPTION_MASTER_SECRET"
	EncryptionGlobalSecretEnvKey    = "ENCRYPTION_GLOBAL_SECRET"
	DatabaseTypeEnvKey              = "DATABASE_TYPE"
	DatabaseStoragePathEnvKey       = "DATABASE_STORAGE_PATH"
	DatabaseHostEnvKey              = "DATABASE_HOST"
	DatabasePortEnvKey              = "DATABASE_PORT"
	DatabaseDatabaseEnvKey          = "DATABASE_NAME"
	DatabaseUsernameEnvKey          = "DATABASE_USERNAME"
	DatabasePasswordEnvKey          = "DATABASE_PASSWORD"
	DatabaseSSLModeEnvKey           = "DATABASE_SSL_MODE"
	DatabaseMigrateEnvKey           = "DATABASE_MIGRATE"
	ActivityRetentionDaysEnvKey     = "ACTIVITY_RETENTION_DAYS"
	PaginationDefaultPageSizeEnvKey = "PAGINATION_DEFAULT_PAGE_SIZE"
	CorsAllowOriginsEnvKey          = "CORS_ALLOW_ORIGINS"
	CorsAllowMethodsEnvKey          = "CORS_ALLOW_METHODS"
	CorsAllowHeadersEnvKey          = "CORS_ALLOW_HEADERS"
	CorsExposeHeadersEnvKey         = "CORS_EXPOSE_HEADERS"

	SecurityPasswordMinLengthEnvKey        = "SECURITY_PASSWORD_MIN_LENGTH"
	SecurityPasswordRequireNumberEnvKey    = "SECURITY_PASSWORD_REQUIRE_NUMBER"
	SecurityPasswordRequireSpecialEnvKey   = "SECURITY_PASSWORD_REQUIRE_SPECIAL"
	SecurityPasswordRequireUppercaseEnvKey = "SECURITY_PASSWORD_REQUIRE_UPPERCASE"

	SecurityHeadersEnabledEnvKey               = "SECURITY_HEADERS_ENABLED"
	SecurityHeadersHSTSMaxAgeEnvKey            = "SECURITY_HEADERS_HSTS_MAX_AGE"
	SecurityHeadersHSTSIncludeSubdomainsEnvKey = "SECURITY_HEADERS_HSTS_INCLUDE_SUBDOMAINS"
	SecurityHeadersCSPEnvKey                   = "SECURITY_HEADERS_CSP"
	SecurityHeadersFrameOptionsEnvKey          = "SECURITY_HEADERS_FRAME_OPTIONS"
	SecurityHeadersContentTypeOptionsEnvKey    = "SECURITY_HEADERS_CONTENT_TYPE_OPTIONS"
	SecurityHeadersReferrerPolicyEnvKey        = "SECURITY_HEADERS_REFERRER_POLICY"
	SecurityHeadersPermissionsPolicyEnvKey     = "SECURITY_HEADERS_PERMISSIONS_POLICY"

	// Flags
	FlagDebug       = "debug"
	FlagEnvironment = "environment"
	// Logger Flags
	FlagLogLevel  = "log-level"
	FlagLogFormat = "log-format"

	FlagAPIPort                   = "port"
	FlagBindTo                    = "bind-to"
	FlagBaseURL                   = "base-url"
	FlagAPIPrefix                 = "api-prefix"
	FlagRootPassword              = "root-password"
	FlagJwtAuthSecret             = "jwt-secret"
	FlagJwtIssuer                 = "jwt-issuer"
	FlagJwtExpiration             = "jwt-expiration"
	FlagEncryptionMasterSecret    = "encryption-master-secret"
	FlagEncryptionGlobalSecret    = "encryption-global-secret"
	FlagDatabaseType              = "db-type"
	FlagDatabaseStoragePath       = "db-path"
	FlagDatabaseHost              = "db-host"
	FlagDatabasePort              = "db-port"
	FlagDatabaseDatabase          = "db-name"
	FlagDatabaseUsername          = "db-user"
	FlagDatabasePassword          = "db-password"
	FlagDatabaseSSLMode           = "db-ssl-mode"
	FlagDatabaseMigrate           = "db-migrate"
	FlagActivityRetentionDays     = "activity-retention-days"
	FlagPaginationDefaultPageSize = "default-page-size"
	FlagCorsAllowOrigins          = "cors-allow-origins"
	FlagCorsAllowMethods          = "cors-allow-methods"
	FlagCorsAllowHeaders          = "cors-allow-headers"
	FlagCorsExposeHeaders         = "cors-expose-headers"

	FlagSecurityPasswordMinLength        = "password-min-length"
	FlagSecurityPasswordRequireNumber    = "password-require-number"
	FlagSecurityPasswordRequireSpecial   = "password-require-special"
	FlagSecurityPasswordRequireUppercase = "password-require-uppercase"

	FlagSecurityHeadersEnabled    = "security-headers-enabled"
	FlagSecurityHeadersHSTSMaxAge = "hsts-max-age"
	FlagSecurityHeadersCSP        = "csp"

	// Additional keys from default.go
	RootUserUsernameKey    = "root_user.username"
	RootUserPasswordKey    = "root_user.password"
	RootUserUsernameEnvKey = "ROOT_USER_USERNAME"
	RootUserPasswordEnvKey = "ROOT_USER_PASSWORD"
	FlagRootUserUsername   = "root-username"
	FlagRootUserPassword   = "root-user-password"

	DomainNameKey                 = "domain.name"
	DomainAdminEmailAddressKey    = "domain.admin_email"
	DomainNameEnvKey              = "DOMAIN_NAME"
	DomainAdminEmailAddressEnvKey = "DOMAIN_ADMIN_EMAIL"

	SeedDemoDataKey    = "seed.demo_data"
	SeedDemoDataEnvKey = "SEED_DEMO_DATA"
	FlagSeedDemoData   = "seed-demo-data"

	// Certificate Keys
	CertificateExpiresInYearsKey     = "certificate.expires_in_years"
	CertificateKeySizeKey            = "certificate.key_size"
	CertificateSignatureAlgorithmKey = "certificate.signature_algorithm"
	CertificateCountryKey            = "certificate.country"
	CertificateStateKey              = "certificate.state"
	CertificateCityKey               = "certificate.city"
	CertificateOrganizationKey       = "certificate.organization"
	CertificateOrganizationalUnitKey = "certificate.organizational_unit"
	CertificateAdminEmailAddressKey  = "certificate.admin_email"

	CertificateExpiresInYearsEnvKey     = "CERTIFICATE_EXPIRES_IN_YEARS"
	CertificateKeySizeEnvKey            = "CERTIFICATE_KEY_SIZE"
	CertificateSignatureAlgorithmEnvKey = "CERTIFICATE_SIGNATURE_ALGORITHM"
	CertificateCountryEnvKey            = "CERTIFICATE_COUNTRY"
	CertificateStateEnvKey              = "CERTIFICATE_STATE"
	CertificateCityEnvKey               = "CERTIFICATE_CITY"
	CertificateOrganizationEnvKey       = "CERTIFICATE_ORGANIZATION"
	CertificateOrganizationalUnitEnvKey = "CERTIFICATE_ORGANIZATIONAL_UNIT"
	CertificateAdminEmailAddressEnvKey  = "CERTIFICATE_ADMIN_EMAIL"

	APIKey = "api_key"

	// MessageProcessor
	MessageProcessorDefaultMaxRetriesKey    = "message_processor.default_max_retries"
	MessageProcessorPollIntervalKey         = "message_processor.poll_interval"
	MessageProcessorProcessingTimeoutKey    = "message_processor.processing_timeout"
	MessageProcessorRecoveryEnabledKey      = "message_processor.recovery_enabled"
	MessageProcessorMaxProcessingAgeKey     = "message_processor.max_processing_age"
	MessageProcessorCleanupEnabledKey       = "message_processor.cleanup_enabled"
	MessageProcessorCleanupMaxAgeKey        = "message_processor.cleanup_max_age"
	MessageProcessorCleanupIntervalKey      = "message_processor.cleanup_interval"
	MessageProcessorKeepCompleteMessagesKey = "message_processor.keep_complete_messages"
	MessageProcessorDebugKey                = "message_processor.debug"

	MessageProcessorDefaultMaxRetriesEnvKey    = "MESSAGE_PROCESSOR_MAX_RETRIES"
	MessageProcessorPollIntervalEnvKey         = "MESSAGE_PROCESSOR_POLL_INTERVAL"
	MessageProcessorProcessingTimeoutEnvKey    = "MESSAGE_PROCESSOR_PROCESSING_TIMEOUT"
	MessageProcessorRecoveryEnabledEnvKey      = "MESSAGE_PROCESSOR_RECOVERY_ENABLED"
	MessageProcessorMaxProcessingAgeEnvKey     = "MESSAGE_PROCESSOR_MAX_PROCESSING_AGE"
	MessageProcessorCleanupEnabledEnvKey       = "MESSAGE_PROCESSOR_CLEANUP_ENABLED"
	MessageProcessorCleanupMaxAgeEnvKey        = "MESSAGE_PROCESSOR_CLEANUP_MAX_AGE"
	MessageProcessorCleanupIntervalEnvKey      = "MESSAGE_PROCESSOR_CLEANUP_INTERVAL"
	MessageProcessorKeepCompleteMessagesEnvKey = "MESSAGE_PROCESSOR_KEEP_COMPLETE_MESSAGES"
	MessageProcessorDebugEnvKey                = "MESSAGE_PROCESSOR_DEBUG"

	FlagMessageProcessorDefaultMaxRetries    = "mp-max-retries"
	FlagMessageProcessorPollInterval         = "mp-poll-interval"
	FlagMessageProcessorProcessingTimeout    = "mp-timeout"
	FlagMessageProcessorRecoveryEnabled      = "mp-recovery"
	FlagMessageProcessorMaxProcessingAge     = "mp-max-age"
	FlagMessageProcessorCleanupEnabled       = "mp-cleanup"
	FlagMessageProcessorCleanupMaxAge        = "mp-cleanup-max-age"
	FlagMessageProcessorCleanupInterval      = "mp-cleanup-interval"
	FlagMessageProcessorKeepCompleteMessages = "mp-keep-complete"
	FlagMessageProcessorDebug                = "mp-debug"
)
