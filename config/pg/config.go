package pg

type Config struct {
	User           string `mapstructure:"user"`
	Password       string `mapstructure:"password"`
	Host           string `mapstructure:"host"`
	Port           string `mapstructure:"port"`
	Database       string `mapstructure:"database"`
	MigrationsPath string `mapstructure:"migrations_path"`

	RetryConnectionAttempts int `mapstructure:"retry_connection_attempts"`
	RetryConnectionTimeout  int `mapstructure:"retry_connection_timeout"`
}
