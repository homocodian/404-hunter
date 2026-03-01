package config

type Config struct {
	Workers int
	Port    int
}

func NewConfig(workers int, port int) *Config {
	return &Config{
		Workers: workers,
		Port:    port,
	}
}
