package config

type App struct {
	Host     string `mapstructure:"APP_HOST"`
	Port     string `mapstructure:"APP_PORT"`
	LogLevel string `mapstructure:"APP_LOG_LEVEL"`
}

type Database struct {
	Host string `mapstructure:"POSTGRES_HOST"`
	Port string `mapstructure:"POSTGRES_PORT"`
	User string `mapstructure:"POSTGRES_USER"`
	Pass string `mapstructure:"POSTGRES_PASS"`
	Name string `mapstructure:"POSTGRES_NAME"`
}

type MinIO struct {
	Endpoint  string `mapstructure:"MINIO_ENDPOINT"`
	AccessKey string `mapstructure:"MINIO_ACCESS_KEY"`
	SecretKey string `mapstructure:"MINIO_SECRET_KEY"`
	Bucket    string `mapstructure:"MINIO_BUCKET"`
	UseSSL    bool   `mapstructure:"MINIO_USE_SSL"`
}

type Redis struct {
	Host string `mapstructure:"REDIS_HOST"`
	Port string `mapstructure:"REDIS_PORT"`
}

type Config struct {
	App      App      `mapstructure:",squash"`
	Database Database `mapstructure:",squash"`
	MinIO    MinIO    `mapstructure:",squash"`
	Redis    Redis    `mapstructure:",squash"`
}
