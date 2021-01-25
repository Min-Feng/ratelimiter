package limiter

type Limiter interface {
	Allow(key string) (count int32, err error)
	Delete(key string) error
}
