package db

type Config struct {
	Database    string `yaml:"database"`
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	SslMode     string `yaml:"ssl_mode"`
	MaxAttempts int    `yaml:"max_attempts"`
}
