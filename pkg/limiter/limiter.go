package limiter

type Limiter interface {
	Allow(key string) (count int32, err error)
	Close(key string) error
	CloseAll() error
}
