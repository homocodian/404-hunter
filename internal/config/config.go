package config

type Config struct {
	Workers    int
	Port       int
	OutputFile string
}

type Option func(*Config)

func WithPort(port int) Option {
	return func(c *Config) {
		c.Port = port
	}
}

func WithOutputFile(path string) Option {
	return func(c *Config) {
		c.OutputFile = path
	}
}

func NewServerConfig(workers int, port int) *Config {
	return NewConfig(workers, WithPort(port))
}

func NewConfig(workers int, opts ...Option) *Config {
	cfg := Config{
		Workers: workers,
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	return &cfg
}
