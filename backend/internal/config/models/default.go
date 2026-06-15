package models

func DefaultConfig() *Config {
	return &Config{
		Items: []ConfigItem{
			{Key: DebugKey, Value: "false", EnvName: DebugEnvKey, FlagName: FlagDebug},
			{Key: EnvironmentKey, Value: "development", EnvName: EnvironmentEnvKey, FlagName: FlagEnvironment},

			// Logger Default Values
			{Key: LogLevelKey, Value: "info", EnvName: LogLevelEnvKey, FlagName: FlagLogLevel},
			{Key: LogFormatKey, Value: "text", EnvName: LogFormatEnvKey, FlagName: FlagLogFormat},
			{Key: LogEnableCallerKey, Value: "true", EnvName: LogEnableCallerEnvKey},
			{Key: LogUseStdoutKey, Value: "true", EnvName: LogUseStdoutEnvKey},
			{Key: LogFilePathKey, Value: "", EnvName: LogFilePathEnvKey},

			{Key: ServerAPIPortKey, Value: "5000", EnvName: ServerAPIPortEnvKey, FlagName: FlagAPIPort},
			{Key: ServerBindAddressKey, Value: "0.0.0.0", EnvName: ServerBindAddressEnvKey, FlagName: FlagBindTo},
			{Key: ServerBaseURLKey, Value: "http://localhost:5000", EnvName: ServerBaseURLEnvKey, FlagName: FlagBaseURL},
			{Key: ServerAPIPrefixKey, Value: "/api", EnvName: ServerAPIPrefixEnvKey, FlagName: FlagAPIPrefix},
			{Key: AuthRootPasswordKey, Value: "root", EnvName: AuthRootPasswordEnvKey, FlagName: FlagRootPassword},
			{Key: JwtAuthSecretKey, Value: "secret", EnvName: JwtAuthSecretEnvKey, FlagName: FlagJwtAuthSecret},
			{Key: JwtIssuerKey, Value: "locally-cli", EnvName: JwtIssuerEnvKey, FlagName: FlagJwtIssuer},
			{Key: EncryptionMasterSecretKey, Value: "default-master-secret-change-in-production", EnvName: EncryptionMasterSecretEnvKey, FlagName: FlagEncryptionMasterSecret},
			{Key: EncryptionGlobalSecretKey, Value: "default-global-secret-change-in-production", EnvName: EncryptionGlobalSecretEnvKey, FlagName: FlagEncryptionGlobalSecret},

			// Root User Default Values
			{Key: RootUserUsernameKey, Value: "root", EnvName: RootUserUsernameEnvKey, FlagName: FlagRootUserUsername},
			{Key: RootUserPasswordKey, Value: "root", EnvName: RootUserPasswordEnvKey, FlagName: FlagRootUserPassword},

			// Domain Default Values
			{Key: DomainNameKey, Value: DefaultLocallyDomain, EnvName: DomainNameEnvKey},
			{Key: DomainAdminEmailAddressKey, Value: DefaultCertificateAdminEmailAddress, EnvName: DomainAdminEmailAddressEnvKey},

			// Seeding Default Values
			{Key: SeedDemoDataKey, Value: "false", EnvName: SeedDemoDataEnvKey, FlagName: FlagSeedDemoData},

			// Pagination Default Values
			{Key: PaginationDefaultPageSizeKey, Value: DefaultPageSize, EnvName: PaginationDefaultPageSizeEnvKey, FlagName: FlagPaginationDefaultPageSize},

			// Security Default Values
			{Key: SecurityPasswordMinLengthKey, Value: "8", EnvName: SecurityPasswordMinLengthEnvKey, FlagName: FlagSecurityPasswordMinLength},
			{Key: SecurityPasswordRequireNumberKey, Value: "true", EnvName: SecurityPasswordRequireNumberEnvKey, FlagName: FlagSecurityPasswordRequireNumber},
			{Key: SecurityPasswordRequireSpecialKey, Value: "true", EnvName: SecurityPasswordRequireSpecialEnvKey, FlagName: FlagSecurityPasswordRequireSpecial},
			{Key: SecurityPasswordRequireUppercaseKey, Value: "true", EnvName: SecurityPasswordRequireUppercaseEnvKey, FlagName: FlagSecurityPasswordRequireUppercase},

			// Security Headers Default Values
			{Key: SecurityHeadersEnabledKey, Value: "true", EnvName: SecurityHeadersEnabledEnvKey, FlagName: FlagSecurityHeadersEnabled},
			{Key: SecurityHeadersHSTSMaxAgeKey, Value: "63072000", EnvName: SecurityHeadersHSTSMaxAgeEnvKey, FlagName: FlagSecurityHeadersHSTSMaxAge},
			{Key: SecurityHeadersHSTSIncludeSubdomainsKey, Value: "true", EnvName: SecurityHeadersHSTSIncludeSubdomainsEnvKey},
			{Key: SecurityHeadersCSPKey, Value: "default-src 'self'", EnvName: SecurityHeadersCSPEnvKey, FlagName: FlagSecurityHeadersCSP},
			{Key: SecurityHeadersFrameOptionsKey, Value: "DENY", EnvName: SecurityHeadersFrameOptionsEnvKey},
			{Key: SecurityHeadersContentTypeOptionsKey, Value: "nosniff", EnvName: SecurityHeadersContentTypeOptionsEnvKey},
			{Key: SecurityHeadersReferrerPolicyKey, Value: "strict-origin-when-cross-origin", EnvName: SecurityHeadersReferrerPolicyEnvKey},
			{Key: SecurityHeadersPermissionsPolicyKey, Value: "geolocation=(), microphone=(), camera=()", EnvName: SecurityHeadersPermissionsPolicyEnvKey},

			// API Key
			{Key: APIKey, Value: "sk-locally-"},

			// Cors
			{Key: CorsAllowOriginsKey, Value: "http://localhost:3000, http://127.0.0.1:3000", EnvName: CorsAllowOriginsEnvKey, FlagName: FlagCorsAllowOrigins},
			{Key: CorsAllowMethodsKey, Value: "GET, POST, PUT, DELETE, PATCH, OPTIONS, HEAD", EnvName: CorsAllowMethodsEnvKey, FlagName: FlagCorsAllowMethods},
			{Key: CorsAllowHeadersKey, Value: "Accept, Accept-Language, Content-Type, Content-Language, Origin, Authorization, X-Requested-With, X-Request-ID, X-HTTP-Method-Override, Cache-Control, X-Tenant-ID", EnvName: CorsAllowHeadersEnvKey, FlagName: FlagCorsAllowHeaders},
			{Key: CorsExposeHeadersKey, Value: "X-Request-ID", EnvName: CorsExposeHeadersEnvKey, FlagName: FlagCorsExposeHeaders},

			// Database Default Values
			{Key: DatabaseTypeKey, Value: "sqlite", EnvName: DatabaseTypeEnvKey, FlagName: FlagDatabaseType},
			{Key: DatabaseStoragePathKey, Value: "", EnvName: DatabaseStoragePathEnvKey, FlagName: FlagDatabaseStoragePath},
			{Key: DatabaseHostKey, Value: "localhost", EnvName: DatabaseHostEnvKey, FlagName: FlagDatabaseHost},
			{Key: DatabasePortKey, Value: "5432", EnvName: DatabasePortEnvKey, FlagName: FlagDatabasePort},
			{Key: DatabaseDatabaseKey, Value: "locally", EnvName: DatabaseDatabaseEnvKey, FlagName: FlagDatabaseDatabase},
			{Key: DatabaseUsernameKey, Value: "locally", EnvName: DatabaseUsernameEnvKey, FlagName: FlagDatabaseUsername},
			{Key: DatabasePasswordKey, Value: "locally", EnvName: DatabasePasswordEnvKey, FlagName: FlagDatabasePassword},
			{Key: DatabaseSSLModeKey, Value: "false", EnvName: DatabaseSSLModeEnvKey, FlagName: FlagDatabaseSSLMode},
			{Key: DatabaseMigrateKey, Value: "false", EnvName: DatabaseMigrateEnvKey, FlagName: FlagDatabaseMigrate},

			// Activity Default Values
			{Key: ActivityRetentionDaysKey, Value: "90", EnvName: ActivityRetentionDaysEnvKey, FlagName: FlagActivityRetentionDays},

			// vault
			{Key: VaultsConfigKey, Value: ""},
		},
	}
}
