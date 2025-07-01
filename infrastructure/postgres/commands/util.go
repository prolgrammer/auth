package commands

const (
	UserTable          = "users"
	UserIdField        = "id"
	UserEmailField     = "email"
	UserPasswordField  = "password"
	UserCreatedAtField = "created_at"
)

const (
	SessionTable            = "sessions"
	SessionIdField          = "id"
	SessionUserIdField      = "user_id"
	SessionRefreshTokenHash = "refresh_token_hash"
	SessionUserAgentField   = "user_agent"
	SessionIPField          = "ip_address"
	SessionCreatedAtField   = "created_at"
	SessionExpiresAtField   = "expires_at"
)
