package config

type Config struct {
	ServiceCfg Service  `yaml:"service"`
	DBCfg      Database `yaml:"database"`
	LoggerCfg  Logger   `yaml:"logger"`
}

type Service struct {
	Host            string `yaml:"host" env-default:"localhost"`
	Port            string `yaml:"port" env-default:"80"`
	CacheTtlSeconds int64  `yaml:"cacheTtlSeconds" env-default:"60"`
}

type Database struct {
	Host         string `yaml:"host" env-default:"localhost"`
	Port         string `yaml:"port" env-default:"5432"`
	DatabaseName string `yaml:"databaseName" env-default:"cars"`
	User         string `yaml:"user" env-default:"postgres"`
	Password     string `yaml:"password"`
}

type Logger struct {
	Level string `yaml:"level" env-default:"info"`
}
