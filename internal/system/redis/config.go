package redis

type Config struct {
	Addr     string
	DB       int
	Password string

	SentinelAddrs  []string
	SentinelMaster string
}
