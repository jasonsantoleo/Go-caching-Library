package factory

// Type of cache,chooses the backend it want to use
type BackendType string

const (
	Memory    BackendType = "memory"
	Redis     BackendType = "redis"
	Memcached BackendType = "memcached"
)

// configuration of the avaliable backend
type Config struct {
	// In-Memory  config
	MemoryMaxSize int

	// Redis  config
	RedisAddr     string
	RedisPassword string
	RedisDB       int

	// Memcached  config
	MemcachedServers []string
}

// returns default config
func DefaultMemoryConfig() Config {
	return Config{
		MemoryMaxSize: 0,
	}
}
